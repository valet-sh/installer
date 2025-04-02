package utils

import (
	"fmt"
)

var DebugMode bool

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
	if DebugMode {
		fmt.Printf(data, args)
	}
}
