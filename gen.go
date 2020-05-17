// Gen 提供了方法，让用户可以根据自己需要来决定生成哪些方法
// 如：
// 有结构体
//  type User struct {
// 		ID int
// 		Name string
// }
// gen.Register(&User{}, gen.Option{})
package gen

// Type 类型
type Type struct{}

// Register 注册类型，将实例传入，以供后续生成时使用
func Register(t interface{}, opt Option) Type {
	return Type{}
}

// Option 选项，说明方法和字段
type Option struct {
	JoinType interface{} // 供Join使用，如果为nil则使用自身，需要是列表结构体，没有的话请先生成
	JoinCond Cond        // join条件

	TatgetType interface{} // 表示目标类型，供Map, Reduce, Join使用，如果为nil则使用自身，需要是列表结构体，没有的话请先生成

	Methods      []Method // 与field无关方法，为nil表示全部
	FieldMethods []Method // 与field有关方法，必须指定方法和字段
}

// Cond 条件
type Cond string

// Method 方法
type Method struct {
	Name   MethodName
	Fields []string // 字段名，与具体结构体有关
}

// MethodName 方法名
type MethodName string

// 方法名枚举
const (
	Where MethodName = "Where"
	// ...
)

// 方法枚举
var ()
