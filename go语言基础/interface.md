# 接口(interface)

golang中接口（interface)是被非常精心设计的，利用接口golang可以实现很多类似面向对象的设计，而且比传统的面向对象更加方便。

接口是一组方法的定义的集合，接口中没有变量，。

接口的格式：

```go
type 接口名 interface {
    M1(参数列表) （返回值列表） //方法1
    M2(参数列表) （返回值列表） //方法2
    ...
}
```

某个类型`全部实现接口中定义的方法`，我们称作该类型实现了该接口，而不需要在代码中显示声明该类型实现某个接口。

```go
type Sayer interface{
    Say()
}

type Dog struct{}
type Cat struct{}

func (d *Dog) Say() {
	fmt.Println("wang ~")
}

func (Cat *Cat) Say() {
	fmt.Println("miao ~")
}

func Say(s Sayer) { //函数使用接口类型做参数
	s.Say()
}

var sayer Sayer

dog:=new(Dog)
cat:=new(Cat)

sayer = dog  //这样赋值给接口变量是ok的，Dog实现了Say方法
sayer.Say()

sayer = cat
sayer.Say()

Say(dog)  //这样子传参也ok
Say(cat)
```

从上面的例子上看，golang的`interface`可以实现面向对象语言多态的特性，而且更加简洁，高效。


## 接口嵌套

golang的接口是可以嵌套的，一个接口可以嵌套一个或者多个其他接口。

```go

type Reader interface{
    Read()
}

type Writer interface {
    Write()
} 


type ReadWrite interface { //ReadWrite 嵌套了 Reader 和 Writer
    Reader
    Writer
}
```

## 空接口

在接口中，空接口有其独特的位置，空接口没有定义任何方法，那么在golang中任何类型都实现了这个接口。

```go
interface{}
```

我们再看下`fmt`包中`Println`的定义：
```go
//func Println(a ...interface{}) (n int, err error)
fmt.Println(1,"aa",true) //由于参数是 interface{} 所以可以传任意类型
```

## 类型断言

接口类型的变量支持类型断言，通过类型断言我们可以检测接口类型变量底层真实的类型。

类型断言表达式如下：

```go
v，ok := varInterface.(T) 
```


```go
var i interface{}
a := 1
i = a

if temp, ok := i.(int); ok { //实际开发中我们经常会这样子 做类型转换
	fmt.Println(temp)
}


dog:=new(Dog)
_,ok := dog.(Sayer) //我们还可以判某个类似是否实现了某个接口
```

我们可以使用`switch case`语句来更多规模的类型检测

比如`fmt`包中`printArg`的这一段代码

```go
switch f := arg.(type) {
	case bool:
		p.fmtBool(f, verb)
	case float32:
		p.fmtFloat(float64(f), 32, verb)
	case float64:
		p.fmtFloat(f, 64, verb)
    case complex64:
        p.fmtComplex(complex128(f), 64, verb)
    ...
```