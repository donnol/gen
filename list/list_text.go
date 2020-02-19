package list

var typText = `
	// {{.typName}}List 列表结构体
	type {{.typName}}List []{{.typNameWithPath}}
`

var columnMethodText = `
	// Column{{.fieldName}} {{.fieldName}}列
	func (list {{.typName}}List) Column{{.fieldName}}() []{{.fieldType}} {
		result := make([]{{.fieldType}}, len(list), len(list))
		for i, single := range list {
			result[i] = single.{{.fieldName}}
		}
		return result
	}
	`
