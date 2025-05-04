# Go汇编基础

## Plan9汇编基本概念

Go汇编是基于Plan9汇编改编而来，它是一种独特的汇编方言。Go汇编的设计目标是提供一种跨平台的汇编语言抽象，使得同一段汇编代码可以在不同的硬件架构上运行 <mcreference link="https://github.com/yangyuqian/technical-articles/blob/master/asm/golang-plan9-assembly-cn.md" index="1">1</mcreference>。

### 虚拟寄存器

Go汇编中引入了四个重要的虚拟寄存器 <mcreference link="https://hopehook.com/post/golang_assembly/" index="2">2</mcreference>：

1. FP (Frame pointer)：用于访问函数的参数和局部变量
2. PC (Program counter)：程序计数器，用于跳转和分支
3. SB (Static base pointer)：用于引用全局符号
4. SP (Stack pointer)：指向栈顶

### 操作数和指令

在Go汇编中，操作数的方向与Intel汇编相反，是从左到右的。主要的指令类型包括 <mcreference link="https://hopehook.com/post/golang_assembly/" index="2">2</mcreference>：

- 数据移动指令：
  - MOVB：移动1字节
  - MOVW：移动2字节
  - MOVD：移动4字节
  - MOVQ：移动8字节

- 运算指令：
  - ADD：加法运算
  - SUB：减法运算
  - IMUL：乘法运算

- 栈操作：
  - 不使用PUSH/POP，而是通过SUB/ADD指令操作SP实现

## 内存布局

### 数据定义

Go汇编中使用DATA和GLOBL指令来定义变量 <mcreference link="https://github.com/yangyuqian/technical-articles/blob/master/asm/golang-plan9-assembly-cn.md" index="1">1</mcreference>：

```asm
// 定义全局变量
DATA symbol+offset(SB)/width, value
GLOBL symbol(SB), NOPTR, $size
```

### 函数声明

函数声明的基本格式 <mcreference link="https://xargin.com/go-and-plan9-asm/" index="3">3</mcreference>：

```asm
TEXT pkgname·funcname(SB), NOSPLIT, $framesize-argsize
```

- pkgname：包名
- funcname：函数名
- NOSPLIT：不进行栈分裂检查
- framesize：栈帧大小
- argsize：参数和返回值的大小

## 寄存器使用

除了虚拟寄存器外，Go汇编还可以使用目标平台的物理寄存器。在AMD64架构下，常用的通用寄存器包括：

- RAX, RBX, RCX, RDX：通用数据寄存器
- RSI, RDI：源索引和目标索引寄存器
- R8-R15：额外的通用寄存器

## 注意事项

1. Go汇编不能独立使用，必须与Go源码文件一起编译
2. 在编写汇编代码时，需要特别注意栈的管理和内存对齐
3. 使用Go汇编时，应该重点关注函数调用中的参数传递和栈帧布局
4. 编写汇编代码时应遵循Go的调用约定和ABI规范

## 调试技巧

可以使用以下命令查看Go代码对应的汇编代码：

```bash
go build -gcflags "-S" main.go    # 查看汇编代码
go tool compile -S main.go      # 查看优化后的汇编代码
```