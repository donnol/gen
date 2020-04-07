package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/donnol/gen/list"
	"github.com/donnol/gen/parser"
	"github.com/donnol/gen/template"
	"golang.org/x/mod/modfile"
)

type excludeFlags []string

func (i *excludeFlags) String() string {
	return "my string representation"
}

func (i *excludeFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	fmt.Printf("Gen is a tool for code generate.\n")

	// 解析标签
	var rFlag bool
	flag.BoolVar(&rFlag, "r", false, "recursive parse dir")
	var excludeFlags excludeFlags
	flag.Var(&excludeFlags, "exclude", "exclude dir")
	flag.Parse()
	fmt.Printf("recursive: %v, exclude: %+v\n", rFlag, excludeFlags)

	// 获取目录
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("dir: %s\n", dir)

	// 解析目录里的go.mod文件，获取模块名
	modfilePath := filepath.Join(dir, "go.mod")
	content, err := ioutil.ReadFile(modfilePath)
	if err != nil {
		panic(err)
	}
	modPath := modfile.ModulePath(content)
	fmt.Printf("modPath: %s\n", modPath)

	// 获取模块里的所有包
	allPkgPath := []string{modPath}
	if rFlag {
		if err := filepath.Walk(dir, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
			if path == dir {
				return nil
			}
			// 获取所需目录
			if info.IsDir() {
				// 排除
				for _, exd := range excludeFlags {
					if strings.Contains(path, exd) {
						return nil
					}
				}

				// 过滤没有go文件的
				fileInfos, err := ioutil.ReadDir(path)
				if err != nil {
					panic(err)
				}
				haveGoFile := false
				for _, fi := range fileInfos {
					ext := filepath.Ext(fi.Name())
					if ext == ".go" {
						haveGoFile = true
						break
					}
				}
				if !haveGoFile {
					return nil
				}

				// 替换系统目录为包路径
				pkgPath := strings.Replace(path, dir, modPath, -1)

				allPkgPath = append(allPkgPath, pkgPath)
			}

			return nil
		})); err != nil {
			panic(err)
		}
	}

	// 生成代码
	p := parser.New(parser.Option{
		IgnoreFileSuffix: []string{
			list.GenFileSuffix,
		},
	})
	t := &template.Template{}
	l := list.New(p, t)
	var pkgNum int
	for _, pkgPath := range allPkgPath {
		fmt.Printf("=== parse path: %s\n", pkgPath)
		if err := l.Parse(pkgPath); err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		pkgNum++
	}

	// 结果统计
	fmt.Printf("Job done: %d/%d\n", pkgNum, len(allPkgPath))
}
