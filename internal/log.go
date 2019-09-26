package papersave

import (
	"fmt"
	"os"
)

func _log(prefix, msg string, args ...interface{}) {
	fmt.Printf("[%s] ", prefix)
	fmt.Printf(msg, args...)
	fmt.Printf("\n")
}

func Info(msg string, args ...interface{}) {
	_log("INFO", msg, args...)
}

func Debug(msg string, args ...interface{}) {
	_log("DEBUG", msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	_log("FATAL", msg, args...)
}

func Panicp(err error) {
	if err != nil {
		Fatal("%s", err)
		os.Exit(1)
	}
}
