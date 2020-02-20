package parser

import (
	"go/types"
)

type FieldList []Field

func (list FieldList) ColumnInfo() []Info {
	result := make([]Info, len(list), len(list))
	for i, single := range list {
		result[i] = single.Info
	}
	return result
}

type InfoList []Info

func (list InfoList) ColumnName() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Name
	}
	return result
}

func (list InfoList) ColumnTypesType() []types.Type {
	result := make([]types.Type, len(list), len(list))
	for i, single := range list {
		result[i] = single.TypesType
	}
	return result
}

func (list InfoList) ColumnType() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Type
	}
	return result
}

func (list InfoList) ColumnUnderType() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.UnderType
	}
	return result
}

func (list InfoList) ColumnImportPath() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.ImportPath
	}
	return result
}

func (list InfoList) ColumnTypName() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.TypName
	}
	return result
}

func (list InfoList) ColumnTypNameWithPath() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.TypNameWithPath
	}
	return result
}

func (list InfoList) ColumnComment() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Comment
	}
	return result
}

func (list InfoList) ColumnDoc() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Doc
	}
	return result
}

type OptionList []Option

type ParserList []Parser

type PkgList []Pkg

func (list PkgList) ColumnName() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Name
	}
	return result
}

func (list PkgList) ColumnDir() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Dir
	}
	return result
}

func (list PkgList) ColumnImportPath() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.ImportPath
	}
	return result
}

func (list PkgList) ColumnStructs() [][]Struct {
	result := make([][]Struct, len(list), len(list))
	for i, single := range list {
		result[i] = single.Structs
	}
	return result
}

type StructList []Struct

func (list StructList) ColumnInfo() []Info {
	result := make([]Info, len(list), len(list))
	for i, single := range list {
		result[i] = single.Info
	}
	return result
}

func (list StructList) ColumnFields() [][]Field {
	result := make([][]Field, len(list), len(list))
	for i, single := range list {
		result[i] = single.Fields
	}
	return result
}
