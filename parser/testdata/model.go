package testdata

import (
	"time"

	"github.com/pkg/errors"
)

// Model 模型
type Model struct {
	// 唯一
	ID int // id

	// 长度不限
	Name string // 名称

	// 凡事皆有状态
	Status Status // 状态

	// 默认为false
	Old bool // 是否已旧

	// 默认1.0
	Height float64 // 高度

	// 列表
	Addrs []errors.Frame // 地址

	// 知时知地知人知物
	CreatedAt time.Time // 创建时间

	Err errors.Frame
}

// Status 状态
type Status int
