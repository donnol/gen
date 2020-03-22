package testdata1

import (
	"github.com/donnol/gen/list/testdata1/content"
	"github.com/donnol/gen/list/testdata2"
)

// ModelList 列表结构体
type ModelList []Model

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
func (list ModelList) MapListByName() map[string][]Model {
	result := make(map[string][]Model)
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
func (list ModelList) MapListByAge() map[float64][]Model {
	result := make(map[float64][]Model)
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
