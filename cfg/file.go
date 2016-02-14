package cfg

import(
	"text/template"
	"bytes"
	"bufio"
)

type ConfigFileAction struct {
	Filename string
	Contents []byte
	Target	 string
	Template *template.Template
	ValidateCmd string
}

func (a *ConfigFileAction) GetIo() VagabondIo {
	if a.Target == "" {
		return VagabondIoLocal{}
	}
	return VagabondIoMachine{a.Target}
}

func (a *ConfigFileAction) NeedsRun() (bool, error) {
	exists := a.GetIo().FileExists(a.Filename)
	return !exists, nil
}

func (a *ConfigFileAction) Run(vars interface{}) (err error) {
	var contents []byte
	if a.Template != nil {
		contents, err = getTemplateContents(a.Template, vars)
		if err != nil {
			return
		}
	} else {
		contents = a.Contents
	}

	err = a.GetIo().FileWrite(a.Filename, contents)
	if err != nil {
		return
	}
	if a.ValidateCmd != "" {
		err = a.GetIo().Exec(a.ValidateCmd)
	}
	return
}

func getTemplateContents(tpl *template.Template, vars interface{}) (contents []byte, err error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err = tpl.Execute(writer, vars)
	if err != nil {
		return
	}
	writer.Flush()
	return buf.Bytes(), err
}