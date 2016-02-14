package cfg

import(
	"testing"
	"bytes"
)

func TestLocalIo(t *testing.T) {
	dir, err := testDir()
	if err != nil {
		t.Error("Unable to create temp dir")
	}
	filename := dir + "localio.conf"

	f := VagabondIoLocal{}
	exists := f.FileExists(filename);
	if exists {
		t.Error(
			"expected", false,
			"got", exists,
		)
	}

	expectedContents := []byte("foobarbaz")
	err = f.FileWrite(filename, expectedContents)
	if err != nil {
		t.Error(err)
	}
	contents, err := f.FileRead(filename)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(contents, expectedContents) {
		t.Error(
			"expected", expectedContents,
			"got", contents,
		)
	}
	err = f.FileDelete(filename)
	if err != nil {
		t.Error(err)
	}
}

func TestLocalExec(t *testing.T) {
	i := VagabondIoLocal{}
	err := i.Exec("exit 0")
	if err != nil {
		t.Error(
			"expected", "no error",
			"got", err,
		)
	}
	err = i.Exec("exit 1")
	if err == nil {
		t.Error(
			"expected", "error",
			"got", "no error",
		)
	}

}
