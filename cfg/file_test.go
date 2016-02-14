package cfg

import(
	"testing"
	"io/ioutil"
	"math/rand"
	"text/template"
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestConfigFileAction(t *testing.T) {
	dir, err := testDir()
	assert.NoError(t, err)

	a :=ConfigFileAction{
		Filename: dir + "file.conf",
	}
	nr, err := a.NeedsRun()
	assert.NoError(t, err)
	assert.True(t, nr)

	err = a.Run()
	assert.NoError(t, err)
}

func TestConfigFileNeedsRun(t *testing.T) {
	dir, err := testDir()
	assertNoErr(t, err)

	filename := dir + "needsrun.conf"

	a := ConfigFileAction{
		Filename: filename,
		Contents: []byte("foo"),
	}
	nr, err := a.NeedsRun()
	assert.NoError(t, err)
	assert.True(t, nr)

	ioutil.WriteFile(filename, []byte("baz"), 0777)
	nr, err = a.NeedsRun()
	assert.NoError(t, err)
	assert.True(t, nr)

	ioutil.WriteFile(filename, []byte("foo"), 0777)
	nr, err = a.NeedsRun()
	assert.NoError(t, err)
	assert.False(t, nr)

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
		TemplateVars:"baz",
		Contents: []byte("bazbar"),
	}

	a.Run()
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

	a.Run()
	if err := assertFileContentsEqual([]byte("bazbar"), filename); err !=nil {
		t.Error(err)
	}
}

func TestConfigFileActionAppend(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Error("Unable to create temp dir")
	}
	filename := dir + "TestConfigFileAppend.conf"
	a := ConfigFileAction{
		Filename: filename,
		Contents: []byte("bazbar"),
		Append: true,
	}
	a.Run()
	expected := []byte("#VAGABONDAUTOCONFIG\nbazbar\n#VAGABONDAUTOCONFIGEND\n")
	if err := assertFileContentsEqual(expected, filename); err != nil {
		t.Error(err)
	}
}

func TestConfigFileActionAppendToExistingFile(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Error("Unable to create temp dir")
	}
	filename := dir + "TestConfigFileAppendToExisting.conf"
	ioutil.WriteFile(filename, []byte("sample\n"), 0777)

	a := ConfigFileAction{
		Filename: filename,
		Contents: []byte("bazbar"),
		Append: true,
	}
	a.Run()

	expected := []byte("sample\n#VAGABONDAUTOCONFIG\nbazbar\n#VAGABONDAUTOCONFIGEND\n")
	if err := assertFileContentsEqual(expected, filename); err != nil {
		t.Error(err)
	}
}

func assertBoolEquals(t *testing.T, expected bool, val bool) {
	if val != expected {
		t.Error(
			"expected", expected,
			"got", val,
		)
	}
}

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
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