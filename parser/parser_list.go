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

func (list FieldList) MapInfo() map[Info]Field {
	result := make(map[Info]Field)
	for _, single := range list {
		result[single.Info] = single
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

func (list InfoList) MapName() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.Name] = single
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

func (list InfoList) MapTypesType() map[types.Type]Info {
	result := make(map[types.Type]Info)
	for _, single := range list {
		result[single.TypesType] = single
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

func (list InfoList) MapType() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.Type] = single
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

func (list InfoList) MapUnderType() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.UnderType] = single
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

func (list InfoList) MapImportPath() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.ImportPath] = single
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

func (list InfoList) MapTypName() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.TypName] = single
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

func (list InfoList) MapTypNameWithPath() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.TypNameWithPath] = single
	}
	return result
}

func (list InfoList) ColumnCanUseAsMapKey() []bool {
	result := make([]bool, len(list), len(list))
	for i, single := range list {
		result[i] = single.CanUseAsMapKey
	}
	return result
}

func (list InfoList) MapCanUseAsMapKey() map[bool]Info {
	result := make(map[bool]Info)
	for _, single := range list {
		result[single.CanUseAsMapKey] = single
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

func (list InfoList) MapComment() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.Comment] = single
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

func (list InfoList) MapDoc() map[string]Info {
	result := make(map[string]Info)
	for _, single := range list {
		result[single.Doc] = single
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

func (list PkgList) MapName() map[string]Pkg {
	result := make(map[string]Pkg)
	for _, single := range list {
		result[single.Name] = single
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

func (list PkgList) MapDir() map[string]Pkg {
	result := make(map[string]Pkg)
	for _, single := range list {
		result[single.Dir] = single
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

func (list PkgList) MapImportPath() map[string]Pkg {
	result := make(map[string]Pkg)
	for _, single := range list {
		result[single.ImportPath] = single
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

func (list StructList) MapInfo() map[Info]Struct {
	result := make(map[Info]Struct)
	for _, single := range list {
		result[single.Info] = single
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
