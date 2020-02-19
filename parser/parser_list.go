package parser

import (
	"go/types"
)

// FieldList 列表结构体
type FieldList []Field

// ColumnInfo Info列
func (list FieldList) ColumnInfo() []Info {
	result := make([]Info, len(list), len(list))
	for i, single := range list {
		result[i] = single.Info
	}
	return result
}

// InfoList 列表结构体
type InfoList []Info

// ColumnName Name列
func (list InfoList) ColumnName() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Name
	}
	return result
}

// ColumnTypesType TypesType列
func (list InfoList) ColumnTypesType() []types.Type {
	result := make([]types.Type, len(list), len(list))
	for i, single := range list {
		result[i] = single.TypesType
	}
	return result
}

// ColumnType Type列
func (list InfoList) ColumnType() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Type
	}
	return result
}

// ColumnUnderType UnderType列
func (list InfoList) ColumnUnderType() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.UnderType
	}
	return result
}

// ColumnImportPath ImportPath列
func (list InfoList) ColumnImportPath() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.ImportPath
	}
	return result
}

// ColumnTypName TypName列
func (list InfoList) ColumnTypName() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.TypName
	}
	return result
}

// ColumnTypNameWithPath TypNameWithPath列
func (list InfoList) ColumnTypNameWithPath() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.TypNameWithPath
	}
	return result
}

// ColumnComment Comment列
func (list InfoList) ColumnComment() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Comment
	}
	return result
}

// ColumnDoc Doc列
func (list InfoList) ColumnDoc() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Doc
	}
	return result
}

// OptionList 列表结构体
type OptionList []Option

// ParserList 列表结构体
type ParserList []Parser

// PkgList 列表结构体
type PkgList []Pkg

// ColumnName Name列
func (list PkgList) ColumnName() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Name
	}
	return result
}

// ColumnDir Dir列
func (list PkgList) ColumnDir() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Dir
	}
	return result
}

// ColumnImportPath ImportPath列
func (list PkgList) ColumnImportPath() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.ImportPath
	}
	return result
}

// ColumnStructs Structs列
func (list PkgList) ColumnStructs() [][]Struct {
	result := make([][]Struct, len(list), len(list))
	for i, single := range list {
		result[i] = single.Structs
	}
	return result
}

// StructList 列表结构体
type StructList []Struct

// ColumnInfo Info列
func (list StructList) ColumnInfo() []Info {
	result := make([]Info, len(list), len(list))
	for i, single := range list {
		result[i] = single.Info
	}
	return result
}

// ColumnFields Fields列
func (list StructList) ColumnFields() [][]Field {
	result := make([][]Field, len(list), len(list))
	for i, single := range list {
		result[i] = single.Fields
	}
	return result
}
