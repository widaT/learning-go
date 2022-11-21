# errgroup

## 什么是errrgroup

在开发并发程序时，错误的收集和传播往往比较繁琐，有时候当一个错误发声时，我们需要停止所有相关任务，有时候却不是。`sync.ErrGroup`刚好可以解决我们上述的痛点，它提供错误传播，以及利用`context`的方式来决定是否要停止相关任务。

`errrgroup.Group`结构体
```go
type Group struct {
	cancel func()
	wg sync.WaitGroup
	errOnce sync.Once
	err     error
}
``

三个对外api

```go
unc WithContext(ctx context.Context) (*Group, context.Context)
func (g *Group) Go(f func() error)
func (g *Group) Wait() error
```

## 使用errrgroup

### 只返回错误

```go
package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func main() {
	var g errgroup.Group
	var urls = []string{
		"http://www.golang.org/",
		"http://www.111111111111111111111111.com/", //这个地址不存在
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for _, url := range urls {
		url := url 
		g.Go(func() error {
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}else {
		fmt.Println(err)
	}
}
```

```bash
$ go run main.go
www.111111111111111111111111.com:80: unknown error host unreachable
```

### 使用 errgroup.WithContext

```go
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

func main() {
	ctx ,_:=context.WithTimeout(context.Background(),3*time.Second)
	var g ,_=  errgroup.WithContext(ctx)
	var urls = []string{
		"http://www.golang.org/",
		"http://www.111111111111111111111111.com/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for _, url := range urls {
		url := url
		g.Go(func() error {
			ch := make(chan error)
			go func() {
				time.Sleep(4e9)
				resp, err := http.Get(url)
				if err == nil {
					resp.Body.Close()
				}
				ch <- err
			}()
			select {
            case err:= <-ch :
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}else {
		fmt.Println(err)
	}
}
```

```bash
$ go run main.go
context deadline exceeded
```