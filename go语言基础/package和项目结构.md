# package和go项目结构

## go源码文件

go项目的源码文件有如下三种

- `.go`为文件名的go代码文件，在实际项目中这个是最常见的。
- `.s`为结尾的go汇编代码文件，需要汇编加速的项目会用到，go内核源码比较多，实际项目中非常少见。
- `.h .cpp .c`为结尾的c/c++代码，需要和c/c++交互的项目有看到，在实际项目也不较少见。

一般go编译器在编译时会扫描相关目录下的这些go源文件进行编译.

如下情况下go源码文件会忽略编译
-  `_test.go`为结尾的文件,这些文件是go的单元测试代码,用于测试go程序,不参与编译。
- `_$GOOS_$GOARCH.go` `$GOOS`代码操作系统（windows，linux等），`$GOARCH`代表cpu架构平台（arm64，adm64等），go编译时会符合环境变量相关文件编译。
- 编译过程中指定了条件编译参数`go build -tags tag.list`,编译器会选择性忽略一些文件不参与编译,参考[go build](../go开发环境搭建/go命令.md)。


## 包

包（package）是go语言模块化的实现方式，每个go代码需要在第一行非注释的代码注明package。

go的package是基于目录的，同一个目录下的go代码只能用相同的package。

## main package

go中`main package`是go可执行程序必须包含的一个package。该package下的`main`方法是go语言程序执行入口。


## go可见性规则

golang中常量、变量、类型名称、函数名、结构字段名 以一个大写字母开头，这个对象就可以被外部包的代码所使用；标识符如果以小写字母开头，则对包外是不可见的，整个包的内部是可见。


```go

//a目录下
package a

const COST="aa" //对外可见

var A :=0  //对外可见
var a :=0  //对外不可见

func F1() { //对外可见

}

func f1(){ //对外不可见

}

type Student struct{ //对外可见
    Name string   //对外可见
    age int       //对外不可见
} 

type student struct {//对外不可见

}


//main package

package main

import "a"

func main(){

   println(a.COST) //访问a package下的常量

    a.A ==1 //访问a package下 A变量
    a.a //编译报错
    a.F1()
    a.f1()//编译报错


    var stu  a.Student   //访问 a package下的Student结构体

    var stu1 a.student //student不可见 编译报错

    stu.Name = "nam1" //ok
    stu.age = 18 //不可见 编译报错
}

```