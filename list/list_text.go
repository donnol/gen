package list

// 类型定义，需要指定新的类型名，以及已有类型（如果是其它包的，还需要带上包名，如errors.Frame）
var typText = `
	// {{.typName}}List 列表结构体
	type {{.typName}}List []{{.typNameWithPath}}
`

// 取列方法定义，需要指定类型名，字段名，字段类型
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

// 列取映射方法定义
var mapMethodText = `
	// Map{{.fieldName}} {{.fieldName}}映射
	func (list {{.typName}}List) Map{{.fieldName}}() map[{{.fieldType}}]{{.typNameWithPath}} {
		result := make(map[{{.fieldType}}]{{.typNameWithPath}})
		for _, single := range list {
			result[single.{{.fieldName}}] = single
		}
		return result
	}
	`
