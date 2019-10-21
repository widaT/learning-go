# go 错误处理

在go1.13之前，go的错误处理方式代码写起来相当繁琐。go 1.13吸收了go社区一些优秀的错误处理方式（[pkg/errors](https://github.com/pkg/errors)），彻底解决被人诟病的问题。本文主要介绍的错误处理方式是基于go1.13的。

go 1.13的`error`包 增加的`errors.Unwrap`，`errors.As`，`errors.Is`三个方法。
同时 `fmt` 包增加 `fmt.Errorf("%w", err)`的方式来wrap一个错误。

我们通过代码来了解它们的用法。
```golang
package main

import (
	"fmt"
	"errors"
)
type Err struct {
	Code int
	Msg string
}
func (e *Err) Error() string  {
	return fmt.Sprintf("code : %d ,msg:%s",e.Code,e.Msg)
}
var A_ERR = &Err{-1,"error"}
func a()  error {
	return A_ERR
}

func b()  error {
	err := a()
	return fmt.Errorf("access denied: %w", err) //使用fmt.Errorf wrap 另一个错误
}

func main()  {
	err := b()
	er := errors.Unwrap(err)  //如果一个错误包含 Unwrap 方法则返回这个错误，如果没有则返回nil
	fmt.Println(er ==A_ERR )


	fmt.Println(errors.Is(err,A_ERR)) // 递归调用Unwrap判断是否包含 A_ERR
	var e = &Err{}
    fmt.Println( errors.As(err, &e))
    
	if errors.As(err, &e) {         // 递归调用Unwrap是否包含A_ERR，如果有这赋值给e
		fmt.Printf("code : %d ,msg:%s",e.Code,e.Msg)
	}
}
```

运行代码
```bash
$ go run main.go
true
true
true
code : -1 ,msg:error
```

错误为什么为要被wrap？

在一个函数A中错误发生的时候，我们会返回这个错误，函数B调用函数A拿到这个错，但是函数B不想做其他处理，它也返回错误，但是要打上自己的信息，说明这个错误经过了B函数，所以Wrap err就有了使用场景。

用了wrap后，错误是链状结构，我们用`errors.Unwrap`，逐级遍历err。还有我们有时候不一定会关心所有链条上的错误类型，我们只判断是否包含某种特点错误类型，所以 `errors.Is`和`errors.As` 方法就出现了。

## 带上函数调用栈信息

标准库的错误处理基本上能我们日常的开发需求，而且基本上能做到很优雅的错误处理。但是有时候我们还想带上更多信息，比如函数调用栈。我们使用第三方库[pingcap errors](https://github.com/pingcap/errors)来实现。

```golang
package main
import (
	"fmt"
	pkgerr "github.com/pingcap/errors"
)
type Err struct {
	Code int
	Msg string
}
func (e *Err) Error() string  {
	return fmt.Sprintf("code : %d ,msg:%s",e.Code,e.Msg)
}
var A_ERR = &Err{-1,"error"}

func stackfn1() error {
	return  pkgerr.WithStack(A_ERR)
}

func main()  {
	err := stackfn1()
	fmt.Printf("%+v",err) //这边使用 “%+v”
}
```

```bash
$ go run main.go
code : -1 ,msg:error
main.stackfn1
	/home/wida/gocode/goerrors-demo/main.go:18
main.main
	/home/wida/gocode/goerrors-demo/main.go:22
runtime.main
	/home/wida/go/src/runtime/proc.go:203
runtime.goexit
	/home/wida/go/src/runtime/asm_amd64.s:1357
Process finished with exit code 0
```
有了函数调用栈信息，我们可以更好的定位错误。