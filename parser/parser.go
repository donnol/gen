package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	testSuffix = "_test"
)

// Parser 解析器，用于解析源码
type Parser struct {
	ignoreFileSuffix []string
}

// Option 选项
type Option struct {
	IgnoreFileSuffix []string // 忽略文件名含有指定后缀的文件
}

// New 新建
func New(opts ...Option) *Parser {
	p := &Parser{}
	for _, opt := range opts {
		p.ignoreFileSuffix = append(p.ignoreFileSuffix, opt.IgnoreFileSuffix...)
	}
	return p
}

// ParsePkg 解析包，如：github.com/pkg/errors，返回包信息
func (p *Parser) ParsePkg(pkg string) (Pkg, error) {
	var result Pkg

	// 获取包所在目录
	buildPkg, err := build.Import(pkg, "", build.ImportComment)
	if err != nil {
		return result, errors.WithStack(err)
	}
	result.Dir = buildPkg.Dir

	fset := token.NewFileSet()

	// 解析目录，获得ast.File
	pkgs, err := parser.ParseDir(fset, buildPkg.Dir, func(fi os.FileInfo) bool {
		li := strings.LastIndex(fi.Name(), filepath.Ext(fi.Name()))

		// 跳过test文件
		testi := strings.LastIndex(fi.Name(), testSuffix)
		if testi != -1 && li-testi == len(testSuffix) {
			return false
		}

		// 跳过指定文件后缀
		if len(p.ignoreFileSuffix) > 0 {
			for _, suf := range p.ignoreFileSuffix {
				sufi := strings.LastIndex(fi.Name(), suf)
				if sufi != -1 && li-sufi == len(suf) {
					return false
				}
			}
		}

		return true
	}, parser.ParseComments)
	if err != nil {
		return result, errors.WithStack(err)
	}

	// 收集包信息
	conf := &types.Config{
		IgnoreFuncBodies: true,
		FakeImportC:      false,
		Error: func(err error) {
			log.Printf("types check failed: %+v\n", err)
		},
		Importer:                 importer.Default(),
		DisableUnusedImportCheck: true,
	}

	for pkgName, parserPkg := range pkgs {
		result.Name = pkgName
		result.ImportPath = pkg

		// 包内文件解析
		files := make([]*ast.File, 0, len(parserPkg.Files))
		for _, file := range parserPkg.Files {
			files = append(files, file)
		}

		info := types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
			Scopes:     make(map[ast.Node]*types.Scope),
		}
		var typPkg *types.Package
		typPkg, err = conf.Check(pkg, fset, files, &info)
		if err != nil {
			return result, errors.WithStack(err)
		}

		// NOTE:这里会不会存在相同的type，导致键冲突呢？
		typsExprMap := make(map[types.Type]ast.Expr)
		for expr, tv := range info.Types {
			typsExprMap[tv.Type] = expr
		}
		result.typsExprMap = typsExprMap

		// 获取文档和注释
		docMap, commentMap := p.getStructDocAndComment(pkg, parserPkg)

		// 遍历作用域，返回结构体信息
		result.Structs = p.getStructInfo(&result, typPkg, info, typsExprMap, docMap, commentMap)
	}

	return result, nil
}

func (p *Parser) getStructDocAndComment(pkg string, parserPkg *ast.Package) (
	map[string]string,
	map[string]string,
) {
	var docMap = make(map[string]string)
	var commentMap = make(map[string]string)
	docPkg := doc.New(parserPkg, pkg, doc.AllDecls)
	for _, typ := range docPkg.Types {
		docMap[typ.Name] = typ.Doc

		for _, spec := range typ.Decl.Specs {
			// 找到结构体
			typSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, field := range structType.Fields.List {
				var fieldName, fieldDoc, fieldComment string
				// 内嵌结构体需要使用ast.Ident来获取名称
				if len(field.Names) != 0 {
					for _, name := range field.Names {
						fieldName = name.Name
					}
				} else {
					astField, ok := field.Type.(*ast.Ident)
					if ok {
						fieldName = astField.Name
					} else if selField, ok := field.Type.(*ast.SelectorExpr); ok {
						fieldName = selField.Sel.Name
					}
				}
				if fieldName == "" {
					fmt.Printf("Can't find field name: %+v\n", field)
				}
				if field.Doc != nil {
					fieldDoc = strings.TrimSpace(field.Doc.Text())
				}
				if field.Comment != nil {
					fieldComment = strings.TrimSpace(field.Comment.Text())
				}
				key := p.getStructDocAndCommentMapKey(typ.Name, fieldName)
				docMap[key] = fieldDoc
				commentMap[key] = fieldComment
			}
		}
	}

	return docMap, commentMap
}

func (p *Parser) getStructDocAndCommentMapKey(typName, fieldName string) string {
	return fmt.Sprintf("%s;%s", typName, fieldName)
}

func (p *Parser) getStructInfo(
	pkg *Pkg,
	typPkg *types.Package,
	typesInfo types.Info,
	typsExprMap map[types.Type]ast.Expr,
	docMap map[string]string,
	commentMap map[string]string,
) []Struct {
	var result []Struct

	for ident, obj := range typesInfo.Defs {
		exprStr := types.ExprString(ident)

		if obj == nil {
			continue
		}
		if ident.Obj != nil && ident.Obj.Kind != ast.Typ {
			continue
		}
		debug("ident: %+v, %+v\n", ident, ident.Obj)

		objType := obj.Type()
		objStruct, ok := objType.Underlying().(*types.Struct)
		if !ok {
			continue
		}
		expr, ok := typsExprMap[objType]
		if ok {
			debug("==expr: %+v, %v\n", expr, types.ExprString(expr))
		}
		debug("===tv: %+v, %+v, %+v\n", exprStr, obj.Name(), objStruct)

		result = append(result, p.getStructFromObject(pkg, obj, docMap, commentMap))
	}

	return result
}

func (p *Parser) getStructFromPkgScope(
	pkg *Pkg,
	scope *types.Scope,
	docMap map[string]string,
	commentMap map[string]string,
) []Struct {
	var result []Struct

	for _, scopeName := range scope.Names() {
		// 找到对象
		obj := scope.Lookup(scopeName)

		tmp := p.getStructFromObject(pkg, obj, docMap, commentMap)
		if tmp.Name == "" {
			continue
		}

		result = append(result, tmp)
	}

	return result
}

func (p *Parser) getStructFromObject(
	pkg *Pkg,
	obj types.Object,
	docMap map[string]string,
	commentMap map[string]string,
) Struct {
	tmp := Struct{}

	// 解析对象
	// 找出结构体
	objType := obj.Type()
	objStruct, ok := objType.Underlying().(*types.Struct)
	if !ok {
		return tmp
	}

	// 为结果赋值
	tmp.Info.InitWithTypes(pkg, obj.Name(), obj.Pkg().Path(), docMap[obj.Name()], "", objType)
	// 字段
	for i := 0; i < objStruct.NumFields(); i++ {
		field := objStruct.Field(i)
		key := p.getStructDocAndCommentMapKey(obj.Name(), field.Name())
		tmpField := Field{}
		tmpField.Info.InitWithTypes(pkg, field.Name(), "", docMap[key], commentMap[key], field.Type())
		tmpField.Anonymous = field.Anonymous()
		tmp.Fields = append(tmp.Fields, tmpField)
	}

	return tmp
}
