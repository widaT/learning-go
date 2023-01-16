# Go语言基础

## go的25个语言关键字

```
break      default       func     interface   select
case       defer         go       map         struct
chan       else          goto     package     switch
const      fallthrough   if       range       type
continue   for           import   return      var
```

这个25个关键字，不能自定义使用。变量和函数应避免和上面的25个重名。

## 预定的36个符号

### 内建常量

    true false iota nil

### 内建类型
    
    int int8 int16 int32 int64
    uint uint8 uint16 uint32 uint64 uintptr
    float32 float64 complex128 complex64
    bool byte rune string error

### 内建函数

    make len cap new append copy close delete
    complex real imag
    panic recover


## go语言命名规范

Go程序员应该使用驼峰命名，当名字由几个单词组成时使用大小写分隔，而不是用下划线分隔。例如“NewObject”或者“newObject”，而非“new_object”.


