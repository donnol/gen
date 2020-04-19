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
	func (list {{.typName}}List) MapListBy{{.fieldName}}() map[{{.fieldType}}]{{.typName}}List {
		result := make(map[{{.fieldType}}]{{.typName}}List, len(list))
		for _, single := range list {
			result[single.{{.fieldNameWithInner}}] = append(result[single.{{.fieldNameWithInner}}], single)
		}
		return result
	}
	`

// ===
// === 以下方法不需要指定字段，生成列表结构体时即可同时生成
// ===

// where, 返回符合条件的行
var whereMethodText = `
	// Where 返回符合条件的行
	func (list {{.typName}}List) Where(f func(u {{.typName}}) bool) {{.typName}}List {
		result := make({{.typName}}List, 0, len(list))
		for _, single := range list {
			if !f(single) {
				continue
			}
			result = append(result, single)
		}
		return result
	}
	`

// len, 长度
var lenMethodText = `
	// Len 长度
	func (list {{.typName}}List) Len() int {
		return len(list)
	}
	`

// sort, 排序
var sortMethodText = `
	// Sort 排序
	func (list {{.typName}}List) Sort(f func(i, j int) bool) {{.typName}}List {
		sort.Slice(list, f)
		return list
	}
	`

// limit, 从offset位置开始的前几个
var limitMethodText = `
	// Limit 获取从offset位置开始的前几个
	func (list {{.typName}}List) Limit(offset, n int) {{.typName}}List {
		l := len(list)
		result := make({{.typName}}List, 0, l)

		if offset < 0 || offset >= l {
			return result
		}
		if n > l - offset {
			n = l - offset
		}
		for i:=offset; i<offset+n; i++ {
			result = append(result, list[i])
		}

		return result
	}
	`

// reduce, 降维，从数组变为单个，对数组中的每个元素执行函数(升序执行)，将其结果汇总为单个返回值
var reduceMethodText = `
	// Reduce 降维，从数组变为单个
	func (list {{.typName}}List) Reduce(f func(u {{.typName}}, nu {{.typName}}) {{.typName}}) {{.typName}} {
		var u {{.typName}}
		for i, nu := range list {
			if i == 0 {
				u = nu
				continue
			}
			u = f(u, nu)
		}
		return u
	}
	`

// reverse, 反转
var reverseMethodText = `
	// Reverse 反转
	func (list {{.typName}}List) Reverse() {{.typName}}List {
		result := make({{.typName}}List, 0, len(list))
		for i := len(list)-1; i>=0 ; i-- {
			result = append(result, list[i])
		}
		return result
	}
	`

// first, 取首个，如果没有数据，会返回结构体零值
var firstMethodText = `
	// First 取首个
	func (list {{.typName}}List) First() {{.typName}} {
		if len(list) == 0 {
			return {{.typName}}{}
		}
		return list[0]
	}
	`

// last, 取最后一个，如果没有数据，会返回结构体零值
var lastMethodText = `
	// Last 取最后一个，如果没有数据，会返回结构体零值
	func (list {{.typName}}List) Last() {{.typName}} {
		if len(list) == 0 {
			return {{.typName}}{}
		}
		return list[len(list)-1]
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
