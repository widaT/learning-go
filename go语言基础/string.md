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

## 


## rune




## 