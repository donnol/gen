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

func main() {
	fmt.Printf("Gen is a tool for code generate.\n")

	// 解析标签
	rFlag := flag.Bool("r", false, "gen -r")
	flag.Parse()
	fmt.Printf("recursive: %v\n", *rFlag)

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
	if *rFlag {
		if err := filepath.Walk(dir, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
			if path == dir {
				return nil
			}
			// 获取所需目录
			if info.IsDir() && !strings.Contains(path, ".git") {
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
