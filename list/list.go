package list

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/donnol/gen/parser"
	"github.com/donnol/gen/template"
	toollist "github.com/donnol/tools/list"
	"github.com/pkg/errors"
)

// 生成文件后缀
const (
	GenFileSuffix = "_gen_list"
)

// 指令
const (
	commandName  = "list"
	attrColumn   = "column"
	attrMap      = "map"
	attrSliceMap = "slicemap"

	joinCommandName = "join"
	joinEqual       = "="
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

	// 没有内容，则不生成文件
	if len(content) == 0 {
		return nil
	}

	// 写入，添加包名和依赖导入
	importPaths := parser.ImportPathMap(importPathMap).Keys()
	importPaths = toollist.Filter(importPaths, pkg.ImportPath)
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

func (list *List) output(pkg parser.Pkg) (string, parser.ImportPathMap, []byte, error) {
	var pkgName = pkg.Name
	var importPathMap = make(parser.ImportPathMap)
	var content []byte

	var buf = bytes.NewBuffer(content)
	for _, singleStruct := range pkg.Structs {
		importPathMap = importPathMap.MergeKey(
			singleStruct.Info.GetNotEmptyImportPathMap())

		structName := singleStruct.Name
		typNameWithPath := singleStruct.Info.GetTypNameWithPath(pkg.ImportPath)

		existFieldCommand := false
		for _, singleField := range singleStruct.Fields {
			if singleField.Info.Commands.ExistCommand(commandName) {
				existFieldCommand = true
				break
			}
		}
		if singleStruct.Info.Commands.ExistCommand(commandName) ||
			existFieldCommand {
			if err := list.template.Execute(buf, "List", typText, map[string]interface{}{
				"typName":         structName,
				"typNameWithPath": typNameWithPath,
			}); err != nil {
				return pkgName, importPathMap, content, errors.WithStack(err)
			}
		}

		for _, singleField := range singleStruct.Fields {
			importPathMap = importPathMap.MergeKey(singleField.Info.GetNotEmptyImportPathMap())

			fieldName := singleField.Name
			fieldTypNameWithPath := singleField.Info.GetTypNameWithPath(pkg.ImportPath)

			// 取列
			if singleField.Info.Commands.ExistCommandAttr(commandName, attrColumn) {
				if err := list.template.Execute(buf, "List", columnMethodText, map[string]interface{}{
					"typName":   structName,
					"fieldName": fieldName,
					"fieldType": fieldTypNameWithPath,
				}); err != nil {
					return pkgName, importPathMap, content, errors.WithStack(err)
				}
			}
			// 映射
			if singleField.Info.CanUseAsMapKey {
				for _, methodText := range []struct {
					attr string
					text string
				}{
					{attrMap, mapMethodText},           // 列取映射
					{attrSliceMap, sliceMapMethodText}, // 列取数组映射
				} {
					if singleField.Info.Commands.ExistCommandAttr(commandName, parser.Attr(methodText.attr)) {
						if err := list.template.Execute(buf, "List", methodText.text, map[string]interface{}{
							"typName":         structName,
							"typNameWithPath": typNameWithPath,
							"fieldName":       fieldName,
							"fieldType":       fieldTypNameWithPath,
						}); err != nil {
							return pkgName, importPathMap, content, errors.WithStack(err)
						}
					}
				}
			}

			// join
			if singleField.Info.Commands.ExistCommand(joinCommandName) {
				joinTyp,
					joinTypField,
					joinTypWithPath,
					joinTypFieldTyp := singleField.Info.Commands.GetJoinTyp(
					list.parser,
					singleStruct.Info.ImportPath,
					pkg.Structs,
				)

				if err := list.template.Execute(buf, "List", joinMethodText, map[string]interface{}{
					"typName":         structName,
					"typNameWithPath": typNameWithPath,
					"typFieldName":    fieldName,
					"joinTyp":         joinTyp,
					"joinTypWithPath": joinTypWithPath,
					"joinTypField":    joinTypField,
					"joinTypFieldTyp": joinTypFieldTyp,
				}); err != nil {
					return pkgName, importPathMap, content, errors.WithStack(err)
				}
			}
		}
	}
	content = buf.Bytes()

	return pkgName, importPathMap, content, nil
}

func (list *List) getFileName(pkgDir, pkgName string) string {
	filename := fmt.Sprintf("%s%s.go", pkgName, GenFileSuffix)
	filename = filepath.Join(pkgDir, filename)
	return filename
}
