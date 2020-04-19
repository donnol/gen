package testdata1

import (
	"sort"

	"github.com/donnol/gen/list/testdata1/content"
	"github.com/donnol/gen/list/testdata2"
	"github.com/donnol/gen/list/testdata3"
	"github.com/gofrs/uuid"
)

// ModelList 列表结构体
type ModelList []Model

// Where 返回符合条件的行
func (list ModelList) Where(f func(u Model) bool) ModelList {
	result := make(ModelList, 0, len(list))
	for _, single := range list {
		if !f(single) {
			continue
		}
		result = append(result, single)
	}
	return result
}

// Sort 排序
func (list ModelList) Sort(f func(i, j int) bool) ModelList {
	sort.Slice(list, f)
	return list
}

// Limit 获取从offset位置开始的前几个
func (list ModelList) Limit(offset, n int) ModelList {
	l := len(list)
	result := make(ModelList, 0, l)

	if offset < 0 || offset >= l {
		return result
	}
	if n > l-offset {
		n = l - offset
	}
	for i := offset; i < offset+n; i++ {
		result = append(result, list[i])
	}

	return result
}

// Reduce 降维，从数组变为单个
func (list ModelList) Reduce(f func(u Model, nu Model) Model) Model {
	var u Model
	for i, nu := range list {
		if i == 0 {
			u = nu
			continue
		}
		u = f(u, nu)
	}
	return u
}

// Reverse 反转
func (list ModelList) Reverse() ModelList {
	result := make(ModelList, 0, len(list))
	for i := len(list) - 1; i >= 0; i-- {
		result = append(result, list[i])
	}
	return result
}

// First 取首个
func (list ModelList) First() Model {
	if len(list) == 0 {
		return Model{}
	}
	return list[0]
}

// Last 取最后一个，如果没有数据，会返回结构体零值
func (list ModelList) Last() Model {
	if len(list) == 0 {
		return Model{}
	}
	return list[len(list)-1]
}

// ColumnInnerCode InnerCode列
func (list ModelList) ColumnInnerCode() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Inner.Code
	}
	return result
}

// MapInnerCode InnerCode映射
func (list ModelList) MapInnerCode() map[string]Model {
	result := make(map[string]Model)
	for _, single := range list {
		result[single.Inner.Code] = single
	}
	return result
}

// MapListByInnerCode InnerCode数组映射
func (list ModelList) MapListByInnerCode() map[string]ModelList {
	result := make(map[string]ModelList, len(list))
	for _, single := range list {
		result[single.Inner.Code] = append(result[single.Inner.Code], single)
	}
	return result
}

// ColumnBaseUUID2 BaseUUID2列
func (list ModelList) ColumnBaseUUID2() []uuid.UUID {
	result := make([]uuid.UUID, len(list), len(list))
	for i, single := range list {
		result[i] = single.Base.UUID2
	}
	return result
}

// MapBaseUUID2 BaseUUID2映射
func (list ModelList) MapBaseUUID2() map[uuid.UUID]Model {
	result := make(map[uuid.UUID]Model)
	for _, single := range list {
		result[single.Base.UUID2] = single
	}
	return result
}

// MapListByBaseUUID2 BaseUUID2数组映射
func (list ModelList) MapListByBaseUUID2() map[uuid.UUID]ModelList {
	result := make(map[uuid.UUID]ModelList, len(list))
	for _, single := range list {
		result[single.Base.UUID2] = append(result[single.Base.UUID2], single)
	}
	return result
}

// ColumnID ID列
func (list ModelList) ColumnID() []int {
	result := make([]int, len(list), len(list))
	for i, single := range list {
		result[i] = single.ID
	}
	return result
}

// MapName Name映射
func (list ModelList) MapName() map[string]Model {
	result := make(map[string]Model)
	for _, single := range list {
		result[single.Name] = single
	}
	return result
}

// MapListByName Name数组映射
func (list ModelList) MapListByName() map[string]ModelList {
	result := make(map[string]ModelList, len(list))
	for _, single := range list {
		result[single.Name] = append(result[single.Name], single)
	}
	return result
}

// MapAge Age映射
func (list ModelList) MapAge() map[float64]Model {
	result := make(map[float64]Model)
	for _, single := range list {
		result[single.Age] = single
	}
	return result
}

// MapListByAge Age数组映射
func (list ModelList) MapListByAge() map[float64]ModelList {
	result := make(map[float64]ModelList, len(list))
	for _, single := range list {
		result[single.Age] = append(result[single.Age], single)
	}
	return result
}

// JoinUserByUserIDEqualID 连表
func (list ModelList) JoinUserByUserIDEqualID(
	ol []User,
	f func(
		Model,
		User,
	) Model,
) ModelList {

	oMap := make(map[int]User)
	for _, single := range ol {
		oMap[single.ID] = single
	}

	result := make(ModelList, len(list), len(list))
	for i, single := range list {
		tmp := f(single, oMap[single.UserID])

		result[i] = tmp
	}

	return result
}

// DeriveByUserIDEqualID 衍生
func (list ModelList) DeriveByUserIDEqualID(
	ol []User,
	f func(
		Model,
		User,
	) ModelUser,
) []ModelUser {

	oMap := make(map[int]User)
	for _, single := range ol {
		oMap[single.ID] = single
	}

	result := make([]ModelUser, len(list), len(list))
	for i, single := range list {
		tmp := f(single, oMap[single.UserID])

		result[i] = tmp
	}

	return result
}

// JoinContentByContentIDEqualID 连表
func (list ModelList) JoinContentByContentIDEqualID(
	ol []content.Content,
	f func(
		Model,
		content.Content,
	) Model,
) ModelList {

	oMap := make(map[int]content.Content)
	for _, single := range ol {
		oMap[single.ID] = single
	}

	result := make(ModelList, len(list), len(list))
	for i, single := range list {
		tmp := f(single, oMap[single.ContentID])

		result[i] = tmp
	}

	return result
}

// JoinModelByModelIDEqualID 连表
func (list ModelList) JoinModelByModelIDEqualID(
	ol []testdata2.Model,
	f func(
		Model,
		testdata2.Model,
	) Model,
) ModelList {

	oMap := make(map[int]testdata2.Model)
	for _, single := range ol {
		oMap[single.ID] = single
	}

	result := make(ModelList, len(list), len(list))
	for i, single := range list {
		tmp := f(single, oMap[single.ModelID])

		result[i] = tmp
	}

	return result
}

// JoinAddrByAddrIDEqualID 连表
func (list ModelList) JoinAddrByAddrIDEqualID(
	ol []testdata3.Addr,
	f func(
		Model,
		testdata3.Addr,
	) Model,
) ModelList {

	oMap := make(map[int]testdata3.Addr)
	for _, single := range ol {
		oMap[single.ID] = single
	}

	result := make(ModelList, len(list), len(list))
	for i, single := range list {
		tmp := f(single, oMap[single.AddrID])

		result[i] = tmp
	}

	return result
}
