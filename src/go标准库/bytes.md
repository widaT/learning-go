# bytes包

`[]byte`字节数组操作是`string`外的另外一种高频操作，在golang中也专门提供了`bytess`包来做这种工作。功能大多数和`strings`类似。


## Compare和Equal方法

`Compare`方法安装字典顺序比较`a`和`b`，返回值1为a>b，0为a==b，-1为a<b。
`Equal`方法判断a和b的长度一致，并且a的b的byte值是一样的则为true，否则为false。

```go
func Compare(a, b []byte) int
func Equal(a, b []byte) bool
```

```go
bytes.Compare([]byte("a"), []byte("b")) //-1
bytes.Compare([]byte("a"), []byte("a")) //0
bytes.Compare([]byte("b"), []byte("a")) //-1
```


## Contains 方法

`Contains`方法用来字节数组`b`中是否包含子字节数组`subslice`

```go
func Contains(b, subslice []byte) bool
```

```go
bytes.Contains([]byte("seafood"), []byte("foo")) //true
bytes.Contains([]byte("seafood"), []byte("bar")) //false
```


## HasPrefix和HasSuffix方法

`HasPrefix`和`HasSuffix`用来判断字节数组中是否包含前缀和后缀

```go
func HasPrefix(s, prefix []]byte)) bool
func HasSuffix(s, suffix []]byte)) bool
```

```go
bytes.HasPrefix([]byte("Gopher"), []byte("Go")) //true
bytes.HasSuffix([]byte("Amigo"), []byte("go")) //true
```

## Index 方法和 LastIndex方法

`Index`返回子字节数组`sep`在字节数组`s`中第一次出现的位置，如果找不到，则返回 -1。
`LastIndex`返回子字节数组`sep`在字节数组`s`中最后一次出现的位置，如果找不到，则返回 -1。

```go
func Index(s, sep []byte) int
func LastIndex(s, sep []byte) int
```

```go
bytes.Index([]byte("chicken"), []byte("ken")) //4
bytes.Index([]byte("chicken"), []byte("dmr")) //-1

bytes.Index([]byte("go gopher"), []byte("go")) //0
bytes.LastIndex([]byte("go gopher"), []byte("go")) //3
bytes.LastIndex([]byte("go gopher"), []byte("rodent")) //-1
```


## Join方法

`Join`将`a`中的子串连接成一个单独的字符串，子串之间用`sep`拼接

```go
func Join(s [][]byte, sep []byte) []byte
```

```go
s := [][]byte{[]byte("foo"), []byte("bar"), []byte("baz")}
fmt.Printf("%s", bytes.Join(s, []byte(", "))) //foo, bar, baz
```


## Repeat方法

`Repeat`将`count`个字节数组`b`连接成一个新的字节数组。

```go
func Repeat(b []byte, count int) []byte
```

```go 
bytes.Repeat([]byte("na"), 2)//nana
```

## Replace和ReplaceAll方法

`Replace`和`ReplaceAll`方法为字节数组替换方法，`Replace`返回`s`的副本，并将副本中的`old`字符串替换为`new`字节数组，替换次数为`n`次，如果`n`为 -1则全部替换。

```go
func Replace(s, old, new []byte, n int) []byte
func ReplaceAll(s, old, new []byte) []byte
```

```go
bytes.Replace([]byte("oink oink oink"), []byte("k"), []byte("ky"), 2)  //oinky oinky oink
bytes.Replace([]byte("oink oink oink"), []byte("oink"), []byte("moo"), -1) // moo moo moo
bytes.ReplaceAll([]byte("oink oink oink"), []byte("oink"), []byte("moo")) // moo moo moo 功能同上
```

## Split 方法

`Split`方法以`sep`为分隔符，将`s`切分成多个子字节数组，结果中不包含 sep 本身

```go
func Split(s, sep []byte) [][]byte
```

```go
bytes.Split([]byte("a,b,c"), []byte(",")) //["a" "b" "c"]
bytes.Split([]byte("a man a plan a canal panama"), []byte("a ")) //["" "man " "plan " "canal panama"]
bytes.Split([]byte(" xyz "), []byte("")) //[" " "x" "y" "z" " "]
bytes.Split([]byte(""), []byte("Bernardo O'Higgins")) //[""]
```


## Trim、TrimSpace、TrimPrefix，TrimSuffix方法

`Trim`将删除`s`首尾连续的包含在`cutset`中的字符
`TrimSpace`将删除`s`首尾连续的的空白字符
`TrimPrefix`删除`s`头部的`prefix`字符串
`TrimSuffix` 删除`s`尾部的`suffix`字符串

```go
func Trim(s []byte, cutset string) []byte
func TrimSpace(s []byte) []byte
func TrimPrefix(s, prefix []byte) []byte
func TrimSuffix(s, suffix []byte) []byte
```

```go
bytes.Trim([]byte(" !!! Achtung! Achtung! !!! "), "! ") //Achtung! Achtung
bytes.TrimSpace([]byte(" \t\n Hello, Gophers \n\t\r\n"))//Hello, Gophers

var s = []byte("¡¡¡Hello, Gophers!!!")
bytes.TrimPrefix(s, []byte("¡¡¡Hello, "))  //Gophers!!!
bytes.TrimSuffix(s, []byte(", Gophers!!!")) //¡¡¡Hello
```

# 参考资料

[go bytes官方文档](https://golang.org/pkg/bytes/)