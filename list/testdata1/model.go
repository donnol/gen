package testdata1

import (
	"go/types"

	"github.com/donnol/gen/list/testdata1/base"
	"github.com/pkg/errors"
)

// Model 模型
//
// @gen list
type Model struct {
	// @gen list [column, map, slicemap] .Code
	Inner

	// @gen list [column, map, slicemap] .UUID
	base.Base

	// @gen list column
	// 唯一
	ID int // id

	// 长度不限
	// @gen list map
	// @gen list slicemap
	Name string // 名称

	// @gen list [map, slicemap]
	Age float64 // 年龄

	// @gen join =User.ID
	// @gen derive =User.ID ModelUser
	UserID   int
	UserName string

	// @gen join =./content.Content.ID
	ContentID    int
	ContentTitle string

	// @gen join =../testdata2.Model.ID
	ModelID   int
	ModelName string

	// @gen join =github.com/donnol/gen/list/testdata3.Addr.ID
	AddrID   int
	AddrName string

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

// Inner 内嵌
type Inner struct {
	Code string
}

// User 用户
type User struct {
	ID   int
	Name string
}

// ModelUser 模型用户
type ModelUser struct {
	ModelID   int
	ModelName string
	UserID    int
	UserName  string
}
