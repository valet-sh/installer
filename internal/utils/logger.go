package utils

import (
	"fmt"
	"os"
)

var DebugMode bool
var LogFile *os.File

func Println(data any, args ...any) {
	if DebugMode {
		if len(args) == 0 {
			fmt.Println(data)
		} else {
			fmt.Println(data, args)
		}
	}
}

func Printf(data string, args ...any) {
	debugMessage := fmt.Sprintf(data, args...)
	if DebugMode {
		fmt.Println(debugMessage)
	}
	if LogFile != nil {
		_, errr := LogFile.WriteString(debugMessage)
		if errr != nil {
			fmt.Println("Error writing to log file:", errr)
		}
	}
}
