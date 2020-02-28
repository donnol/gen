package testdata

import (
	"go/types"

	"github.com/pkg/errors"
)

type ModelList []Model

func (list ModelList) ColumnID() []int {
	result := make([]int, len(list), len(list))
	for i, single := range list {
		result[i] = single.ID
	}
	return result
}

func (list ModelList) MapID() map[int]Model {
	result := make(map[int]Model)
	for _, single := range list {
		result[single.ID] = single
	}
	return result
}

func (list ModelList) MapListByID() map[int][]Model {
	result := make(map[int][]Model)
	for _, single := range list {
		result[single.ID] = append(result[single.ID], single)
	}
	return result
}

func (list ModelList) ColumnName() []string {
	result := make([]string, len(list), len(list))
	for i, single := range list {
		result[i] = single.Name
	}
	return result
}

func (list ModelList) MapName() map[string]Model {
	result := make(map[string]Model)
	for _, single := range list {
		result[single.Name] = single
	}
	return result
}

func (list ModelList) MapListByName() map[string][]Model {
	result := make(map[string][]Model)
	for _, single := range list {
		result[single.Name] = append(result[single.Name], single)
	}
	return result
}

func (list ModelList) ColumnArray() [][4]int {
	result := make([][4]int, len(list), len(list))
	for i, single := range list {
		result[i] = single.Array
	}
	return result
}

func (list ModelList) MapArray() map[[4]int]Model {
	result := make(map[[4]int]Model)
	for _, single := range list {
		result[single.Array] = single
	}
	return result
}

func (list ModelList) MapListByArray() map[[4]int][]Model {
	result := make(map[[4]int][]Model)
	for _, single := range list {
		result[single.Array] = append(result[single.Array], single)
	}
	return result
}

func (list ModelList) ColumnChan() []chan int {
	result := make([]chan int, len(list), len(list))
	for i, single := range list {
		result[i] = single.Chan
	}
	return result
}

func (list ModelList) MapChan() map[chan int]Model {
	result := make(map[chan int]Model)
	for _, single := range list {
		result[single.Chan] = single
	}
	return result
}

func (list ModelList) MapListByChan() map[chan int][]Model {
	result := make(map[chan int][]Model)
	for _, single := range list {
		result[single.Chan] = append(result[single.Chan], single)
	}
	return result
}

func (list ModelList) ColumnMap() []map[int]int {
	result := make([]map[int]int, len(list), len(list))
	for i, single := range list {
		result[i] = single.Map
	}
	return result
}

func (list ModelList) ColumnPointer() []*int {
	result := make([]*int, len(list), len(list))
	for i, single := range list {
		result[i] = single.Pointer
	}
	return result
}

func (list ModelList) MapPointer() map[*int]Model {
	result := make(map[*int]Model)
	for _, single := range list {
		result[single.Pointer] = single
	}
	return result
}

func (list ModelList) MapListByPointer() map[*int][]Model {
	result := make(map[*int][]Model)
	for _, single := range list {
		result[single.Pointer] = append(result[single.Pointer], single)
	}
	return result
}

func (list ModelList) ColumnSlice() [][]int {
	result := make([][]int, len(list), len(list))
	for i, single := range list {
		result[i] = single.Slice
	}
	return result
}

func (list ModelList) ColumnOutArray() [][4]errors.Frame {
	result := make([][4]errors.Frame, len(list), len(list))
	for i, single := range list {
		result[i] = single.OutArray
	}
	return result
}

func (list ModelList) MapOutArray() map[[4]errors.Frame]Model {
	result := make(map[[4]errors.Frame]Model)
	for _, single := range list {
		result[single.OutArray] = single
	}
	return result
}

func (list ModelList) MapListByOutArray() map[[4]errors.Frame][]Model {
	result := make(map[[4]errors.Frame][]Model)
	for _, single := range list {
		result[single.OutArray] = append(result[single.OutArray], single)
	}
	return result
}

func (list ModelList) ColumnOutChan() []chan errors.Frame {
	result := make([]chan errors.Frame, len(list), len(list))
	for i, single := range list {
		result[i] = single.OutChan
	}
	return result
}

func (list ModelList) MapOutChan() map[chan errors.Frame]Model {
	result := make(map[chan errors.Frame]Model)
	for _, single := range list {
		result[single.OutChan] = single
	}
	return result
}

func (list ModelList) MapListByOutChan() map[chan errors.Frame][]Model {
	result := make(map[chan errors.Frame][]Model)
	for _, single := range list {
		result[single.OutChan] = append(result[single.OutChan], single)
	}
	return result
}

func (list ModelList) ColumnOutMap() []map[int]errors.Frame {
	result := make([]map[int]errors.Frame, len(list), len(list))
	for i, single := range list {
		result[i] = single.OutMap
	}
	return result
}

func (list ModelList) ColumnOutMap2() []map[errors.Frame]errors.Frame {
	result := make([]map[errors.Frame]errors.Frame, len(list), len(list))
	for i, single := range list {
		result[i] = single.OutMap2
	}
	return result
}

func (list ModelList) ColumnOutMap3() []map[types.Type]errors.Frame {
	result := make([]map[types.Type]errors.Frame, len(list), len(list))
	for i, single := range list {
		result[i] = single.OutMap3
	}
	return result
}

func (list ModelList) ColumnOutPointer() []*errors.Frame {
	result := make([]*errors.Frame, len(list), len(list))
	for i, single := range list {
		result[i] = single.OutPointer
	}
	return result
}

func (list ModelList) MapOutPointer() map[*errors.Frame]Model {
	result := make(map[*errors.Frame]Model)
	for _, single := range list {
		result[single.OutPointer] = single
	}
	return result
}

func (list ModelList) MapListByOutPointer() map[*errors.Frame][]Model {
	result := make(map[*errors.Frame][]Model)
	for _, single := range list {
		result[single.OutPointer] = append(result[single.OutPointer], single)
	}
	return result
}

func (list ModelList) ColumnOutSlice() [][]errors.Frame {
	result := make([][]errors.Frame, len(list), len(list))
	for i, single := range list {
		result[i] = single.OutSlice
	}
	return result
}
