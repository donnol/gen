package list

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/donnol/gen/parser"
	"github.com/donnol/gen/template"
	toollist "github.com/donnol/tools/list"
	"github.com/pkg/errors"
)

// 生成文件后缀
const (
	GenFileSuffix = "_gen_list"
)

// 指令
const (
	commandName  = "list"
	attrColumn   = "column"
	attrMap      = "map"
	attrSliceMap = "slicemap"

	joinCommandName = "join"
	joinEqual       = "="

	deriveCommandName = "derive"
)

// List 列表，解析结构体，生成列表结构体及取列，列映射，列数组映射等方法
type List struct {
	parser   *parser.Parser
	template *template.Template

	useAnnotation bool // 是否使用注解
	methods       []Method
}

// New 新建
func New(p *parser.Parser, t *template.Template, opt Option) *List {
	return &List{
		parser:        p,
		template:      t,
		useAnnotation: opt.UseAnnotation,
		methods:       opt.Methods,
	}
}

// Option 选项
type Option struct {
	UseAnnotation bool     // 是否使用注解
	Methods       []Method // 方法
}

type Method struct {
	Name  string
	Attrs []string
}

// Parse 解析，对输入的导入路径进行解析，生成需要的结构体和方法，再写到同目录下的特定文件内
func (list *List) Parse(importPath string) error {
	// 解析
	pkg, err := list.parser.ParsePkg(importPath)
	if err != nil {
		return errors.WithStack(err)
	}

	// 生成
	pkgName, importPathMap, content, err := list.output(pkg)
	if err != nil {
		return err
	}

	// 没有内容，则不生成文件
	if len(content) == 0 {
		return nil
	}

	// 写入，添加包名和依赖导入
	importPaths := parser.ImportPathMap(importPathMap).Keys()
	importPaths = toollist.Filter(importPaths, pkg.ImportPath)
	fileContent := list.template.SpliceFile(template.FileOption{
		PkgName:     pkgName,
		ImportPaths: importPaths,
		Content:     string(content),
	})
	err = list.template.WriteFile(list.getFileName(pkg.Dir, pkgName), []byte(fileContent))
	if err != nil {
		return err
	}

	return nil
}

func (list *List) output(pkg parser.Pkg) (string, parser.ImportPathMap, []byte, error) {
	var pkgName = pkg.Name
	var importPathMap = make(parser.ImportPathMap)
	var content []byte

	var buf = bytes.NewBuffer(content)
	for _, singleStruct := range pkg.Structs {
		importPathMap = importPathMap.MergeKey(
			singleStruct.Info.GetNotEmptyImportPathMap())

		structName := singleStruct.Name
		typNameWithPath := singleStruct.Info.GetTypNameWithPath(pkg.ImportPath)

		existFieldCommand := false
		for _, singleField := range singleStruct.Fields {
			if singleField.Info.Commands.ExistCommand(commandName) ||
				singleField.Info.Commands.ExistCommand(joinCommandName) {
				existFieldCommand = true
				break
			}
		}
		if !list.useAnnotation || singleStruct.Info.Commands.ExistCommand(commandName) ||
			existFieldCommand {
			for _, text := range []string{
				typText,
				whereMethodText,
				lenMethodText,
				sortMethodText,
				limitMethodText,
				mapMethodText,
				reduceMethodText,
				eachMethodText,
				shuffleMethodText,
				reverseMethodText,
				firstMethodText,
				lastMethodText,
			} {
				if err := list.template.Execute(buf, "List", text, map[string]interface{}{
					"typName":         structName,
					"typNameWithPath": typNameWithPath,
				}); err != nil {
					return pkgName, importPathMap, content, errors.WithStack(err)
				}
			}
		}

		for _, singleField := range singleStruct.Fields {
			importPathMap = importPathMap.MergeKey(singleField.Info.GetNotEmptyImportPathMap())

			fieldName := singleField.Name
			fieldTypNameWithPath := singleField.Info.GetTypNameWithPath(pkg.ImportPath)
			fieldNameWithInner := fieldName
			extraTyp, extraTypWithPath := singleField.Info.Commands.GetExtraTyp(list.parser)
			if singleField.Anonymous {
				innerFieldName := strings.TrimLeft(extraTyp, ".")
				fieldName += innerFieldName
				fieldNameWithInner += extraTyp
				_, fieldTypNameWithPath = singleField.GetTypesTypeStructField(innerFieldName)
			}

			// 取列
			if singleField.Info.Commands.ExistCommandAttr(commandName, attrColumn) {
				if err := list.template.Execute(buf, "List", columnMethodText, map[string]interface{}{
					"typName":            structName,
					"fieldName":          fieldName,
					"fieldNameWithInner": fieldNameWithInner,
					"fieldType":          fieldTypNameWithPath,
				}); err != nil {
					return pkgName, importPathMap, content, errors.WithStack(err)
				}
			}
			// 映射
			if singleField.Info.CanUseAsMapKey {
				for _, methodText := range []struct {
					attr string
					text string
				}{
					{attrMap, mapByMethodText},         // 列取映射
					{attrSliceMap, sliceMapMethodText}, // 列取数组映射
				} {
					if singleField.Info.Commands.ExistCommandAttr(commandName, parser.Attr(methodText.attr)) {
						if err := list.template.Execute(buf, "List", methodText.text, map[string]interface{}{
							"typName":            structName,
							"typNameWithPath":    typNameWithPath,
							"fieldName":          fieldName,
							"fieldNameWithInner": fieldNameWithInner,
							"fieldType":          fieldTypNameWithPath,
						}); err != nil {
							return pkgName, importPathMap, content, errors.WithStack(err)
						}
					}
				}
			}

			// join
			if singleField.Info.Commands.ExistCommand(joinCommandName) {
				joinTyp,
					joinTypField,
					joinTypWithPath,
					joinTypFieldTyp := singleField.Info.Commands.GetJoinTyp(
					list.parser,
					singleStruct.Info.ImportPath,
					pkg.Structs,
				)

				if err := list.template.Execute(buf, "List", joinMethodText, map[string]interface{}{
					"typName":         structName,
					"typNameWithPath": typNameWithPath,
					"typFieldName":    fieldName,
					"joinTyp":         joinTyp,
					"joinTypWithPath": joinTypWithPath,
					"joinTypField":    joinTypField,
					"joinTypFieldTyp": joinTypFieldTyp,
				}); err != nil {
					return pkgName, importPathMap, content, errors.WithStack(err)
				}
			}

			// derive
			if singleField.Info.Commands.ExistCommand(deriveCommandName) {
				joinTyp,
					joinTypField,
					joinTypWithPath,
					joinTypFieldTyp := singleField.Info.Commands.GetJoinTyp(
					list.parser,
					singleStruct.Info.ImportPath,
					pkg.Structs,
				)

				if err := list.template.Execute(buf, "List", deriveMethodText, map[string]interface{}{
					"typName":           structName,
					"typNameWithPath":   typNameWithPath,
					"typFieldName":      fieldName,
					"joinTyp":           joinTyp,
					"joinTypWithPath":   joinTypWithPath,
					"joinTypField":      joinTypField,
					"joinTypFieldTyp":   joinTypFieldTyp,
					"deriveTypWithPath": extraTypWithPath,
				}); err != nil {
					return pkgName, importPathMap, content, errors.WithStack(err)
				}
			}
		}
	}
	content = buf.Bytes()

	return pkgName, importPathMap, content, nil
}

func (list *List) getFileName(pkgDir, pkgName string) string {
	filename := fmt.Sprintf("%s%s.go", pkgName, GenFileSuffix)
	filename = filepath.Join(pkgDir, filename)
	return filename
}
