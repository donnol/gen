package parser

import (
	"fmt"
	"go/types"
	"os"
	"strconv"
	"strings"
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

// Field 字段
type Field struct {
	Info
}

// Info 信息
type Info struct {
	Name string // 名称，如：Info

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
	Name string // 名称，如：list
	Attr string // 属性，如: column
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
func (list CommandList) ExistCommandAttr(name string, attr string) bool {
	for _, single := range list {
		if single.Name == name && single.Attr == attr {
			return true
		}
	}
	return false
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
		for _, p := range attrParts {
			p := strings.TrimSpace(p)
			cmds = append(cmds, Command{
				Name: parts[1],
				Attr: p,
			})
		}
	} else {
		var attr string
		if len(parts) > 2 {
			attr = parts[2]
		}
		cmds = append(cmds, Command{
			Name: parts[1],
			Attr: attr,
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
