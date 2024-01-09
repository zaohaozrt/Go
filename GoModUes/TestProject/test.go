package main

import (
	"fmt"
	testtttt "mypackage" //mypackage 模块里面的 testtttt 包(test1 包的重命名)
)

func main() {
	fmt.Println(testtttt.Sum(1, 2))
}
