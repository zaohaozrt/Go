package gotest_test

import (
	"gotest"
	"testing"
)

// 子测试适用于单元测试和性能测试；
// 子测试可以控制并发；
// 子测试提供一种类似table-driven风格的测试；
// 子测试可以共享setup和tear-down；
// sub1 为子测试，只做加法测试
// go test subunit_test.go -v -run Sub/A=
func sub1(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := gotest.Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

// sub2 为子测试，只做加法测试
func sub2(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := gotest.Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

// sub3 为子测试，只做加法测试
func sub3(t *testing.T) {
	var a = 1
	var b = 2
	var expected = 3

	actual := gotest.Add(a, b)
	if actual != expected {
		t.Errorf("Add(%d, %d) = %d; expected: %d", a, b, actual, expected)
	}
}

// TestSub 内部调用sub1、sub2和sub3三个子测试
func TestSub(t *testing.T) {
	// setup code

	t.Run("A=1", sub1)
	t.Run("A=2", sub2)
	t.Run("B=1", sub3)

	// tear-down code
}
