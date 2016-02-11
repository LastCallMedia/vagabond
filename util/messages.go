package util

import (
	"fmt"
	"os"
)

// Formats and prints an error message then exits
func Fatal(a ...interface{}) {
	fmt.Print(FgRed, a, Reset, "\n")
	os.Exit(1)
}

// Formats and prints an error message then exits
func Fatalf(message string, a ...interface{}) {
	fmt.Printf(FgRed+message+Reset+"\n", a...)
	os.Exit(1)
}

// Prints a success message
func Success(a ...interface{}) {
	fmt.Print(FgGreen, a, Reset, "\n")
}

// Formats and prints a success message
func Successf(message string, a ...interface{}) {
	fmt.Printf(FgGreen+message+Reset+"\n", a...)
}
