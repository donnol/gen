# gen

[![Go Report Card](https://goreportcard.com/badge/github.com/donnol/gen)](https://goreportcard.com/report/github.com/donnol/gen)
[![GoDoc](https://pkg.go.dev/mod/github.com/donnol/gen?status.svg)](https://pkg.go.dev/mod/github.com/donnol/gen)

生成。

## 为什么

一切都是为了复用。

并在达到复用目标的同时，还要保证代码的显式、清晰。

## 结构体列表

给定一个结构体，生成相应的列表结构体，并在列表结构体上生成方法。方法包括取列，列映射，列数组映射等

如：

给定结构体：

```go
type Model struct {
    ID int
    Name string
    Old bool
    Height float64
    CreatedAt time.Time
}
```

生成以下代码，文件名为 xxx_list.go：

```go
type ModelList []Model

func (list ModelList) ColumnID() []int {
    result := make([]int, 0, len(list))
    for _, single := range list {
        result = append(result, single.ID)
    }
    return result
}

// 注意，这个写法，如果列表里存在多个相同id数据，只会取最后一个
func (list ModelList) MapByID() map[int]Model {
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

// more...
```

如果觉得这里的方法不够用了，比如取多个列值，直接用这里生成的方法就需要遍历多次，这种可以通过遍历一次来优化的，就可以在包里新建文件 xxx_list_ext.go，然后在文件里给列表结构体添加相应的方法。

```go
// 结果使用匿名结构体，调用之后，后续使用的是里面的字段
func (list ModelList) ColumnIDAndName() struct{
    IDs []int
    Names []string
} {
    l := len(list)
    ids := make([]int, 0, l)
    names := make([]string, 0, l)
    for _, single := range list {
        ids = append(ids, single.ID)
        names = append(names, single.Name)
    }
    return struct{
        IDs []int
        Names []string
    }{
        IDs: ids,
        Names: names,
    }
}

```

[playground](https://play.golang.org/p/RTHKlv8WqyO)
