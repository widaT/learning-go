# strings包

字符串操作通常是一种高频操作，在golang中专门提供了`strings`包来做这种工作。


## Compare 方法

``比较两个字符串的字典序，返回值0表示a==b，-1表示a < b，1表示a > b.通常我们应该运算符==，<，>来比较连个字符串，代码更简介，意思更明细，效率也更高。

```go
func Compare(a, b string) int
```

```go
strings.Compare("a", "b") //-1
strings.Compare("a", "a") //0
strings.Compare("b", "a") //1
```

## Contains 方法

`Contains`方法用来判断字符串`s`中是否包含子串`substr`

```go
func Contains(s, substr string) bool
```

```go
strings.Contains("seafood", "foo") //true
```

## HasPrefix和HasSuffix方法

`HasPrefix`和`HasSuffix`用来判断字符串中是否包含包含前缀和后缀

```go
func HasPrefix(s, prefix string) bool
func HasSuffix(s, suffix string) bool
```

```go
strings.HasPrefix("Gopher", "Go") //true
strings.HasSuffix("Amigo", "go") //true
```

## Index 方法和 LastIndex方法

`Index`返回子串`sep`在字符串`s`中第一次出现的位置，如果找不到，则返回 -1。
`LastIndex`返回子串`sep`在字符串`s`中最后一次出现的位置，如果找不到，则返回 -1。

```go
func Index(s, substr string) int
func LastIndex(s, substr string) int
```

```go
strings.Index("chicken", "ken") //4
strings.Index("chicken", "dmr") //-1

strings.Index("go gopher", "go") //0
strings.LastIndex("go gopher", "go") //3
strings.LastIndex("go gopher", "rodent") //-1
```

## Join方法

`Join`将`a`中的子串连接成一个单独的字符串，子串之间用`sep`拼接

```go
func Join(elems []string, sep string) string
```

```go
s := []string{"foo", "bar", "baz"}
fmt.Println(strings.Join(s, ", "))  //foo, bar, baz
```

## Repeat方法

`Repeat`将`count`个字符串`s`连接成一个新的字符串。

```go
func Repeat(s string, count int) string
```

```go 
strings.Repeat("na", 2)//nana
```

## Replace和ReplaceAll方法

`Replace`和`ReplaceAll`方法为字符串替换方法，`Replace`返回`s`的副本，并将副本中的`old`字符串替换为`new`字符串，替换次数为`n`次，如果`n`为 -1则全部替换。

```go
func Replace(s, old, new string, n int) string
func ReplaceAll(s, old, new string) string
```

```go
strings.Replace("oink oink oink", "k", "ky", 2)  //oinky oinky oink
strings.Replace("oink oink oink", "oink", "moo", -1) // moo moo moo
strings.ReplaceAll("oink oink oink", "oink", "moo") // moo moo moo 功能同上
```

## Split 方法

`Split`方法以`sep`为分隔符，将`s`切分成多个子串，结果中不包含 sep 本身

```go
func Split(s, sep string) []string
```

```go
strings.Split("a,b,c", ",") //["a" "b" "c"]
strings.Split("a man a plan a canal panama", "a ") //["" "man " "plan " "canal panama"]
strings.Split(" xyz ", "") //[" " "x" "y" "z" " "]
strings.Split("", "Bernardo O'Higgins") //[""]
```


## Trim、TrimSpace、TrimPrefix，TrimSuffix方法

`Trim`将删除`s`首尾连续的包含在`cutset`中的字符
`TrimSpace`将删除`s`首尾连续的的空白字符
`TrimPrefix`删除`s`头部的`prefix`字符串
`TrimSuffix` 删除`s`尾部的`suffix`字符串

```go
func Trim(s, cutset string) string
func TrimSpace(s string) string
func TrimPrefix(s, prefix string) string
func TrimSuffix(s, suffix string) string
```

```go
strings.Trim("¡¡¡Hello, Gophers!!!", "!¡") //Hello, Gophers
strings.TrimSpace(" \t\n Hello, Gophers \n\t\r\n")//Hello, Gophers

var s = "¡¡¡Hello, Gophers!!!"
strings.TrimPrefix(s, "¡¡¡Hello, ")  //Gophers!!!
strings.TrimSuffix(s, ", Gophers!!!") //¡¡¡Hello
```

# 参考资料

[go strings官方文档](https://golang.org/pkg/strings/)