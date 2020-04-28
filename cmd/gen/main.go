package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/donnol/gen/list"
	"github.com/donnol/gen/parser"
	"github.com/donnol/gen/template"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

type excludeFlags []string

func (i *excludeFlags) String() string {
	return "exclude dir"
}

func (i *excludeFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *excludeFlags) Type() string {
	return ""
}

var (
	rootCmd = cobra.Command{
		Use:   "gen",
		Short: "a tool for code generate",
		Long:  "gen struct method, something like: ColumnXXX etc.",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func main() {
	fmt.Printf("Gen is a tool for code generate.\n")

	// 获取目录
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("dir: %s\n", dir)

	// 解析标签
	var rFlag bool
	rootCmd.PersistentFlags().BoolVarP(&rFlag, "r", "r", false, "recursive parse dir")
	var excludeFlags excludeFlags
	rootCmd.PersistentFlags().VarP(&excludeFlags, "exclude", "e", "exclude dir")
	var typ string
	rootCmd.PersistentFlags().StringVarP(&typ, "type", "t", "", "specify type with path, like: github.com/pkg/errors.Frame")
	var method string
	rootCmd.PersistentFlags().StringVar(&method, "method", "", "specify method will gen for the struct list, like: Where")
	var saveToFile bool
	rootCmd.PersistentFlags().BoolVarP(&saveToFile, "w", "w", false, "save to file")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "gen struct list method",
		Long:  "gen struct list method, like: ColumnXXX etc.",
		Run: func(cmd *cobra.Command, args []string) {
			// 没有指定类型，则从当前目录开始遍历，为含有@gen标记的结构体生成方法
			if strings.TrimSpace(typ) == "" {
				genPkg(dir, rFlag, excludeFlags)
			} else { // 指定了类型，则只对该类型生成方法
				fmt.Printf("=== list %+v, %s\n", typ, method)

				// TODO:
				// 指定了结构体，则不需要遍历目录，解析type内容，如果有包路径则找路径内的结构体信息，如果没有包路径则直接在当前目录找结构体信息
				importPath, typName, typNameWithPath := parser.GetImportPathAndTypeName(typ)
				p := parser.New()
				t := &template.Template{}
				opt := list.Option{}

				if strings.TrimSpace(method) != "" {
					methods := strings.Split(method, ";")
					for _, met := range methods {
						opt.Methods = append(opt.Methods, list.Method{
							Name: met,
						})
					}
				}

				l := list.New(p, t, opt)
				fmt.Printf("=== path: %s, %s, %s, %+v\n", importPath, typName, typNameWithPath, opt)
				if err := l.Parse(importPath); err != nil {
					fmt.Printf("=== parse err: %+v\n", err)
				}
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("recursive: %v, exclude: %+v\n", rFlag, excludeFlags)
}

func genPkg(dir string, rFlag bool, excludeFlags []string) {
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
	p := parser.New()
	t := &template.Template{}
	opt := list.Option{UseAnnotation: true}
	l := list.New(p, t, opt)
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
