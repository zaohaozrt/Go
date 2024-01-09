package gotest_test

// 文件名必须以“_test.go”结尾；
// 函数名必须以“BenchmarkXxx”开始；
// 使用命令“go test -bench.”即可开始性能测试；
import (
	gotestt "gotest1" //gotest1 模块里面的 gotestt包(gotest包的重命名)
	"testing"
)

func BenchmarkMakeSliceWithoutAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gotestt.MakeSliceWithoutAlloc()
	}
}

func BenchmarkMakeSliceWithPreAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gotestt.MakeSliceWithPreAlloc()
	}
}
