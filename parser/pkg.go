package parser

import "go/types"

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

	Comment string // 注释
	Doc     string // 文档
}
