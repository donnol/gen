# gen

[![Go Report Card](https://goreportcard.com/badge/github.com/donnol/gen)](https://goreportcard.com/report/github.com/donnol/gen)
[![GoDoc](https://pkg.go.dev/mod/github.com/donnol/gen?status.svg)](https://pkg.go.dev/mod/github.com/donnol/gen)

生成。

## 为什么

一切都是为了复用。

并在达到复用目标的同时，还要保证代码的显式、清晰。

NOTE: join 和 derive 有点复杂了，不建议使用。

## 使用

安装工具`gen`:

```sh
go get -u github.com/donnol/gen/cmd/gen
```

先要在写结构体的时候在文档做好标志(写的时候请忽略标记后的说明)，如下：

命令格式：

| 命令 | 子命令 | [属性]    | [参数]   |
| ---- | ------ | --------- | -------- |
| gen  | list   | column... | .ID      |
| gen  | join   |           | =User.ID |

```go
// Model 模型
//
// @gen list: 出现这个标记表示这个结构体需要生成相应的列表结构体，这个是可选的，因为字段标记里出现了也要生成
type Model struct {
    // 内嵌结构体
    // @gen list column .ID
    Inner

    // @gen list column: 出现这个标记表示列表结构体应该有取ID列值的方法，如果只有这个标记，没有结构体标记，也要生成相应的列表结构体
    ID int // 记录id

    // @gen list map: 表示要取Name列映射
    // @gen list slicemap: 表示要取Name列的数组映射
    // @gen list [map, slicemap]: 多个写在一起
    // @gen list [map, slicemap]; xxx [xxx, xxx]: 如果后面有多种指令，使用这个格式(NOTE:未实现)
    Name string // 名称

    // @gen join =User.ID: 表示关联User结构体的ID字段，这样就需要知道结构体名和字段名，如果结构体在其它包，还需要有包路径，如：github.com/pkg/errors.XXX.YYY，或相对路径：./pkgpath.XXX.YYY, ../pkgpath.XXX.YYY
    UserID int // 用户ID
    UserName string // 用户名

    Old bool // 旧
    Height float64 // 高度
    CreatedAt time.Time // 创建时间
}

type Inner struct {
    ID int
}

type User struct {
    ID int
    Name string
}
```

然后再在项目根目录运行`gen -r`，命令就会根据 go.mod 文件找到模块信息，遍历项目里的包，在包目录里生成包含所需结构和方法的文件。

生成内容可看下面的说明。

## 结构体列表

给定一个结构体，生成相应的列表结构体，并在列表结构体上生成方法。方法包括取列，列映射，列数组映射等

如：

给定结构体：

```go
type Model struct {
    ID int
    Name string
    UserID int
    UserName string
    Old bool
    Height float64
    CreatedAt time.Time
}

type User struct {
    ID int
    Name string
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

func (list ModelList) JoinUserListByUserIDEqualID(u UserList, f func(Model, User) Model) ModelList {
    userMap := u.MapByID()

    result := make(ModelList, len(list), len(list))
    for i, single := range list {
        tmp := f(single, userMap[single.UserID])

        result[i] = tmp
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
