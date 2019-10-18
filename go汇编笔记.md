# go 汇编笔记


# JLS 指令
```assembly
CMPQ CX，$3
JLS  48  
```

JLS(通JBE一个意思)转移条件：`CMPQ CX，$3 ；JLS 48` 当CX的值小于或等于3的时候转跳到48。


# 编译器优化的例子

go标准库`binary.BigEndian`中的`bigEndian.PutUint32`
```golang
func PutUint32(b []byte, v uint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}
```

这个例子 _ = b[3] 这个语句对于的汇编是
```bash
$ go tool compile -S main.go |grep main.go:10 
        0x000e 00014 (main.go:10)       PCDATA  $0, $0
        0x000e 00014 (main.go:10)       PCDATA  $1, $0
        0x000e 00014 (main.go:10)       MOVQ    "".b+40(SP), CX
        0x0013 00019 (main.go:10)       CMPQ    CX, $3
        0x0017 00023 (main.go:10)       JLS     48
        0x0030 00048 (main.go:10)       MOVL    $3, AX
        0x0035 00053 (main.go:10)       CALL    runtime.panicIndex(SB)
        0x003a 00058 (main.go:10)       XCHGL   AX, AX
```
意思很明显是提前判断一下是不是slice 下标越界。
但是这个 为什么不直接写成如下的样子，不是也会检查 panic，而不需要额外的 `_ = b[3]`
```golang
func PutUint32(b []byte, v uint32) {
	b[3] = byte(v)
	b[2] = byte(v >> 8)
	b[1] = byte(v >> 16)
	b[0] = byte(v >> 24)
}
```

其实这边就是涉及编译器优化的问题，我们看下原先 `PutUint32`的汇编去掉gc的代码

```
"".PutUint32 STEXT nosplit size=59 args=0x20 locals=0x18
        0x0000 00000 (main.go:9)        TEXT    "".PutUint32(SB), NOSPLIT|ABIInternal, $24-32
        0x0000 00000 (main.go:9)        SUBQ    $24, SP
        0x0004 00004 (main.go:9)        MOVQ    BP, 16(SP)
        0x0009 00009 (main.go:9)        LEAQ    16(SP), BP
        0x000e 00014 (main.go:10)       MOVQ    "".b+40(SP), CX //这边是回去slice len的值
        0x0013 00019 (main.go:10)       CMPQ    CX, $3
        0x0017 00023 (main.go:10)       JLS     48
        0x0019 00025 (main.go:11)       MOVL    "".v+56(SP), AX
        0x001d 00029 (main.go:11)       BSWAPL  AX
        0x001f 00031 (main.go:14)       MOVQ    "".b+32(SP), CX
        0x0024 00036 (main.go:14)       MOVL    AX, (CX)
        0x0026 00038 (main.go:15)       MOVQ    16(SP), BP
        0x002b 00043 (main.go:15)       ADDQ    $24, SP
        0x002f 00047 (main.go:15)       RET
        0x0030 00048 (main.go:10)       MOVL    $3, AX
        0x0035 00053 (main.go:10)       CALL    runtime.panicIndex(SB)
        0x003a 00058 (main.go:10)       XCHGL   AX, AX

```
从这段汇编来看，你会看到
```golang
	b[3] = byte(v)
	b[2] = byte(v >> 8)
	b[1] = byte(v >> 16)
    b[0] = byte(v >> 24)
```
这个已经直接被优化掉 
```assembly
MOVL    "".v+56(SP), AX
BSWAPL  AX          //指令作用是：32位寄存器内的字节次序变反。比如：(EAX)=9668 8368H，执行指令：BSWAP EAX ，则(EAX)=6883 6896H。
``` 
BSWAPL 指令做的工作就和上4个做的工作一样。

我们再看看乱序后的汇编代码

```bash
"".PutUint32 STEXT nosplit size=79 args=0x20 locals=0x18
        0x0000 00000 (main.go:9)        TEXT    "".PutUint32(SB), NOSPLIT|ABIInternal, $24-32
        0x0000 00000 (main.go:9)        SUBQ    $24, SP
        0x0004 00004 (main.go:9)        MOVQ    BP, 16(SP)
        0x0009 00009 (main.go:9)        LEAQ    16(SP), BP
        0x000e 00014 (main.go:10)       MOVQ    "".b+40(SP), CX
        0x0013 00019 (main.go:10)       CMPQ    CX, $3
        0x0017 00023 (main.go:10)       JLS     68
        0x0019 00025 (main.go:11)       MOVL    "".v+56(SP), AX
        0x001d 00029 (main.go:11)       MOVQ    "".b+32(SP), CX
        0x0022 00034 (main.go:11)       MOVB    AL, 3(CX)
        0x0025 00037 (main.go:12)       MOVL    AX, DX
        0x0027 00039 (main.go:12)       SHRL    $8, AX
        0x002a 00042 (main.go:12)       MOVB    AL, 2(CX)
        0x002d 00045 (main.go:13)       MOVL    DX, AX
        0x002f 00047 (main.go:13)       SHRL    $16, DX
        0x0032 00050 (main.go:13)       MOVB    DL, 1(CX)
        0x0035 00053 (main.go:14)       SHRL    $24, AX
        0x0038 00056 (main.go:14)       MOVB    AL, (CX)
        0x003a 00058 (main.go:15)       MOVQ    16(SP), BP
        0x003f 00063 (main.go:15)       ADDQ    $24, SP
        0x0043 00067 (main.go:15)       RET
        0x0044 00068 (main.go:10)       MOVL    $3, AX
        0x0049 00073 (main.go:10)       CALL    runtime.panicIndex(SB)
        0x004e 00078 (main.go:10)       XCHGL   AX, AX
```
明显看出来，乱序后没有编译器指令优化。

从这里我们就可以知道为什么要先写`_ = b[3]`这样语句判断下下标边界，而后续的四行代码被编译器优化后 只有对 `v` 参数做`BSWAPL`也就没有检查边界的地方。