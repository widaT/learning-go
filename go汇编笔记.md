# go 汇编笔记


# JLS 指令
```assembly
CMPQ CX，$3
JLS  48  
```

JLS(通JBE一个意思)转移条件：`CMPQ CX，$3 ；JLS 48` 当CX的值小于或等于3的时候转跳到48。