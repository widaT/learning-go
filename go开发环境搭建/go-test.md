# go test 

go test是golang的轻量化单元测试工具，结合testing这个官方包，可以很方便为golang程序写单元测试和基准测试。

在之前go build命令介绍的时候说过，go build 编译包时，会忽略以“_test.go”结尾的文件。这些别忽略的_test.go文件正是go test的一部分。

在以_test.go为结尾的文件中，常见有两种类型的函数：

- 测试函数，以Test为函数名前缀的函数，用于测试程序是否按照正确的方式运行；使用go test命令执行测试函数并报告测试结果是PASS或FAIL。
- 基准测试(benchmark)函数，以Benchmark为函数名前缀的函数，它们用于衡量一些函数的性能；基准测试函数中一般会多次调用被测试的函数，然后收集平均执行时间。


无论是测试函数或者是基准测试函数都必须`import testing` 

## 测试函数

测试函数的签名

```go
func TestA(t *testing.T){

}
```
例如我们看下go的官方包 bytes.Compare 函数的测试


```go
package test

import (
	"bytes"
	"testing"
)

var Compare = bytes.Compare

func TestCompareA(t *testing.T) {
	var b = []byte("Hello Gophers!")
	if Compare(b, b) != 0 {
		t.Error("b != b")
	}
	if Compare(b, b[:1]) != 1 {
		t.Error("b > b[:1] failed")
	}
}
```

我们执行下 ```go test```，这个命令会遍历当前目录下所有的测试函数.
```
$ go test 
PASS
ok      github.com/wida/gocode/test  0.001s
```

参数-v用于打印每个测试函数的名字和运行时间：

```
$ go test -v
=== RUN   TestCompareA
--- PASS: TestCompareA (0.00s)
=== RUN   TestCompareB
--- PASS: TestCompareB (0.00s)
PASS
ok      github.com/wida/gocode/test  0.002s
```

参数-run 用于运行制定的测试函数,

```
$ go run -v -run "TestCompareA"
=== RUN   TestCompareA
--- PASS: TestCompareA (0.00s)
PASS
ok      github.com/wida/gocode/test  0.002s
```
run 后面的参数是正则匹配 -run "TestCompare" 会同时执行 TestCompareA 和 TestCompareB。


### 表驱动测试
在实际编写测试代码时，通常把要测试的输入值和期望的结果写在一起组成一个数据表（table），表（table）中的每条记录代表是一个含有输入值和期望值。还是看官方bytes.Compare的测试例子：

```go
var compareTests = []struct {
	a, b []byte
	i    int
}{
	{[]byte(""), []byte(""), 0},
	{[]byte("a"), []byte(""), 1},
	{[]byte(""), []byte("a"), -1},
	{[]byte("abc"), []byte("abc"), 0},
	{[]byte("abd"), []byte("abc"), 1},
	{[]byte("abc"), []byte("abd"), -1},
	{[]byte("ab"), []byte("abc"), -1},
	{[]byte("abc"), []byte("ab"), 1},
	{[]byte("x"), []byte("ab"), 1},
	{[]byte("ab"), []byte("x"), -1},
	{[]byte("x"), []byte("a"), 1},
	{[]byte("b"), []byte("x"), -1},
	// test runtime·memeq's chunked implementation
	{[]byte("abcdefgh"), []byte("abcdefgh"), 0},
	{[]byte("abcdefghi"), []byte("abcdefghi"), 0},
	{[]byte("abcdefghi"), []byte("abcdefghj"), -1},
	{[]byte("abcdefghj"), []byte("abcdefghi"), 1},
	// nil tests
	{nil, nil, 0},
	{[]byte(""), nil, 0},
	{nil, []byte(""), 0},
	{[]byte("a"), nil, 1},
	{nil, []byte("a"), -1},
}

func TestCompareB(t *testing.T) {
	for _, tt := range compareTests {
		numShifts := 16
		buffer := make([]byte, len(tt.b)+numShifts)
		for offset := 0; offset <= numShifts; offset++ {
			shiftedB := buffer[offset : len(tt.b)+offset]
			copy(shiftedB, tt.b)
			cmp := Compare(tt.a, shiftedB)
			if cmp != tt.i {
				t.Errorf(`Compare(%q, %q), offset %d = %v; want %v`, tt.a, tt.b, offset, cmp, tt.i)
			}
		}
	}
}
```
```
$ go test -v -run "TestCompareB"
=== RUN   TestCompareB
--- PASS: TestCompareB (0.00s)
PASS
ok      github.com/wida/gocode/test  0.004s
```

## 基准测试函数

基准测试函数的函数签名如下

```go
func BenchmarkTestB(b *testing.B) {
  
}
```

我们还是一官方的bytes.Compare 为例写一个基准测试

```go
func BenchmarkComare(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Compare([]byte("abcdefgh"), []byte("abcdefgh"))
	}
}
```

基准测试的运行 需要加 -bench参数

```
$ go test  -bench BenchmarkCompare
goos: linux
goarch: amd64
pkg: github.com/wida/gocode/test
BenchmarkCompare-4      30000000                35.4 ns/op
PASS
ok      github.com/wida/gocode/test  1.111s
```
报告显示我们的测试程序跑了30000000次，每次平均耗时35.4纳秒。

### 使用 b.ResetTimer
有些时候我们的基础测试函数逻辑有点复杂或者在准备测试数据，如下

```go
func BenchmarkComare(b *testing.B) {
    //准备数据
    ...
    b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compare([]byte("abcdefgh"), []byte("abcdefgh"))
	}
}
```
我们可以使用 b.ResetTimer() 将准备数据的时间排除在总统计时间之外。
