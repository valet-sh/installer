package utils

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

var DebugMode bool
var LogFile *os.File

func Println(data any, args ...any) {
	message := fmt.Sprintln(append([]any{data}, args...)...)
	writeToLog(message)
	color.Info.Print(message)
}

func Printf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	if len(message) > 0 && message[len(message)-1] != '\n' {
		message += "\n"
	}
	writeToLog(message)
	color.Info.Print(message)
}

func Debug(data any, args ...any) {
	message := fmt.Sprintln(append([]any{data}, args...)...)
	writeToLog(message)
	if DebugMode {
		color.Comment.Print(message)
	}
}

func Debugf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	if len(message) > 0 && message[len(message)-1] != '\n' {
		message += "\n"
	}
	writeToLog(message)
	if DebugMode {
		color.Comment.Print(message)
	}
}

func writeToLog(message string) {
	if LogFile != nil {
		_, err := LogFile.WriteString(message)
		if err != nil {
			fmt.Printf("Error writing to log file: %v\n", err)
		}
	}
}
