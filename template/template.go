package template

import (
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

// Template 模板，用于生成和格式化文件
type Template struct {
}

// Execute 执行
func (temp *Template) Execute(wr io.Writer, name, text string, data interface{}) error {
	t := template.New(name)

	t, err := t.Parse(text)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := t.Execute(wr, data); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// FileOption 文件选项
type FileOption struct {
	PkgName     string
	ImportPaths []string
	Content     string
}

// SpliceFile 拼接文件
func (temp *Template) SpliceFile(opt FileOption) string {
	// 包
	pkgContent := fmt.Sprintf("package %s\n", opt.PkgName)

	// 导入
	var importContent string
	for _, path := range opt.ImportPaths {
		importContent += fmt.Sprintf("\"%s\"\n", path)
	}

	// 类型/方法等内容
	file := fmt.Sprintf("%s\n import (\n%s\n)\n %s", pkgContent, importContent, opt.Content)

	return file
}

// WriteFile 写文件
func (temp *Template) WriteFile(filename string, content []byte) error {
	// 写入前先format
	formatContent, err := format.Source(content)
	if err != nil {
		return errors.WithMessagef(err, "%s", content)
	}

	if err := ioutil.WriteFile(filename, formatContent, os.ModePerm); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
