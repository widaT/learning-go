# go 语言结构图优雅初始化的一种方式
该方式存于go-micro，go-grpc中

```go
package  main
import "fmt"
type Config struct {
	A string
	B int
}

type option func(c *Config)

func NewConig(o ...option)  *Config{
	c := Config{
		A:"wida",
		B:3,
	}
	for _,op:= range o {
		op(&c)
	}
	return &c
}

func writeA(a string)  option{
	return func(c *Config) {
		c.A=a
	}
}

func writeB(b int) option{
	return func(c *Config) {
		c.B = b
	}
}

func main()  {
	a := writeA("amy")
	b := writeB(6)
	c := NewConig(a,b)
	fmt.Println(c)
}
```

运行：
```bash
# go run main.go
&{amy 6}
```