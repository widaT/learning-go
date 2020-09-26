# 看一下有趣的汇编运用例子

```go
package main

import "fmt"

const s = "Go101.org"
//len(s)==1
//1 << 9 == 512
// 512/128 =4
var a byte = 1 << len(s) / 128
var b byte = 1 << len(s[:]) / 128

func main() {
	fmt.Println(a, b)
}
 ```

 ```bash
$ go run main.go
4 0
 ```

为什么结果会不一样?建议读者运行下这个程序,不着急看下面的结果.

如果你不懂go汇编,你很难知道到底发生了什么? 你懂go汇编的话,你看容易在汇编代码中找到端倪.

```bash
go tool compile -l -S -N main.go
```

在汇编代码中看到
```
"".a SNOPTRDATA size=1
	0x0000 04                                               .
"".b SNOPTRBSS size=1
```

`var a byte = 1 << len(s) / 128` 这段代码被go编译器(go 1.14)直接优化成a初始化为4的变量,b这没被初始化,b的赋值在main package的init函数中

```
	0x0015 00021 (main.go:8)	MOVQ	AX, ""..autotmp_0+8(SP)
	0x001a 00026 (main.go:8)	MOVQ	$9, ""..autotmp_0+16(SP)
	0x0023 00035 (main.go:8)	MOVQ	$9, ""..autotmp_1(SP)
	0x002b 00043 (main.go:8)	JMP	45
	0x002d 00045 (main.go:8)	MOVB	$0, "".b(SB)
```

上面的代码很有趣,b被赋值为0,至于0从哪里冒出来的,完全看不到.
我们很有理由怀疑go的编译器出bug了,bug不是因为b的结果不对,而是应该对len(s)和len(s[:])用同一套规则,假定 `1 << len(s[:])`已经 让byte类型溢出了其结果为0,那么应该`1 << len(s) `这个也应该是一样的结果.但是编译器s是常量s[:]是变量,len(s)还是常量,len(s[:])是变量,常量在运算过程中有隐式类型转换,`1 << len(s)`1会变成int类型`var a byte = int(1) << len(s) / 128`,变量则没有,所以b为`uint8(1) << len(s[:])`的结果为0.


## 参考资料

[go100and1](https://twitter.com/go100and1/status/1309188138015760385)