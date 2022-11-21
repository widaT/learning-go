# log —— 官方的日志库

golang标准库提供了一个`log`包用来实现简单的程序日志记录功能。
`log`包的使用非常简单，函数名字和用法也和fmt包很相似，只是在它的输出默认带了时间。

## 三个基础函数

```go
func Print(v ...interface{})
func Printf(format string, v ...interface{})
func Println(v ...interface{})
```

```go
log.Print("print：", "这是Printf：产生的日志", "\n")
log.Println("Println:", "这是Println产生的日志")
log.Printf("Printf：%s\n", "这是Printf：产生的日志")
```

```bash
$ go run main.go
2020/10/23 10:01:35 print：这是Printf：产生的日志
2020/10/23 10:01:35 Println: 这是Println产生的日志
2020/10/23 10:01:35 Printf：这是Printf：产生的日志
```

## 打印日志后产生`Panic`

```go
func Panic(v ...interface{})   //功能和`Print()`一样，只是后面加了`panic()`.
func Panicf(format string, v ...interface{}) //功能和`Printf()`一样，只是后面加了`panic()`.
func Panicln(v ...interface{}) //功能和` Println()`一样，只是后面加了`panic()`.
```

## 打印日志后产生后退出程序

```go
func Fatal(v ...interface{}) //功能和`Print()`一样，只是后面加了`os.Exit(1)`.
func Fatalf(format string, v ...interface{})   //功能和`Printf()`一样，只是后面加了`os.Exit(1)`.
func Fatalln(v ...interface{})    //功能和` Println()`一样，只是后面加了`os.Exit(1)`.
```


## 修改日志输出格式

`SetFlags`可以改变日志输出的格式，主要改变的是日期时间格式和文件行号格式。

```go
func SetFlags(flag int)
```

```go
Ldate         = 1 << iota     //时间只包含日期： 2009/01/23
Ltime                         //时间只包含时分秒: 01:23:23
Lmicroseconds                 //时间包含时分秒毫秒: 01:23:23.123123.
Llongfile                     //日志产生的代码文件绝对路径和行号: /a/b/c/d.go:23
Lshortfile                    //日志产生的代码文件和行号: d.go:23.
LUTC                          //日期时间转为0时区的
LstdFlags     = Ldate | Ltime //默认值
```

```go
log.SetFlags(log.Ldate | log.Lshortfile)
log.Print("print：", "这是Printf：产生的日志", "\n")
log.Println("Println:", "这是Println产生的日志")
log.Printf("Printf：%s\n", "这是Printf：产生的日志")
```

```bash
$ go run main.go
2020/10/23 test.go:7: print：这是Printf：产生的日志
2020/10/23 test.go:8: Println: 这是Println产生的日志
2020/10/23 test.go:9: Printf：这是Printf：产生的日志
```

## 添加日志前缀

我们还可以给每一行日志加一个前缀。

```go
func SetPrefix(prefix string)
```

```go
log.SetPrefix("[DEBUG] ")
log.Print("print：", "这是Printf：产生的日志", "\n")
```

```bash
$ go run main.go
[DEBUG] 2020/10/23 test.go:10: print：这是Printf：产生的日志
```

## 输出到文件

默认情况下日志会输出到`标准输出`，我们可以使用`SetOutput`修改输出方式。

```go
func SetOutput(w io.Writer)
```

```go
f, _ := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

log.SetOutput(f)
log.Print("print：", "这是Printf：产生的日志", "\n")
log.Println("Println:", "这是Println产生的日志")
log.Printf("Printf：%s\n", "这是Printf：产生的日志")
```    

这样子日志将被写到`log.log`文件中。


# 后记

这个`log`包的功能相对比较单一，缺少日志分级，日志文件切割，日志文件大小和个数控制等功能，通常在实际项目中我们会使用更加强大的第三方包来使用。

