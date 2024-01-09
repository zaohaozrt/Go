package main

import "fmt"

type Int int
type Str struct{}

type Map map[string]string
type Slice []int
type Chan chan int

type Func func(int, int) int

func main() {
	// 1类  只能用定义之后的类型赋值
	a := 10
	var b Int = Int(a)
	var c Int = 10

	var d Str = Str{}

	// 2类  可以用定义之前或者之后的类型赋值
	var e Map = map[string]string{"1": "2"}
	var f Map = Map{"1": "2"}

	var g Slice = []int{1, 2}
	var h Slice = Slice{1, 2}

	var i Chan = make(Chan, 0)
	var j Chan = make(chan int, 0)
	// 3类	只能用定义之前的类型赋值
	var k Func = func(i1, i2 int) int { return i1 + i2 }

	fmt.Println(a, b, c, d, e, f, g, h, i, j, k)
}
