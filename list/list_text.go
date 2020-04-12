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
			result[i] = single.{{.fieldNameWithInner}}
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
			result[single.{{.fieldNameWithInner}}] = single
		}
		return result
	}
	`

// 列取数组映射方法定义
var sliceMapMethodText = `
	// MapListBy{{.fieldName}} {{.fieldName}}数组映射
	func (list {{.typName}}List) MapListBy{{.fieldName}}() map[{{.fieldType}}][]{{.typNameWithPath}} {
		result := make(map[{{.fieldType}}][]{{.typNameWithPath}})
		for _, single := range list {
			result[single.{{.fieldNameWithInner}}] = append(result[single.{{.fieldNameWithInner}}], single)
		}
		return result
	}
	`

// 连表方法定义
var joinMethodText = `
	// Join{{.joinTyp}}By{{.typFieldName}}Equal{{.joinTypField}} 连表
	func (list {{.typName}}List) Join{{.joinTyp}}By{{.typFieldName}}Equal{{.joinTypField}}(
		ol []{{.joinTypWithPath}},
		f func(
			{{.typNameWithPath}},
			{{.joinTypWithPath}},
		) {{.typNameWithPath}},
	) {{.typName}}List {

		oMap := make(map[{{.joinTypFieldTyp}}]{{.joinTypWithPath}})
		for _, single := range ol {
			oMap[single.{{.joinTypField}}] = single
		}

		result := make({{.typName}}List, len(list), len(list))
		for i, single := range list {
			tmp := f(single, oMap[single.{{.typFieldName}}])

			result[i] = tmp
		}

		return result
	}
	`

// 衍生方法定义
var deriveMethodText = `
	// DeriveBy{{.typFieldName}}Equal{{.joinTypField}} 衍生
	func (list {{.typName}}List) DeriveBy{{.typFieldName}}Equal{{.joinTypField}}(
		ol []{{.joinTypWithPath}},
		f func(
			{{.typNameWithPath}},
			{{.joinTypWithPath}},
		) {{.deriveTypWithPath}},
	) []{{.deriveTypWithPath}} {

		oMap := make(map[{{.joinTypFieldTyp}}]{{.joinTypWithPath}})
		for _, single := range ol {
			oMap[single.{{.joinTypField}}] = single
		}
	
		result := make([]{{.deriveTypWithPath}}, len(list), len(list))
		for i, single := range list {
			tmp := f(single, oMap[single.{{.typFieldName}}])
	
			result[i] = tmp
		}
	
		return result
	}
	`

// 取首个
var firstMethodText = `
	// First 取首个
	func (list {{.typName}}List) First() {{.typName}} {
		if len(list) == 0 {
			panic("Empty list")
		}
		return list[0]
	}
	`
