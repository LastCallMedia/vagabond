package util

import (
	"fmt"
	"os"
)

func Fatal(a ...interface{}) {
	fmt.Print(FgRed, a, Reset, "\n")
	os.Exit(1)
}

func Fatalf(message string, a ...interface{}) {
	fmt.Printf(FgRed+message+Reset+"\n", a...)
	os.Exit(1)
}

func Success(a ...interface{}) {
	fmt.Print(FgGreen, a, Reset, "\n")
}

func Successf(message string, a ...interface{}) {
	fmt.Printf(FgGreen+message+Reset+"\n", a...)
}
