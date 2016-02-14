package cfg

import(
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
)

type VagabondIo interface {
	FileExists(string) (bool)
	FileRead(string) ([]byte, error)
	FileWrite(string, []byte) (error)
	FileDelete(string) (error)
	Exec(string) (error)
}

type VagabondIoLocal struct {
}

func (i VagabondIoLocal)FileExists(filename string) (bool){
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist.  No worries.
		return false
	}
	return true
}

func (i VagabondIoLocal)FileRead(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (i VagabondIoLocal)FileWrite(filename string, contents []byte) (err error) {
	cmd := exec.Command("sudo", "tee", filename)
	w, _ := cmd.StdinPipe()
	w.Write(contents)
	w.Close()

	_, err = cmd.Output()
	return
}

func (i VagabondIoLocal)FileDelete(filename string) (err error) {
	err = exec.Command("sudo", "rm", filename).Run()
	return
}

func (i VagabondIoLocal)Exec(command string) (err error) {
	err = exec.Command("bash", "-c", command).Run()
	return
}

type VagabondIoMachine struct {
	Name  string
}
func (i VagabondIoMachine)FileExists(filename string) (bool) {
	err := exec.Command("docker-machine", i.Name, "ssh", "-c", fmt.Sprintf("test -f %s", filename))
	if err != nil {
		return false
	}
	return true
}

func (i VagabondIoMachine)FileRead(filename string) (out []byte, err error) {
	out, err = exec.Command("docker-machine", i.Name, "ssh", "-c", fmt.Sprintf("cat %s", filename)).Output()
	return
}

func (i VagabondIoMachine)FileWrite(filename string, contents []byte) (err error) {
	cmd := exec.Command("docker-machine", i.Name, "ssh", "-c", fmt.Sprintf("sudo tee %s", filename))
	w, _ := cmd.StdinPipe()
	w.Write(contents)
	w.Close()

	_, err = cmd.Output()
	return
}

func (i VagabondIoMachine)FileDelete(filename string) (err error) {
	err = exec.Command("docker-machine", i.Name, "ssh", "-c", fmt.Sprintf("sudo rm %s", filename)).Run()
	return
}

func (i VagabondIoMachine)Exec(command string) (err error) {
	err = exec.Command("docker-machine", i.Name, "ssh", "-c", command).Run()
	return
}
