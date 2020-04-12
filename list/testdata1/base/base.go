package base

import "github.com/gofrs/uuid"

// Base 基底
type Base struct {
	UUID string

	// 包含其它包的类型
	UUID2 uuid.UUID
}
