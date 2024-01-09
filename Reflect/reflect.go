package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

type T struct {
	A int
	B string
}
type a map[string]int
type INT int

func main() {
	var ss a = a{}
	ss["0"] = 0
	ss["1"] = 1

	fmt.Println(reflect.ValueOf(ss))        //map[0:0 1:1] 获取输入参数接口中的数据的值
	fmt.Println(reflect.ValueOf(ss).Type()) //main.a   动态获取输入参数接口中的值的类型
	fmt.Println(reflect.ValueOf(ss).Kind()) //map    某种原始类型的抽象表示

	fmt.Println(reflect.TypeOf(ss))        //main.a 动态获取输入参数接口中的值的类型
	fmt.Println(reflect.TypeOf(ss).Kind()) //map  某种原始类型的抽象表示

	var bb io.Writer //接口
	bb = os.Stdout

	fmt.Println(reflect.ValueOf(bb))        //&{0xc0000cc280}
	fmt.Println(reflect.ValueOf(bb).Type()) //*os.File
	fmt.Println(reflect.ValueOf(bb).Kind()) //ptr

	fmt.Println(reflect.TypeOf(bb))        //*os.File
	fmt.Println(reflect.TypeOf(bb).Kind()) //ptr

	var i INT = 10
	fmt.Println(reflect.ValueOf(i))        //10
	fmt.Println(reflect.ValueOf(i).Type()) //main.INT
	fmt.Println(reflect.ValueOf(i).Kind()) //int

	fmt.Println(reflect.TypeOf(i))        //main.INT
	fmt.Println(reflect.TypeOf(i).Kind()) //int

}
