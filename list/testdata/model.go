package testdata

import (
	"go/types"

	"github.com/pkg/errors"
)

// Model 模型
type Model struct {
	// 唯一
	ID int // id

	// 长度不限
	Name string // 名称

	// Array，Chan，Map，Pointer，Slice
	Array   [4]int
	Chan    chan int
	Map     map[int]int
	Pointer *int
	Slice   []int

	// 来自三方库
	OutArray   [4]errors.Frame
	OutChan    chan errors.Frame
	OutMap     map[int]errors.Frame
	OutMap2    map[errors.Frame]errors.Frame
	OutMap3    map[types.Type]errors.Frame
	OutPointer *errors.Frame
	OutSlice   []errors.Frame
}
