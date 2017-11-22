package main

import (
	//"fmt"

	"github.com/lioneagle/goutil/src/logger"
)

func main() {
	logger.SetLevel(logger.DEBUG)
	//logger.SetStackTraceLevel(logger.DEBUG)
	logger.Print("thor: a web editor")
	logger.Debug("thor: a web editor")
	logger.PrintStack()
}
