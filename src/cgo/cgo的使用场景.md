# cgo的使用场景

在本章的开头，我们建议大家非必要情况不要用cgo。本小节我们详细来看什么时候用cgo。


## 场景1——提升算法效率

本例中要求计算512维的向量的欧式距离。我们用go原生实现了算法，然后使用c/c++平台的avx（Advanced Vector Extensions 高级向量拓展集）实现同样的算法，同时比对下效率。
```go
package main
/*
#cgo CFLAGS: -mavx -std=c99
#include <immintrin.h> //AVX: -mavx
float avx_euclidean_distance(const size_t n, float *x, float *y)
{
    __m256 vsub,vsum={0},v1,v2;
    for(size_t i=0; i < n; i=i+8) {
        v1  = _mm256_loadu_ps(x+i);
        v2  = _mm256_loadu_ps(y+i);
        vsub = _mm256_sub_ps(v1, v2);
        vsum = _mm256_add_ps(vsum, _mm256_mul_ps(vsub, vsub));
    }
    __attribute__((aligned(32))) float t[8] = {0};
    _mm256_store_ps(t, vsum);
    return t[0] + t[1] + t[2] + t[3] + t[4] + t[5] + t[6] + t[7];
}
*/
import "C"

import (
	"fmt"
	"time"
)

func euclideanDistance(size int,x, y []float32) float32 { //cgo实现欧式距离
	dot := C.avx_euclidean_distance((C.size_t)(size), (*C.float)(&x[0]), (*C.float)(&y[0]))
	return float32(dot)
} 
func euclidean(infoA, infoB []float32) float32 { //go原生实现欧式距离
	var distance float32
	for i, number := range infoA {
		a := number - infoB[i]
		distance += a * a
	}
	return distance
}
func main()  {
	size := 512
	x := make([]float32, size)
	y := make([]float32, size)
	for i := 0; i < size; i++ {
		x[i] = float32(i)
		y[i] = float32(i + 1)
	}
	
	stime := time.Now()
	times := 1000
	for i:=0;i<times;i++ {
		euclidean(x,y)
	}
	fmt.Println(time.Now().Sub(stime))

	stime = time.Now()
	for i:=0;i<times;i++ {
		euclideanDistance(size,x,y)
	}
	fmt.Println(time.Now().Sub(stime))
}
```
运行一下
```
$ go run main.go
447.729µs
143.182µs
```
cgo实现的比go原生的算法效率大概快了四倍。其实go实现欧式距离算法还可以更快，后续在go汇编会介绍绕过cgo直接使用汇编实现欧式距离计算的方式。

## 场景2——依赖第三方sdk

在实际开发过程中，我们经常会遇到一些第三方sdk，这些sdk很多为了保护源码会用c和c++编写，然后给你sdk头文件和一个动态或者静态库文件。这个时候只好使用cgo实现自己的业务。

还有就是在c/c++生态产生了很多优秀的项目，比如rocksdb，就是c++实现的本地LSM-Tree实现的内嵌型kv数据库。
很多NewSql的底层数据存储都用rocksdb。所以rocksdb被很多语言集成，有java，python，当然还有go。[gorocksdb](https://github.com/tecbot/gorocksdb)就是用cgo实现的go rocksdb集成。



## 总结

本小节介绍了两个使用cgo的场景，这两个场景通常比较常见，淡然除了这两个场景还有一些不常见的场景比如go访问系统驱动，这时候通常也会用cgo实现。总的来说使用cgo需要谨慎些。
