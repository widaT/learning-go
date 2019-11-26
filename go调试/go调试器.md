# go 调试器
go 没有官方的调式器，有个第三方调试工具[delve](https://github.com/go-delve/delve)在社区非常受欢迎。

## 安装delve
安装到GOBIN目录
```bash
$ go get -u github.com/go-delve/delve/cmd/dlv
```
将GOBIN目录添加到`PATH`环境变量中。

```bash
$ dlv  version
Delve Debugger
Version: 1.3.0
Build: $Id: 2f59bfc686d60989dcef9de40b480d0a34aa2fa5
```
ok,安装成功。

## delve使用
我先写一下要被调试的程序。
```golang
package main
import (
	"fmt"
	"time"
)
func test(a int) int {
	d :=3
	a = d+3
	return a
}
func main()  {
	a ,b:=1,3
	for i:=0;i<5;i++ {
		a++
	}
	test(a)
	go func() {
		a +=1
		b -=1
		time.Sleep(10e9)
		b +=1
	}()
	fmt.Println(a,b)
	time.Sleep(100e9)
}
```
## dlv debug 调试源码
用到`main package`目录下运行
```bash
$ dlv debug
Type 'help' for list of commands.
(dlv) h
The following commands are available:
    args ------------------------ Print function arguments.
    ... //这边省略
Type help followed by a command for full documentation.
```
我们看了下使用帮助
```bash
    args ------------------------ Print function arguments.
    break (alias: b) ------------ Sets a breakpoint.
    breakpoints (alias: bp) ----- Print out info for active breakpoints.
    call ------------------------ Resumes process, injecting a function call (EXPERIMENTAL!!!)
    clear ----------------------- Deletes breakpoint.
    clearall -------------------- Deletes multiple breakpoints.
    condition (alias: cond) ----- Set breakpoint condition.
    config ---------------------- Changes configuration parameters.
    continue (alias: c) --------- Run until breakpoint or program termination.
    deferred -------------------- Executes command in the context of a deferred call.
    disassemble (alias: disass) - Disassembler.
    down ------------------------ Move the current frame down.
    edit (alias: ed) ------------ Open where you are in $DELVE_EDITOR or $EDITOR
    exit (alias: quit | q) ------ Exit the debugger.
    frame ----------------------- Set the current frame, or execute command on a different frame.
    funcs ----------------------- Print list of functions.
    goroutine (alias: gr) ------- Shows or changes current goroutine
    goroutines (alias: grs) ----- List program goroutines.
    help (alias: h) ------------- Prints the help message.
    libraries ------------------- List loaded dynamic libraries
    list (alias: ls | l) -------- Show source code.
    locals ---------------------- Print local variables.
    next (alias: n) ------------- Step over to next source line.
    on -------------------------- Executes a command when a breakpoint is hit.
    print (alias: p) ------------ Evaluate an expression.
    regs ------------------------ Print contents of CPU registers.
    restart (alias: r) ---------- Restart process.
    set ------------------------- Changes the value of a variable.
    source ---------------------- Executes a file containing a list of delve commands
    sources --------------------- Print list of source files.
    stack (alias: bt) ----------- Print stack trace.
    step (alias: s) ------------- Single step through program.
    step-instruction (alias: si)  Single step a single cpu instruction.
    stepout (alias: so) --------- Step out of the current function.
    thread (alias: tr) ---------- Switch to the specified thread.
    threads --------------------- Print out info for every traced thread.
    trace (alias: t) ------------ Set tracepoint.
    types ----------------------- Print list of types
    up -------------------------- Move the current frame up.
    vars ------------------------ Print package variables.
    whatis ---------------------- Prints type of an expression.
```
这个使用帮助的文档英文都不需要翻译，相信大家都能看懂。
你不需要死级记得上面的命令，进入dlv调试模式后，敲你命令搜字母 `tab` 一下就可以自动补全。实在忘记了，你再 `help` 一下就会显示帮助文档。


### 打断点

delve 打断点的命令是`break` 简写`b`，打断点有两种方式
- 我们知道某个包的某个函数比如  `main.main`  我们使用 `break main.main` 就可以设置断点
- 我们知道某个源码文件第几行，我们可以用`break 文件名:行号`的方式设置断点,例如 `break main.go:7`  
使用`breakpoints` 可以查看所用设置的断点

```bash
$ dlv debug
Type 'help' for list of commands.
(dlv) break main.main
Breakpoint 1 set at 0x4aa148 for main.main() ./main.go:5
(dlv) b main.go:7
Breakpoint 2 set at 0x4aa1ab for main.main() ./main.go:7
(dlv) breakpoints
Breakpoint runtime-fatal-throw at 0x42ddc0 for runtime.fatalthrow() /home/wida/go/src/runtime/panic.go:820 (0)
Breakpoint unrecovered-panic at 0x42de30 for runtime.fatalpanic() /home/wida/go/src/runtime/panic.go:847 (0)
        print runtime.curg._panic.arg
Breakpoint 1 at 0x4aa148 for main.main() ./main.go:5 (0)
Breakpoint 2 at 0x4aa1ab for main.main() ./main.go:7 (0)
(dlv)
```

设置好断点后我们可以使用 `continue` 让程序运行到断点。多个断点的话，再按`continue`就能到下一个断点。
 ```bash
 (dlv) c
> main.main() ./main.go:12 (hits goroutine(1):1 total:1) (PC: 0x4aa3f8)
     7: func test(a int) int {
     8:         d :=3
     9:         a = d+3
    10:         return a
    11: }
=>  12: func main()  {
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
    17:         test(a)
(dlv) c
> main.main() ./main.go:15 (hits goroutine(1):1 total:1) (PC: 0x4aa476)
    10:         return a
    11: }
    12: func main()  {
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
=>  15:                 a++
    16:         }
    17:         test(a)
    18:         go func() {
    19:                 a +=1
    20:                 b -=1
(dlv) 
```
### 删除断点

`clear` 和 `clearall` 删除断点，`clear`命令后面有个参数 `breakId`
```bash
(dlv) breakpoints
Breakpoint runtime-fatal-throw at 0x42ddc0 for runtime.fatalthrow() /home/wida/go/src/runtime/panic.go:820 (0)
Breakpoint unrecovered-panic at 0x42de30 for runtime.fatalpanic() /home/wida/go/src/runtime/panic.go:847 (0)
        print runtime.curg._panic.arg
Breakpoint 1 at 0x4aa3f8 for main.main() ./main.go:12 (1)
Breakpoint 2 at 0x4aa476 for main.main() ./main.go:15 (1)
        cond i == 3
(dlv) clear 1
(dlv) breakpoints
Breakpoint runtime-fatal-throw at 0x42ddc0 for runtime.fatalthrow() /home/wida/go/src/runtime/panic.go:820 (0)
Breakpoint unrecovered-panic at 0x42de30 for runtime.fatalpanic() /home/wida/go/src/runtime/panic.go:847 (0)
        print runtime.curg._panic.arg
Breakpoint 2 at 0x4aa476 for main.main() ./main.go:15 (1)
        cond i == 3
(dlv) clearall
Breakpoint 2 cleared at 0x4aa476 for main.main() ./main.go:15
(dlv) breakpoints
Breakpoint runtime-fatal-throw at 0x42ddc0 for runtime.fatalthrow() /home/wida/go/src/runtime/panic.go:820 (0)
Breakpoint unrecovered-panic at 0x42de30 for runtime.fatalpanic() /home/wida/go/src/runtime/panic.go:847 (0)
        print runtime.curg._panic.arg

```

#### 断点上附加条件

`condition` 命令可以已有的断点上加上命令，有连个参数 `condition breakId condition` 第一个参数是已有断点的id（用breakpoints 可以知道），第二个是附加条件。

```bash
(dlv) breakpoints
Breakpoint runtime-fatal-throw at 0x42ddc0 for runtime.fatalthrow() /home/wida/go/src/runtime/panic.go:820 (0)
Breakpoint unrecovered-panic at 0x42de30 for runtime.fatalpanic() /home/wida/go/src/runtime/panic.go:847 (0)
        print runtime.curg._panic.arg
Breakpoint 1 at 0x4aa3f8 for main.main() ./main.go:12 (0)
Breakpoint 2 at 0x4aa476 for main.main() ./main.go:15 (0)
(dlv) condition 2 i==3
(dlv) c
> main.main() ./main.go:12 (hits goroutine(1):1 total:1) (PC: 0x4aa3f8)
     7: func test(a int) int {
     8:         d :=3
     9:         a = d+3
    10:         return a
    11: }
=>  12: func main()  {
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
    17:         test(a)
(dlv) c
> main.main() ./main.go:15 (hits goroutine(1):1 total:1) (PC: 0x4aa476)
    10:         return a
    11: }
    12: func main()  {
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
=>  15:                 a++
    16:         }
    17:         test(a)
    18:         go func() {
    19:                 a +=1
    20:                 b -=1
(dlv) locals
b = 3    
a = 4
i = 3  //condition 命令已经生效
```

### 打印变量

- `locals` 打印局部变量
- `vars` 打印所有包变量 `vars 包名`打印包变量
- `print` 打印变量值

```bash
(dlv) c
> main.main() ./main.go:17 (hits goroutine(1):1 total:1) (PC: 0x4aa4a0)
    12: func main()  {
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
=>  17:         test(a)
    18:         go func() {
    19:                 a +=1
    20:                 b -=1
    21:                 time.Sleep(10e9)
    22:                 b +=1
(dlv) locals
b = 3
a = 6
(dlv) print a
6
(dlv) print &a
(*int)(0xc0000ac010)
(dlv) print b
3
(dlv)
```

### 单步执行

- `next`单步执行（遇到子函数不进入，只执行本文件下一行代码）
- `step`单步进入执行（遇到子函数会进入子函数）。
- `stepout` 退出（执行完）当前函数到该函数被调用的下一行代码。

这边需要注意 `next`和`step`的区别
```bash
(dlv) b main.go:17
Breakpoint 1 set at 0x4aa4a0 for main.main() ./main.go:17
(dlv) c
> main.main() ./main.go:17 (hits goroutine(1):1 total:1) (PC: 0x4aa4a0)
    12: func main()  {
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
=>  17:         test(a)
    18:         go func() {
    19:                 a +=1
    20:                 b -=1
    21:                 time.Sleep(10e9)
    22:                 b +=1
(dlv) n
> main.main() ./main.go:18 (PC: 0x4aa4b9)
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
    17:         test(a)
=>  18:         go func() {
    19:                 a +=1
    20:                 b -=1
    21:                 time.Sleep(10e9)
    22:                 b +=1
    23:         }()
(dlv) r
Process restarted with PID 706
(dlv) c
> main.main() ./main.go:17 (hits goroutine(1):1 total:1) (PC: 0x4aa4a0)
    12: func main()  {
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
=>  17:         test(a)
    18:         go func() {
    19:                 a +=1
    20:                 b -=1
    21:                 time.Sleep(10e9)
    22:                 b +=1
(dlv) s
> main.test() ./main.go:7 (PC: 0x4aa3a0)
     2: 
     3: import (
     4:         "fmt"
     5:         "time"
     6: )
=>   7: func test(a int) int {
     8:         d :=3
     9:         a = d+3
    10:         return a
    11: }
    12: func main()  {
(dlv) n
> main.test() ./main.go:8 (PC: 0x4aa3b7)
     3: import (
     4:         "fmt"
     5:         "time"
     6: )
     7: func test(a int) int {
=>   8:         d :=3
     9:         a = d+3
    10:         return a
    11: }
    12: func main()  {
    13:         a ,b:=1,3
(dlv) so
> main.main() ./main.go:18 (PC: 0x4aa4b9)
Values returned:
        ~r1: 6

    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
    17:         test(a)
=>  18:         go func() {
    19:                 a +=1
    20:                 b -=1
    21:                 time.Sleep(10e9)
    22:                 b +=1
    23:         }()
(dlv)
```

### goroutine 切换

我们在调试的时候不能同事看到所有goroutine的情况，默认我们追踪的代码是在`main goroutine`中，我们有时候想看看非`main goroutine`的情况，就需要goroutine 切换的命令。
- `goroutines` 打印所有的goroutine
- `goroutine` 切换指定id的goroutine执行

```bash
Breakpoint 1 cleared at 0x4aa4a0 for main.main() ./main.go:17
(dlv) b main.go:18
(dlv) c
> main.main() ./main.go:18 (hits goroutine(1):1 total:1) (PC: 0x4aa4b9)
    13:         a ,b:=1,3
    14:         for i:=0;i<5;i++ {
    15:                 a++
    16:         }
    17:         test(a)
=>  18:         go func() {
    19:                 a +=1
    20:                 b -=1
    21:                 time.Sleep(10e9)
    22:                 b +=1
    23:         }()
(dlv) n
> main.main() ./main.go:24 (PC: 0x4aa4f7)
    19:                 a +=1
    20:                 b -=1
    21:                 time.Sleep(10e9)
    22:                 b +=1
    23:         }()
=>  24:         fmt.Println(a,b)
    25:         time.Sleep(100e9)
    26: }
(dlv) goroutines
* Goroutine 1 - User: ./main.go:24 main.main (0x4aa4f7) (thread 1903)
  Goroutine 2 - User: /home/wida/go/src/runtime/proc.go:305 runtime.gopark (0x42faab)
  Goroutine 3 - User: /home/wida/go/src/runtime/proc.go:305 runtime.gopark (0x42faab)
  Goroutine 4 - User: /home/wida/go/src/runtime/proc.go:305 runtime.gopark (0x42faab)
  Goroutine 5 - User: /home/wida/go/src/runtime/lock_futex.go:228 runtime.notetsleepg (0x40ab14)
  Goroutine 17 - User: /home/wida/go/src/runtime/proc.go:305 runtime.gopark (0x42faab)
  Goroutine 18 - User: /home/wida/go/src/runtime/time.go:105 time.Sleep (0x44aa51)
[7 goroutines]
(dlv) goroutine 18
Switched from 1 to 18 (thread 1903)
```

### 重启和退出

- `restart` 重新执行程序，设置的断点还会在。
- `exit` 退出调试程序。

```bash
(dlv) b main.main
Breakpoint 1 set at 0x4aa3f8 for main.main() ./main.go:12
(dlv) breakpoints
Breakpoint runtime-fatal-throw at 0x42ddc0 for runtime.fatalthrow() /home/wida/go/src/runtime/panic.go:820 (0)
Breakpoint unrecovered-panic at 0x42de30 for runtime.fatalpanic() /home/wida/go/src/runtime/panic.go:847 (0)
        print runtime.curg._panic.arg
Breakpoint 1 at 0x4aa3f8 for main.main() ./main.go:12 (0)
(dlv) restart
Process restarted with PID 2622
(dlv) breakpoints
Breakpoint runtime-fatal-throw at 0x42ddc0 for runtime.fatalthrow() /home/wida/go/src/runtime/panic.go:820 (0)
Breakpoint unrecovered-panic at 0x42de30 for runtime.fatalpanic() /home/wida/go/src/runtime/panic.go:847 (0)
        print runtime.curg._panic.arg
Breakpoint 1 at 0x4aa3f8 for main.main() ./main.go:12 (0)
(dlv) exit
```

### 调试汇编代码

- `regs [-a]`查看cpu通用寄存器值，加上`-a`的话查看所有的寄存器的值。
- `disassemble` 查看汇编代码（inter 风格的汇编）

我们在原先代码目录添加一个汇编代码文件 test.s
```assembly
#include "textflag.h"

TEXT ·add(SB), NOSPLIT, $0-8
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    ADDQ AX, BX
    MOVQ BX, ret+16(FP)
    RET
```
修改下main.go的代码
```golang
package main

import (
	"fmt"
	"time"
)

func add(a ,b int) int
func test(a int) int {
	d :=3
	a = d+3
	return a
}
func main()  {
	a ,b:=1,3
	for i:=0;i<5;i++ {
		a++
	}
	test(a)
	add(a,b)
	go func() {
		a +=1
		b -=1
		time.Sleep(10e9)
		b +=1
	}()
	fmt.Println(a,b)
	time.Sleep(100e9)
}
```

```bash
(dlv) c
> main.main() ./main.go:20 (hits goroutine(1):1 total:1) (PC: 0x4aa4bf)
    15:         a ,b:=1,3
    16:         for i:=0;i<5;i++ {
    17:                 a++
    18:         }
    19:         test(a)
=>  20:         add(a,b)
    21:         go func() {
    22:                 a +=1
    23:                 b -=1
    24:                 time.Sleep(10e9)
    25:                 b +=1
(dlv) s
> main.add() ./test.s:4 (PC: 0x4aa6d0)
     1: #include "textflag.h"
     2: 
     3: TEXT ·add(SB), NOSPLIT, $0-8
=>   4:     MOVQ a+0(FP), AX
     5:     MOVQ b+8(FP), BX
     6:     ADDQ AX, BX
     7:     MOVQ BX, ret+16(FP)
     8:     RET
(dlv) regs
     Rip = 0x00000000004aa6d0
     Rsp = 0x000000c000075e78
     Rax = 0x0000000000000003
...
(dlv) n
> main.add() ./test.s:5 (PC: 0x4aa6d5)
     1: #include "textflag.h"
     2: 
     3: TEXT ·add(SB), NOSPLIT, $0-8
     4:     MOVQ a+0(FP), AX
=>   5:     MOVQ b+8(FP), BX
     6:     ADDQ AX, BX
     7:     MOVQ BX, ret+16(FP)
     8:     RET
(dlv) regs
     Rip = 0x00000000004aa6d5
     Rsp = 0x000000c000075e78
     Rax = 0x0000000000000006        //rax 寄存器的值已经发送变化
   ...
(dlv) regs -a
       Rip = 0x00000000004aa6d5
       Rsp = 0x000000c000075e78
       Rax = 0x0000000000000006
      ...
      XMM0 = 0x00000000000000000000000000000000 v2_int={ 0000000000000000 0000000000000000 }    v4_int={ 00000000 00000000 00000000 00000000 }  v8_int={ 0000 0000 0000 0000 0000 0000 0000 0000 }       v16_int={ 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 }     v2_float={ 0 0 }        v4_float={ 0 0 0 0 }
      ...
(dlv) 
```
```bash
(dlv) c
> main.main() ./main.go:20 (hits goroutine(1):1 total:1) (PC: 0x4aa4bf)
    15:         a ,b:=1,3
    16:         for i:=0;i<5;i++ {
    17:                 a++
    18:         }
    19:         test(a)
=>  20:         add(a,b)
    21:         go func() {
    22:                 a +=1
    23:                 b -=1
    24:                 time.Sleep(10e9)
    25:                 b +=1
(dlv) disassemble
TEXT main.main(SB) /home/wida/gocode/debugger-demo/main.go
        main.go:14      0x4aa3e0        64488b0c25f8ffffff              mov rcx, qword ptr fs:[0xfffffff8]
        main.go:14      0x4aa3e9        488d4424a8                      lea rax, ptr [rsp-0x58]
        main.go:14      0x4aa3ee        483b4110                        cmp rax, qword ptr [rcx+0x10]
...
```

## dlv exec 调试可执行文件

我们先编译我们的项目
```bash
$ go build 
$ ll
drwxr-xr-x  3 wida wida    4096 10月 22 11:01 .
drwxr-xr-x 44 wida wida    4096 10月 21 18:25 ..
-rwxr-xr-x  1 wida wida 2013047 10月 22 11:01 debugger-demo
-rw-r--r--  1 wida wida      30 10月 21 18:38 go.mod
-rw-r--r--  1 wida wida     288 10月 22 10:40 main.go
-rw-r--r--  1 wida wida     143 10月 22 10:41 test.s
$ dlv exec ./debugger-demo
(dlv) b main.main
Breakpoint 1 set at 0x48cef3 for main.main() ./main.go:14
(dlv) c
> main.main() ./main.go:14 (hits goroutine(1):1 total:1) (PC: 0x48cef3)
Warning: debugging optimized function
     9: func test(a int) int {
    10:         d :=3
    11:         a = d+3
    12:         return a
    13: }
=>  14: func main()  {
    15:         a ,b:=1,3
    16:         for i:=0;i<5;i++ {
    17:                 a++
    18:         }
    19:         test(a)
(dlv) q
```
后面的调试方式和`dlv exec`类似。


# 参考文档

- [delve Documentation](https://github.com/go-delve/delve/tree/master/Documentation/cli)