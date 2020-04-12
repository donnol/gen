package parser

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Pkg 包
type Pkg struct {
	Name       string // 包名
	Dir        string // 包目录
	ImportPath string // 包导入路径
	Structs    []Struct
}

// Struct 结构体
type Struct struct {
	Info
	Fields []Field
}

// StructList 列表
type StructList []Struct

// Find 找出指定名称的结构体
func (list StructList) Find(name string) Struct {
	for _, single := range list {
		if single.Info.Name == name {
			return single
		}
	}
	return Struct{}
}

// FindField 找出指定结构体的指定字段
func (list StructList) FindField(name, fieldName string) Field {
	for _, single := range list {
		if single.Info.Name == name {
			for _, field := range single.Fields {
				if field.Info.Name == fieldName {
					return field
				}
			}
		}
	}

	return Field{}
}

// Field 字段
type Field struct {
	Info
}

// Info 信息
type Info struct {
	Name string // 名称，如：Info

	Anonymous bool // 是否匿名

	// 可以使用这个字段来断言到具体的类型，如go/types.Struct
	TypesType types.Type // go/types.Type接口

	// 如果是导入的第三方包，会包含包路径，如：github.com/pkg/errors.fundamental
	Type string // 类型，表面类型

	// 如type Status int的底下类型就是int
	UnderType string // 底下类型

	ImportPath      string // 导入路径
	TypName         string // 类型名，如果是结构体则与Name相同，如果是字段，则是字段的类型
	TypNameWithPath string // 带有导入路径的类型名称，如：parser.Info
	CanUseAsMapKey  bool   // 字段类型是否可作为map的键

	Comment string // 注释
	Doc     string // 文档

	Commands CommandList // 指令，如list, list column, list map, list slicemap等中的一个或多个
}

// Command 指令
type Command struct {
	Name  string // 名称，如：list
	Attr  Attr   // 属性，如: column
	Extra string // 额外，如: deriveType
}

// CommandList 列表
type CommandList []Command

// ExistCommand 存在指令
func (list CommandList) ExistCommand(name string) bool {
	for _, single := range list {
		if single.Name == name {
			return true
		}
	}
	return false
}

// ExistCommandAttr 存在指令
func (list CommandList) ExistCommandAttr(name string, attr Attr) bool {
	for _, single := range list {
		if single.Name == name && single.Attr == attr {
			return true
		}
	}
	return false
}

// GetJoinTyp 获取join类型信息
func (list CommandList) GetJoinTyp(parser *Parser, pkgPath string, structs []Struct) (joinTyp, joinTypField, joinTypWithPath, joinTypFieldTyp string) {
	for _, single := range list {
		return single.Attr.GetJoinTyp(parser, pkgPath, structs)
	}

	return
}

// GetExtraTyp 获取额外信息
func (list CommandList) GetExtraTyp(parser *Parser) (extraTyp, extraTypWithPath string) {
	for _, single := range list {
		if single.Extra != "" {
			// TODO:实现带路径类型解析
			extraTyp = single.Extra
			extraTypWithPath = single.Extra
			return
		}
	}

	return
}

// Attr 属性
type Attr string

// GetJoinTyp 获取join类型信息
func (attr Attr) GetJoinTyp(parser *Parser, pkgPath string, structs []Struct) (joinTyp, joinTypField, joinTypWithPath, joinTypFieldTyp string) {
	if strings.TrimSpace(string(attr)) == "" {
		panic(errors.Errorf("Empty attr"))
	}
	typAndField := string(attr[1:])

	// 解析结构体和字段
	parts := strings.Split(typAndField, ".")
	if len(parts) < 2 {
		panic(errors.Errorf("Bad join attr: %s", attr))
	}

	if len(parts) == 2 {
		// User.ID
		joinTyp = parts[0]
		joinTypField = parts[1]
		joinTypWithPath = joinTyp
		field := StructList(structs).FindField(joinTyp, joinTypField)
		joinTypFieldTyp = field.Type

		return
	}

	// 包含路径
	var fieldLeft, joinTypPkgPath string

	if strings.Index(typAndField, ".") == 0 {
		// 相对路径：
		// 去掉第一个斜杆前的内容
		firstSlashIndex := strings.Index(typAndField, "/")
		fieldLeft = typAndField[firstSlashIndex+1:]
	} else {
		// github.com/pkg/errors.XXX.YYY，
		fieldLeft = typAndField
	}

	// 结构体和字段
	lastDotIndex := strings.LastIndex(fieldLeft, ".")
	joinTypField = fieldLeft[lastDotIndex+1:]
	fieldLeft = fieldLeft[:lastDotIndex]
	lastDotIndex = strings.LastIndex(fieldLeft, ".")
	joinTyp = fieldLeft[lastDotIndex+1:]
	fieldLeft = fieldLeft[:lastDotIndex]
	lastSlashIndex := strings.LastIndex(fieldLeft, "/")
	if lastSlashIndex != -1 {
		joinTypWithPath = fieldLeft[lastSlashIndex+1:] + "." + joinTyp
	} else {
		joinTypWithPath = fieldLeft + "." + joinTyp
	}

	// 找到包路径
	if strings.Index(typAndField, "..") == 0 {
		// ../pkgpath.XXX.YYY
		joinTypPkgPath = filepath.Clean(filepath.Join(pkgPath, "../", fieldLeft))
	} else if strings.Index(typAndField, ".") == 0 {
		// ./pkgpath.XXX.YYY
		joinTypPkgPath = filepath.Join(pkgPath, fieldLeft)
	} else {
		// github.com/pkg/errors.XXX.YYY，
		joinTypPkgPath = fieldLeft
	}

	// 字段类型
	pkg, err := parser.ParsePkg(joinTypPkgPath)
	if err != nil {
		panic(err)
	}
	field := StructList(pkg.Structs).FindField(joinTyp, joinTypField)
	joinTypFieldTyp = field.Type

	return
}

// InitWithTypes 初始化
func (info *Info) InitWithTypes(name, pkgPath, doc, comment string, typ types.Type) {
	info.Name = name
	info.TypesType = typ
	info.Type = typ.String()
	info.UnderType = typ.Underlying().String()
	info.ImportPath = pkgPath
	info.Doc = doc
	info.Comment = comment
	info.CanUseAsMapKey = types.Comparable(info.TypesType)

	// 解析指令
	info.initCommandsFromDoc()

	// 路径等信息
	importPath, typName, typNameWithPath := info.GetImportPathAndTypeName(typ.String())
	if info.ImportPath == "" {
		info.ImportPath = importPath
	}
	info.TypName = typName
	info.TypNameWithPath = typNameWithPath
}

func (info *Info) initCommandsFromDoc() {
	// 读取文档行
	lines := strings.Split(info.Doc, "\n")
	commands := make([]Command, 0)
	for _, line := range lines {
		if !strings.Contains(line, "@gen") {
			continue
		}

		cmds := parseGenCommand(line)
		commands = append(commands, cmds...)
	}
	info.Commands = commands
}

func parseGenCommand(line string) (cmds []Command) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		fmt.Printf("Wrong command: %s\n", line)
		return
	}
	// 如果有多个，需要拆成独立的
	leftIndex := strings.Index(line, "[")
	rightIndex := strings.Index(line, "]")
	if leftIndex != -1 && rightIndex != -1 {
		attrs := line[leftIndex+1 : rightIndex]
		attrParts := strings.Split(attrs, ",")

		var extra string
		if rightIndex+2 < len(line) {
			extra = line[rightIndex+2:]
		}
		for _, p := range attrParts {
			p := strings.TrimSpace(p)
			cmds = append(cmds, Command{
				Name:  parts[1],
				Attr:  Attr(p),
				Extra: extra,
			})
		}
	} else {
		var attr, extra string
		if len(parts) > 2 {
			attr = parts[2]
		}
		if len(parts) > 3 {
			extra = parts[3]
		}
		cmds = append(cmds, Command{
			Name:  parts[1],
			Attr:  Attr(attr),
			Extra: extra,
		})
	}

	return
}

// GetNotEmptyImportPathMap 获取非空导入路径映射
func (info Info) GetNotEmptyImportPathMap() ImportPathMap {
	importPathMap := make(ImportPathMap)
	importPath := info.ImportPath
	if importPath != "" {
		importPathMap[importPath] = append(importPathMap[importPath], info.TypNameWithPath)
	}
	return importPathMap
}

// GetTypNameWithPath 获取带路径类型名称
func (info Info) GetTypNameWithPath(pkgImportPath string) string {
	typNameWithPath := info.TypNameWithPath
	if info.ImportPath == pkgImportPath {
		typNameWithPath = info.TypName
	}
	return typNameWithPath
}

// GetTypesTypeStructField 获取结构体字段信息
func (info Info) GetTypesTypeStructField(fieldName string) (
	typName string, typNameWithPath string,
) {
	v, ok := info.TypesType.Underlying().(*types.Struct)
	if !ok {
		return
	}
	for i := 0; i < v.NumFields(); i++ {
		field := v.Field(i)
		if field.Name() != fieldName {
			continue
		}
		typName = field.Type().String()
		typNameWithPath = field.Type().String()
	}

	return
}

// GetImportPathAndTypeName 获取导入路径和类型名称
func (info *Info) GetImportPathAndTypeName(full string) (
	string,
	string,
	string,
) {
	importPath, typeName, typNameWithPath := getImportPathAndTypeName(full)

	switch v := info.TypesType.(type) {
	case *types.Array:
		debug("Array: %+v, %+v, %d\n", v, v.Elem(), v.Len())

		elemTypName := v.Elem().String()
		importPath, typeName, typNameWithPath = getImportPathAndTypeName(elemTypName)
		prefix := "[" + strconv.Itoa(int(v.Len())) + "]"
		typeName = prefix + typeName
		typNameWithPath = prefix + typNameWithPath

	case *types.Basic:
		debug("Basic: %+v, %+v\n", v, v.Info())

	case *types.Chan:
		debug("Chan: %+v, %+v\n", v, v.Elem())

		elemTypName := v.Elem().String()
		importPath, typeName, typNameWithPath = getImportPathAndTypeName(elemTypName)
		prefix := "chan "
		typeName = prefix + typeName
		typNameWithPath = prefix + typNameWithPath

	case *types.Interface:
		debug("Interface: %+v, %+v\n", v, v.NumMethods())

	case *types.Map:
		debug("Map: %+v, %+v\n", v, v.Elem())

		elemTypName := v.Elem().String()
		importPath, typeName, typNameWithPath = getImportPathAndTypeName(elemTypName)
		_, _, keyTypNameWithPath := getImportPathAndTypeName(v.Key().String())
		prefix := "map[" + keyTypNameWithPath + "] "
		typeName = prefix + typeName
		typNameWithPath = prefix + typNameWithPath

	case *types.Named:
		debug("Named: %+v, %+v\n", v, v.NumMethods())

	case *types.Pointer:
		debug("Pointer: %+v, %+v\n", v, v.Elem())

		elemTypName := v.Elem().String()
		importPath, typeName, typNameWithPath = getImportPathAndTypeName(elemTypName)
		prefix := "*"
		typeName = prefix + typeName
		typNameWithPath = prefix + typNameWithPath

	case *types.Signature:
		debug("Signature: %+v, %+v\n", v, v.Params())

	case *types.Slice:
		debug("Slice: %+v, %+v\n", v, v.Elem())

		elemTypName := v.Elem().String()
		importPath, typeName, typNameWithPath = getImportPathAndTypeName(elemTypName)
		prefix := "[]"
		typeName = prefix + typeName
		typNameWithPath = prefix + typNameWithPath

	case *types.Struct:
		debug("Struct: %+v, %+v\n", v, v.NumFields())

	case *types.Tuple:
		debug("Tuple: %+v, %+v\n", v, v.Len())

	}

	return importPath, typeName, typNameWithPath
}

func getImportPathAndTypeName(full string) (string, string, string) {
	importPath := ""
	typeName := full
	typNameWithPath := full

	// 带有包导入路径
	typLastIndex := strings.LastIndex(full, ".")
	if typLastIndex != -1 {
		// 类型名
		typeName = full[typLastIndex+1:]

		// 包名+类型名
		slashLastIndex := strings.LastIndex(full, "/")
		if slashLastIndex != -1 {
			typNameWithPath = typNameWithPath[slashLastIndex+1:]
		}

		// 包路径
		importPath = full[:typLastIndex]
	}

	return importPath, typeName, typNameWithPath
}

// ImportPathMap 路径映射
type ImportPathMap map[string][]string

// Keys 键
func (m ImportPathMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// MergeKey 合并，如果有重复键，om里的值将会覆盖m的
func (m ImportPathMap) MergeKey(om ImportPathMap) ImportPathMap {
	nm := make(ImportPathMap)

	for k, v := range m {
		nm[k] = v
	}
	for k, v := range om {
		nm[k] = v
	}
	return nm
}

func debug(format string, args ...interface{}) {
	if v := os.Getenv("debug"); v != "" {
		fmt.Printf(format, args...)
	}
}
