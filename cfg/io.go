package cfg

import(
	"os"
	"io/ioutil"
	"os/exec"
	"fmt"
	"errors"
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
	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return
}

type VagabondIoMachine struct {
	Name  string
}
func (i VagabondIoMachine)FileExists(filename string) (bool) {
	err := machineCommand(i.Name, fmt.Sprintf("test -f %s", filename)).Run()
	if err != nil {
		return false
	}
	return true
}

func (i VagabondIoMachine)FileRead(filename string) (out []byte, err error) {
	cmd := machineCommand(i.Name, fmt.Sprintf("cat %s", filename))
	stdErr := []byte("")
	ew, _ := cmd.StderrPipe()
	ew.Read(stdErr)
	ew.Close()

	out, err = cmd.Output()
	if err != nil {
		err = errors.New(string(stdErr))
	}
	return
}

func (i VagabondIoMachine)FileWrite(filename string, contents []byte) (err error) {
	cmd := machineCommand(i.Name, fmt.Sprintf("sudo tee %s", filename))
	w, _ := cmd.StdinPipe()
	w.Write(contents)
	w.Close()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	err = machineCommand(i.Name, "sync").Run()
	return
}

func (i VagabondIoMachine)FileDelete(filename string) (err error) {
	out, err := machineCommand(i.Name, fmt.Sprintf("sudo rm %s", filename)).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return
}

func (i VagabondIoMachine)Exec(command string) (err error) {
	out, err := machineCommand(i.Name, command).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return
}

func machineCommand(machineName string, command string) *exec.Cmd {
	return exec.Command("docker-machine", "ssh", machineName, command)
}

