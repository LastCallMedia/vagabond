package cfg

import(
	"text/template"
	"bytes"
	"bufio"
	"regexp"
	"errors"
	"fmt"
)

const (
	autoconfig_start = "#VAGABONDAUTOCONFIG"
	autoconfig_end   = "#VAGABONDAUTOCONFIGEND"
)

type ConfigFileAction struct {
	Filename string
	Contents []byte
	Target	 string
	Template *template.Template
	TemplateVars interface{}
	Append 	bool
}

func (a ConfigFileAction) GetName() string {
	act := "Recreating"
	if a.Append {
		act = "Appending to"
	}
	return fmt.Sprintf("%s %s", act, a.Filename)
}

func (a *ConfigFileAction) GetIo() VagabondIo {
	if a.Target == "" {
		return VagabondIoLocal{}
	}
	return VagabondIoMachine{a.Target}
}

func (a ConfigFileAction) NeedsRun() (bool, error) {
	i := a.GetIo()
	exists := i.FileExists(a.Filename)
	if !exists {
		return true, nil
	}
	existing, err := i.FileRead(a.Filename)
	if err != nil {
		return false, err
	}
	expected, err := a.getExpectedContents()
	if err != nil {
		return false, err
	}
	return !bytes.Equal(existing, expected), err
}

func (a ConfigFileAction) Run() (err error) {
	contents, err := a.getExpectedContents()

	err = a.GetIo().FileWrite(a.Filename, contents)
	if err != nil {
		return errors.New(fmt.Sprintf("Error writing %s: %s", a.Filename, err))
	}
	return
}

func (a ConfigFileAction)getExpectedContents() (contents []byte, err error) {
	if a.Template != nil {
		contents, err = getTemplateContents(a.Template, a.TemplateVars)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error generating file contents for %s: %s", a.Filename, err))
			return
		}
	} else {
		contents = a.Contents
	}

	i := a.GetIo()
	if a.Append {
		existing := []byte{}
		if i.FileExists(a.Filename) {
			existing, err = i.FileRead(a.Filename)
			if err != nil {
				err = errors.New(fmt.Sprintf("Error reading %s: %s", a.Filename, err))
				return contents, err
			}
		}
		contents = appendContents(existing, contents)
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

func appendContents(existing []byte, contents []byte) (newcontents []byte) {
	re := regexp.MustCompile("(?s)" + autoconfig_start + ".*" + autoconfig_end)
	newcontents = re.ReplaceAll(existing, []byte(""))
	newcontents = bytes.TrimRight(newcontents, "\n")
	if !bytes.Equal(newcontents, []byte("")) {
		// Add a newline if this isn't going to be the first line in the file.
		newcontents = append(newcontents, []byte("\n")...)
	}

	newcontents = append(newcontents, autoconfig_start+"\n"...)
	newcontents = append(newcontents, contents...)
	newcontents = append(newcontents, "\n"+autoconfig_end+"\n"...)
	return
}