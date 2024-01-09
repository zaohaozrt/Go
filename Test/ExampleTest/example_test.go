package gotest_test

import "gotest"

// 例子测试函数名需要以”Example”开头；
// 检测单行输出格式为“// Output: <期望字符串>”；
// 检测多行输出格式为“// Output: \ <期望字符串> \ <期望字符串>”，每个期望字符串占一行；
// 检测无序输出格式为”// Unordered output: \ <期望字符串> \ <期望字符串>”，每个期望字符串占一行；
// 测试字符串时会自动忽略字符串前后的空白字符；
// 如果测试函数中没有“Output”标识，则该测试函数不会被执行；
// 执行测试可以使用go test，此时该目录下的其他测试文件也会一并执行；
// 执行测试可以使用go test <xxx_test.go>，此时仅执行特定文件中的测试函数；
// 检测单行输出
func ExampleSayHello() {
	gotest.SayHello()
	// OutPut: Hello World
}

// 检测多行输出
func ExampleSayGoodbye() {
	gotest.SayGoodbye()
	// OutPut:
	// Hello,
	// goodbye
}

// 检测乱序输出
func ExamplePrintNames() {
	gotest.PrintNames()
	// Unordered output:
	// Jim
	// Bob
	// Tom
	// Sue
}
