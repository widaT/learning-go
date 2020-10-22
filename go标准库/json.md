# json包

现代的语言都会支持json这种轻量的序列化方式。在golang中使用`encoding/json`来支持json的各个操作。


## 序列化方法

json序列化方法是`Marshal`函数，函数签名是

```go
func Marshal(v interface{}) ([]byte, error)
```

## 反序列化方法

json序列化方法是`Unmarshal`函数，函数签名是

```go
func Unmarshal(data []byte, v interface{}) error
```

## 用于json的Struct Tag

json序列号和反序列化只针对可见字段，对于不可见字段这两个过程都会直接忽略。
使用`omitempty`来处理空值，使用`-`忽略字段。

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type Stu struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		No     string `json:"no,omitempty"`
		Gender int    `json:"-"`
		group  int    //这个字段不参与序列化
	}

	stu := Stu{
		Name:   "wida",
		Age:    35,
		No:     "8001",
		Gender: 0,
		group:  88,
	}

	b, _ := json.Marshal(&stu)
	fmt.Println(string(b))

	var stu1 Stu

	json.Unmarshal(b, &stu1)
	fmt.Printf("%+v", stu1)
}
```

```bash
$ go run main.go
{"name":"wida","age":35,"no":"8001"}
{Name:wida Age:35 No:8001 Gender:0 group:0}
```

## 定制序列化和反序列化方式

在`encoding/json`包中定义了`Marshaler`和`Unmarshaler`
```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}
type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```
如果一个类型实现了这两个接口就有实现自定义的json序列化和反序列化。

看下官方文档的一个例子：

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Animal int

const (
	Unknown Animal = iota
	Gopher
	Zebra
)

func (a *Animal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*a = Unknown
	case "gopher":
		*a = Gopher
	case "zebra":
		*a = Zebra
	}

	return nil
}

func (a Animal) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	default:
		s = "unknown"
	case Gopher:
		s = "gopher"
	case Zebra:
		s = "zebra"
	}

	return json.Marshal(s)
}

func main() {
	blob := `["gopher","armadillo","zebra","unknown","gopher","bee","gopher","zebra"]`
	var zoo []Animal
	if err := json.Unmarshal([]byte(blob), &zoo); err != nil {
		log.Fatal(err)
	}

	census := make(map[Animal]int)
	for _, animal := range zoo {
		census[animal] += 1
	}

	fmt.Printf("Zoo Census:\n* Gophers: %d\n* Zebras:  %d\n* Unknown: %d\n",
		census[Gopher], census[Zebra], census[Unknown])

}
```