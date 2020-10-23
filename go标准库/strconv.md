# strconv

`strconv`包提供了字符串和其他golang基础类型的互相转换函数。

## 整数和字符串互换

```go
i, err := strconv.Atoi("-42") // -42
s := strconv.Itoa(-42)  //"-42"
```

## 布尔类型，浮点数，整数和字符串转换

```go
b, err := strconv.ParseBool("true")
f, err := strconv.ParseFloat("3.1415", 64)
i, err := strconv.ParseInt("-42", 10, 64)
u, err := strconv.ParseUint("42", 10, 64)

s := strconv.FormatBool(true)
s := strconv.FormatFloat(3.1415, 'E', -1, 64)
s := strconv.FormatInt(-42, 16)
s := strconv.FormatUint(42, 16)
```

## Quote 和 Unquote

`strconv`包中有对字符串加`"`的方法`Quote`，`Unquote`这是给字符串去`"`.

```go
fmt.Println(strconv.Quote(`"Hello   世界"`))         //"\"Hello\t世界\""
fmt.Println(strconv.QuoteRune('世'))                // '世'
fmt.Println(strconv.QuoteRuneToASCII('世'))         // '\u4e16'
fmt.Println(strconv.QuoteToASCII(`"Hello    世界"`)) //"\"Hello\t\u4e16\u754c\""
fmt.Println(strconv.Unquote(`"\"Hello\t世界\""`))    // "Hello    世界" <nil>
```