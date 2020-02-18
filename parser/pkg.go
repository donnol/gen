package parser

// Pkg 包
type Pkg struct {
	Structs []Struct
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
	Name      string // 名称
	Type      string // 类型，表面类型
	UnderType string // 底下类型
	Comment   string // 注释
	Doc       string // 文档
}
