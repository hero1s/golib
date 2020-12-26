package internal

import (
	"bytes"
	goformat "go/format"
	"io/ioutil"
	"text/template"
)

const regularPerm = 0666

type defaultTemplate struct {
	name     string
	text     string
	goFmt    bool
	savePath string
}

func With(name string) *defaultTemplate {
	return &defaultTemplate{
		name: name,
	}
}
func (t *defaultTemplate) Parse(text string) *defaultTemplate {
	t.text = text
	return t
}

func (t *defaultTemplate) GoFmt(format bool) *defaultTemplate {
	t.goFmt = format
	return t
}

func (t *defaultTemplate) SaveTo(data interface{}, path string, forceUpdate bool) error {
	if FileExists(path) && !forceUpdate {
		return nil
	}

	output, err := t.Execute(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, output.Bytes(), regularPerm)
}

func (t *defaultTemplate) Execute(data interface{}) (*bytes.Buffer, error) {
	tem, err := template.New(t.name).Parse(t.text)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err = tem.Execute(buf, data); err != nil {
		return nil, err
	}

	if !t.goFmt {
		return buf, nil
	}

	formatOutput, err := goformat.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	buf.Reset()
	buf.Write(formatOutput)
	return buf, nil
}
