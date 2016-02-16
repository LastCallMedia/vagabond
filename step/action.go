package step

import (
	"bufio"
	"bytes"
	"github.com/LastCallMedia/vagabond/config"
	"io/ioutil"
	"os/exec"
	"regexp"
	"text/template"
	"os"
)

const (
	autoconfig_start = "#VAGABONDAUTOCONFIG"
	autoconfig_end   = "#VAGABONDAUTOCONFIGEND"
)

type ConfigStep struct {
	Name     string
	NeedsRun func(envt *config.Environment) bool
	Run      func(envt *config.Environment) error
}

func (act ConfigStep) GetName() string {
	return act.Name
}

func appendConfigBlock(existing []byte, block []byte) (modified []byte) {
	re := regexp.MustCompile("(?s)" + autoconfig_start + ".*" + autoconfig_end)
	modified = re.ReplaceAll(existing, []byte(""))
	modified = bytes.TrimRight(modified, "\n")
	if !bytes.Equal(modified, []byte("")) {
		// Add a newline if this isn't going to be the first line in the file.
		modified = append(modified, []byte("\n")...)
	}

	modified = append(modified, autoconfig_start+"\n"...)
	modified = append(modified, block...)
	modified = append(modified, "\n"+autoconfig_end+"\n"...)
	return
}

func checkIfFileMatches(filename string, expected []byte) (bool, error) {
	existing, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}
	return bytes.Equal(existing, expected), nil
}

func doTemplateAppend(tplString string, data interface{}, filename string) (out []byte, err error) {
	addition, err := doTemplate(tplString, data)
	if err != nil {
		return
	}
	existing, err := ioutil.ReadFile(filename)
	// It's ok if this file doesn't exist yet.
	if err != nil && os.IsNotExist(err) {
		existing = []byte{}
		err = nil
	}
	out = appendConfigBlock(existing, addition)
	return
}

func doTemplate(tplString string, data interface{}) (out []byte, err error) {
	tpl := template.Must(template.New("bootlocalsh").Parse(tplString))

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err = tpl.Execute(writer, data)
	if err != nil {
		return
	}
	writer.Flush()
	return buf.Bytes(), err
}

func pipeInputToCmd(cmd *exec.Cmd, input []byte) {
	w, _ := cmd.StdinPipe()
	w.Write(input)
	w.Close()
}
