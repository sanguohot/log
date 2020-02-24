package main

import (
	"fmt"
	"github.com/sanguohot/log/test/child"
	"github.com/sanguohot/log/v3"
	"os"
	"runtime"
)

func init() {
	os.Setenv("LOG_FILE", "v3.log")
	os.Setenv("LOG_TYPE", "file")
	log.Init()
}

func main() {
	fmt.Println(runtime.GOARCH[:3])
	fmt.Println(runtime.GOOS)
	fmt.Println(os.Args)
	log.Sugar.Info("hello v3")
	child.TestV3()
}
