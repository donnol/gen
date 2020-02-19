package list

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/donnol/gen/parser"
	"github.com/donnol/gen/template"
	"github.com/pkg/errors"
)

// List 列表，解析结构体，生成列表结构体及取列，列映射，列数组映射等方法
type List struct {
	parser   *parser.Parser
	template *template.Template
}

// New 新建
func New(p *parser.Parser, t *template.Template) *List {
	return &List{
		parser:   p,
		template: t,
	}
}

// Parse 解析，对输入的导入路径进行解析，生成需要的结构体和方法，再写到同目录下的特定文件内
func (list *List) Parse(importPath string) error {
	// 解析
	pkg, err := list.parser.ParsePkg(importPath)
	if err != nil {
		return errors.WithStack(err)
	}

	// 生成
	pkgName, importPathMap, content, err := list.output(pkg)
	if err != nil {
		return err
	}

	// 写入，添加包名和依赖导入
	importPaths := ImportPathMap(importPathMap).Keys()
	importPaths = filter(importPaths, pkg.ImportPath)
	fileContent := list.template.SpliceFile(template.FileOption{
		PkgName:     pkgName,
		ImportPaths: importPaths,
		Content:     string(content),
	})
	err = list.template.WriteFile(list.getFileName(pkg.Dir, pkgName), []byte(fileContent))
	if err != nil {
		return err
	}

	return nil
}

func (list *List) output(pkg parser.Pkg) (string, map[string][]string, []byte, error) {
	var pkgName = pkg.Name
	var importPathMap = make(map[string][]string)
	var content []byte

	var buf = bytes.NewBuffer(content)
	for _, singleStruct := range pkg.Structs {
		importPath := singleStruct.ImportPath
		if importPath != "" {
			importPathMap[importPath] = append(importPathMap[importPath], singleStruct.TypNameWithPath)
		}

		typNameWithPath := singleStruct.TypNameWithPath
		if singleStruct.ImportPath == pkg.ImportPath {
			typNameWithPath = singleStruct.TypName
		}
		if err := list.template.Execute(buf, "List", typText, map[string]interface{}{
			"typName":         singleStruct.Name,
			"typNameWithPath": typNameWithPath,
		}); err != nil {
			return pkgName, importPathMap, content, errors.WithStack(err)
		}

		for _, singleField := range singleStruct.Fields {
			importPath := singleField.ImportPath
			if importPath != "" {
				importPathMap[importPath] = append(importPathMap[importPath], singleField.TypNameWithPath)
			}

			typNameWithPath := singleField.TypNameWithPath
			if singleField.ImportPath == pkg.ImportPath {
				typNameWithPath = singleField.TypName
			}
			if err := list.template.Execute(buf, "List", columnMethodText, map[string]interface{}{
				"typName":         singleStruct.Name,
				"typNameWithPath": singleStruct.TypNameWithPath,
				"fieldName":       singleField.Name,
				"fieldType":       typNameWithPath,
			}); err != nil {
				return pkgName, importPathMap, content, errors.WithStack(err)
			}
		}
	}
	content = buf.Bytes()

	return pkgName, importPathMap, content, nil
}

func (list *List) getFileName(pkgDir, pkgName string) string {
	filename := fmt.Sprintf("%s_list.go", pkgName)
	filename = filepath.Join(pkgDir, filename)
	return filename
}

// ImportPathMap 路径映射
type ImportPathMap map[string][]string

// Keys 键
func (m ImportPathMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func filter(keys []string, s string) []string {
	newKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		if key == s {
			continue
		}
		newKeys = append(newKeys, key)
	}
	return newKeys
}
