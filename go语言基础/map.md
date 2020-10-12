# Map

map结果通常是现代语言中高频使用的结构。了解map的使用方法和map的底层解构是十分必要的，本小节将介绍map的使用，
map的底层原来健在go runtime的章节介绍。

## Map 变量声明

```go
var m map[keytype]valuetype
```
map的key是有要求的，key必须是可比较类型，比如常见的 `string`，`int`类型。而数组，切片，和结构体不能作为key类型。

map的value可以是任何类型。

map的默认初始值为 `nil`


## 创建map和添加元素

golang中创建`map`需要使用`make`方法

格式如下

```go
m := make(map[keytype]valuetype, cap)
```

`cap`参数可以省略。map的容量是根据实际需要自动伸缩。


内置函数`len`可以获取map的长度


map make之后就可以添加元素了

```go
m:=make(map[string]int)
m["aa"]=1
m["bb"]=2
```


## 判断key是否存在

```go
m:=make(map[string]int)
m["aa"]=1
m["bb"]=2

v, found := m["aa"] //found为true 为1
v, found := m["aaa"] //found为false 为0

//通常和if表达是一起，判断key是否存在
if vv,found:=m["bb"];found { //key存在
    //code
}
```


## 删除map中的元素

map使用内置函数`delete`来删除元素。

```go
delete(m,key)
```

```go
m:=make(map[string]int)
m["aa"]=1

delete(m,"aa") 
delete(m,"bb") // key不存在也不会报错 
```


## 遍历map

golang 使用`for-range`遍历map

```go
for k, v := range m {
   //code
}
```

k,v分别代表key和这个key对于的值。

如果只想获取key可以

```go
for k := range m {
    fmt.Println(key)
}
```

如果只想获取value则：

```go
for _, v := range m {
    fmt.Println(v)
}
```

## map的一些特性

- map中key是不排序的，所以你每一次对map进行`for-range`打印的结果可能不同
- golang的map不是`线程安全`的，在并发读写的时候会触发`concurrent map read and map write`的panic。解决的办法是自己加把锁。
- map在删除某个key的时候不是真正的删除，只是标记为空，空间还在，所以在for range删除map元素是安全的。