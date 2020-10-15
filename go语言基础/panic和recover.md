# panic和recover

有些错误编译期就能发现，编译的时候编译器就会报错，而有些运行时的错误编译器是没办法发现的。golang的运行时异常叫做`panic`。

在golang当切片访问越界，空指针引用等会引起`panic`，`panic`如果没有自己手动捕获的话，程序会中断运行并打印`panic`信息。

```go
var a []int = []int{0, 1, 2}
fmt.Println(a[3]) //panic: runtime error: index out of range [3] with length 3
```

```go
type A struct {
		Name string
	}

var a *A
fmt.Println(a.Name) //panic: runtime error: invalid memory address or nil pointer dereference
```

## 手动panic

当程序运行中出现一些异常我们需要程序中断的时候我们可以手动`panic`

```go
var model = os.Getenv("MODEL") //获取环境变量 MODEL
if model == "" {
    panic("no value for $MODEL") //MODEL 环境变量未设置程序中断
}
```

## 使用recover捕获panic

我们在`defer`修饰的函数里面使用`recover`可以捕获的`panic`。

```go
func main() {
	func() {
		defer func() {
			if err := recover(); err != nil { //recover 函数没有异常的是返回 nil
				fmt.Printf("panic: %v", err)
			}
		}()
		func() {
			panic("panic") //函数执行出现异常
		}()
	}()
}
```

`panic`的产生后会终止当前函数运行，然后去检测当前函数的`defer`是否有`recover`，没有的话会一直往上层冒泡直至最顶层；如果中间某个函数的defer有`recover`则这个向上冒泡过程到这个函数就会终止。
