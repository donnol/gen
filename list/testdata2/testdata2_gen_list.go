package testdata2

import (
	"math/rand"
	"sort"
	"time"
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

// Len 长度
func (list ModelList) Len() int {
	return len(list)
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

// Map 对列表里的每个元素执行指定操作
func (list ModelList) Map(f func(u Model) Model) ModelList {
	r := make(ModelList, 0, len(list))
	for _, single := range list {
		r = append(r, f(single))
	}
	return r
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

// Each 逐个元素遍历
func (list ModelList) Each(f func(u Model, i int)) {
	for i, single := range list {
		f(single, i)
	}
}

// Shuffle 洗牌
func (list ModelList) Shuffle() ModelList {
	r := make(ModelList, len(list))
	copy(r, list)
	ras := rand.NewSource(time.Now().Unix())
	ra := rand.New(ras)
	ra.Shuffle(len(r), func(i, j int) {
		r[i], r[j] = r[j], r[i]
	})
	return r
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
