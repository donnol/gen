package testdata

type ModelList []Model

func (list ModelList) ColumnID() []int {
	result := make([]int, len(list), len(list))
	for i, single := range list {
		result[i] = single.ID
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

func (list ModelList) MapAge() map[float64]Model {
	result := make(map[float64]Model)
	for _, single := range list {
		result[single.Age] = single
	}
	return result
}

func (list ModelList) MapListByAge() map[float64][]Model {
	result := make(map[float64][]Model)
	for _, single := range list {
		result[single.Age] = append(result[single.Age], single)
	}
	return result
}
