package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type ConfigFileIo interface {
	Read(env *Environment, filename string) []byte
	Write(env *Environment, filename string, block []byte)
}

type ConfigFileLocalIo struct{}

func (io ConfigFileLocalIo) Read(env *Environment, filename string) (block []byte) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist.  No worries.
		return
	}

	block, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed reading ", filename, err)
	}
	return
}
func (io ConfigFileLocalIo) Write(env *Environment, filename string, block []byte) {
	cmd := exec.Command("sudo", "tee", filename)
	w, _ := cmd.StdinPipe()
	w.Write(block)
	w.Close()

	_, err := cmd.Output()

	if err != nil {
		log.Fatal("Unable to write file: ", filename)
	}
}

type ConfigFileMachineIo struct{}

func (io ConfigFileMachineIo) Read(env *Environment, filename string) (block []byte) {
	// Allow the file not to exist.  cat returns 1 if the file doesn't exist
	block, err := env.GetMachine().Exec(fmt.Sprintf("cat %s || return 0", filename)).Output()
	if err != nil {
		log.Fatalf("Failed reading %s: %s", filename, err)
	}
	return
}

func (io ConfigFileMachineIo) Write(env *Environment, filename string, block []byte) {
	machine := env.GetMachine()
	cmd := machine.Exec("sudo tee " + filename)
	w, _ := cmd.StdinPipe()
	w.Write(block)
	w.Close()

	_, err := cmd.Output()

	if err != nil {
		log.Fatal("Unable to write file: ", filename)
	}

	_, err = machine.Exec("sync").Output()
	if err != nil {
		log.Fatalf("Unable to sync file %s", filename)
	}
}
