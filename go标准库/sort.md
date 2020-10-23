# sort——golang排序操作

有时候我们需要对一些结果集进行排序，有序的数据集适合我们查找元素。
golang的`sort`包默认提供了对`[]int`、`[]float64`和`[]string`排序的支持。还有提供一个`sort.Interface`排序接口，函数`sort.Sort`将对实现`sort.Interface`的数据类型举行排序。

## []int排序

```go
s := []int{5, 2, 6, 3, 1, 4} 
sort.Ints(s)
fmt.Println(s)  //[1 2 3 4 5 6]
```

## []float64排序

```go
s := []float64{5.2, -1.3, 0.7, -3.8, 2.6}
sort.Float64s(s)
fmt.Println(s) //[-3.8 -1.3 0.7 2.6 5.2]
```


## []strings排序

```go
s := []string{"Go", "Bravo", "Gopher", "Alpha", "Grin", "Delta"}
sort.Strings(s)
fmt.Println(s) //[Alpha Bravo Delta Go Gopher Grin]
```

## `Search`在排序好的结果集查找元素

```go
a := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
x := 6
i := sort.Search(len(a), func(i int) bool { return a[i] >= x })
if i < len(a) && a[i] == x {
	fmt.Printf("found %d at index %d in %v\n", x, i, a) //found 6 at index 2 in [1 3 6 10 15 21 28 36 45 55]
}
```

## 自建类型实现`sort.Interface`然后排序

```go
type Float32Slice []float32

func (p Float32Slice) Len() int           { return len(p) }
func (p Float32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Float32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

s := []float32{5.3, 9.2, 6.0, 3.8, 1.1, 4.5} // unsorted
sort.Sort(sort.Reverse(Float32Slice(s)))
fmt.Println(s) //[9.2 6 5.3 4.5 3.8 1.1]
```


## `Reverse`可以实现逆序

```go
s := []int{5, 2, 6, 3, 1, 4} // unsorted
sort.Sort(sort.Reverse(sort.IntSlice(s)))
fmt.Println(s) //[6 5 4 3 2 1]
```

`sort.Reverse`函数包含了`reverse`结构体他继承排序类型的`sort.Interface`，但是修改了 `Less(i, j int) bool`的方法。

```go
func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
```