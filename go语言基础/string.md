# 字符串

## 字符串的本质
go内嵌`string`类型来表示字符串。我们来看`string`的本质

```
type StringHeader struct {
	Data uintptr
	Len  int
}
```

本质是`string`是个结构体，有连个字段`Data`是个指针指向一段字节系列（连续内存）开始的位置。`Len`则代表长度，内嵌函数`len（s）`可以获取这个长度值。


`Data`字段指向的是一段连续内存的起始文字，这一段内存的值是不允许改变的。

```go 
a:="aaa"  // a{Data:&"aaa",Len:3}  &"aaa"代表为“aaa”内存起始地址	
b:=a     // b{Data:&"aaa",Len:3}
c:=a[2:] // c{Data:&"aaa"+2,Len:3}
```


## 字符串拼接

go字符串拼接使用`+`

```go

a：="hello" + " world"

b := a+"ccc"

```

## 字符串切片操作

字符串支持下标索引的方式索引到字符，`s[i]` i需要满足`0<=i<len(s)`,小于0会引起编译错误，大于等于`len(s)`会引起运行时错误`panic: index out of range`。

```go
s := "hello world"
fmt.Println(len(s))     //11
fmt.Println(s[2], s[7]) //108 111（'l'，'o'）
```

`s[i]`的本质是获取 `Data`指针指向内容第`i`个字节的值。由于`Data`指向的只是不允许改变的，所以试图改变它的值会引起编译错误

```
s[2] ='d' // compile error
```

字符串支持切片操作来产生新的字符串

```go

s := "hello world"

fmt.Println(s[0:5]) //"hello"
fmt.Println(s[:5]) //"hello" 等同 s[0:5]
fmt.Println(s[5:])  //" world"  等同 s[5:11]
fmt.Println(s[:]) //"hello world" 等同 s[0:11]
```

## rune

go语言中使用`rune`处理`Unicode`，go的源码文件使用的是UTF8编码。go语言使用`unicode/utf8`来处理utf8。

本质上`rune`是`int32`的别名。

```go
import "unicode/utf8"

s := "Hello,世界"
fmt.Println(len(s))                    // "12"
fmt.Println(utf8.RuneCountInString(s)) // "8"

```
需要主要字符串在用range变量的时候 会转换成`[]rune`
```
for _, v := range s {
		fmt.Println(string(v))
}
```
结果是
```
H
e
l
l
o
,
世
界
```
