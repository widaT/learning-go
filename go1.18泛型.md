# go 1.18泛型初体验

go.1.18beta版发布，众所周知这个go.1.18版本会默认启用go泛型。这个版本也号称go最大的版本改动。

## 初识golang的泛型

我们写一个demo来看看go的泛型是长啥样

```go

package main

import (
	"fmt"
)

type OrderTypes interface {
	~int | ~float32 | ~string
}

func max[T OrderTypes](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func main() {
	fmt.Println(max(1, 11), max("abc", "eff"))
}

```
ok run 一下代码
```bash
$ go run main.go
11 eff
```

`~int | ~float32 | ~string`我们看到了新的语法，`~`是新的操作符，主要用来做类型约束使用， `~int`代表类型约束为`int`类型,`~int | ~float32 | ~string`则代表约束为 `int` 或者 `float32` 或者 `string`。上面额例子中，这三个类型刚好是可以比较的能进行 ">" 操作的。

当然上面的代码是演示用的，在真正的项目中我们应该使用标准`constraints`提供的`Ordered`来做约束。

```go
import (
	"constraints"
)
func max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}
```

`constraints`标准库定义了一下常用的类型约束，如`Ordered`,`Signed`,`Unsigned`,`Integer`,`Float`。



## 提高生产力的泛型

我们通过下面的例子来看看泛型，如何提高我们的生产力。我们将为所有`slice`类型添加三件套`map`,`reduce`,`filter`

```go
func Map[Elem1, Elem2 any](s []Elem1, f func(Elem1) Elem2) []Elem2 {
	r := make([]Elem2, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func Reduce[Elem1, Elem2 any](s []Elem1, initializer Elem2, f func(Elem2, Elem1) Elem2) Elem2 {
	r := initializer
	for _, v := range s {
		r = f(r, v)
	}
	return r
}

func Filter[Elem any](s []Elem, f func(Elem) bool) []Elem {
	var r []Elem
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func Silce() {
	sliceA := []int{3, 99, 31, 63}
	//通过sliceA 生成sliceB
	sliceB := Map(sliceA, func(e int) float32 {
		return float32(e) + 1.3
	})
	fmt.Println(sliceB)
	//找最大值
	max := Reduce(sliceB, 0.0, func(a, b float32) float32 {
		if a > b {
			return a
		}
		return b
	})
	fmt.Println(max)
	//过滤sliceA中大于30的组成新的slice
	sliceC := Filter(sliceA, func(e int) bool {
		if e > 30 {
			return true
		}
		return false
	})
	fmt.Println(sliceC)
}

func main() {
	Silce()
}
```

```bash
$ go run main.go 
[4.3 100.3 32.3 64.3]
100.3
[99 31 63]
```

## 总结

go的泛型目前还没有官方推荐的最佳实践，标准库的代码也基本没改成泛型。但总归走出支持泛型这一步，后续丰富标准库应该是后面版本的事情了。再看[go2](https://github.com/golang/go/blob/dev.go2go/src/cmd/go2go/testdata/go2path/src/)代码的时候发现一个有意思的东西--`orderedmap`。感兴趣的同学可以去看看。