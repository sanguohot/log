package main

import (
	"fmt"
	"github.com/sanguohot/log/test/child"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOARCH[:3])
	fmt.Println(runtime.GOOS)
	child.Test()
}
