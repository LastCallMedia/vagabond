package cfg

import(
	"testing"
	"io/ioutil"
	"math/rand"
	"text/template"
	"bytes"
	"errors"
	"fmt"
)

func TestConfigFileAction(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Error("Unable to create temp dir")
	}

	a :=ConfigFileAction{
		Filename: dir + "file.conf",
	}
	nr, err := a.NeedsRun()
	if err != nil {
		t.Error(err)
	}
	if nr != true {
		t.Error(
			"expected", true,
			"got", nr,
		)
	}

	err = a.Run("foo")
	if err != nil {
		t.Error("got error handling run")
	}
}

func TestConfigFileActionTemplate(t *testing.T) {
	dir, err := testDir();
	if err != nil {
		t.Error("Unable to create temp dir")
	}
	filename := dir + "TestConfigFileActionTemplate.conf"
	a := ConfigFileAction{
		Filename: filename,
		Template: template.Must(template.New("testtmp").Parse("foobar{{.}}")),
		Contents: []byte("bazbar"),
	}

	a.Run("baz")
	if err := assertFileContentsEqual([]byte("foobarbaz"), filename); err !=nil {
		t.Error(err)
	}
}

func TestConfigFileActionContents(t *testing.T) {
	dir, err := testDir();
	if err != nil {
		t.Error("Unable to create temp dir")
	}
	filename := dir + "TestConfigFileActionTemplate.conf"
	a := ConfigFileAction{
		Filename: filename,
		Contents: []byte("bazbar"),
	}

	a.Run("baz")
	if err := assertFileContentsEqual([]byte("bazbar"), filename); err !=nil {
		t.Error(err)
	}
}

func assertFileContentsEqual(expected []byte, filename string) (err error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	if ! bytes.Equal(out, expected) {
		return errors.New(fmt.Sprintln(
			"expected", string(expected),
			"got", string(out),
		))
	}
	return
}

func testDir() (dir string, err error) {
	dir, err = ioutil.TempDir("", string(rand.Int()))
	return
}