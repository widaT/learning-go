# os 包

`os`包提供了与平台无关操作系统函数接口。


## 程序运行环境相关函数

```go
func Hostname() (name string, err error) //返回内核提供的主机名。
func Getpagesize() int //返回底层的系统内存页的尺寸。
func Environ() []string //返回表示环境变量的格式为"key=value"的字符串的切片拷贝。
func Getenv(key string) string //检索并返回名为key的环境变量的值。如果不存在该环境变量会返回空字符串。
func Setenv(key, value string) error //设置名为key的环境变量。如果出错会返回该错误。
func Exit(code int) //让当前程序以给出的状态码code退出。状态码0表示成功，非0表示出错,程序会立刻终止，defer的函数不会被执行。

func Getuid() int //返回调用者的用户ID。
func Geteuid() int //返回调用者的有效用户ID。
func Getgid() int //返回调用者的组ID。
func Getegid() int //返回调用者的有效组ID。
func Getgroups() ([]int, error) //返回调用者所属的所有用户组的组ID。
func Getpid() int //返回调用者所在进程的进程ID。
func Getppid() int //返回调用者所在进程的父进程的进程ID。
func Getwd() (dir string, err error) //返回程序当前的工作目录
``` 

## 目录和文件操作

### 创建文件夹

```go
func Mkdir(name string, perm FileMode) error //使用指定的权限和名称创建一个目录
func MkdirAll(path string, perm FileMode) error //MkdirAll使用指定的权限和名称创建一个目录，包括任何必要的上级目录.
```

```go
os.Mkdir("aa", 0755)
fmt.Println(os.MkdirAll("aa/bb/", 0755))
```

### 文件读写

```go
func OpenFile(name string, flag int, perm FileMode) (file *File, err error) //通用的打开文件方法
func Remove(name string) error //删除文件或者目录
```

golang中用`*File`代表文件句柄。`*File`定义了很多操作它的方法。

```go
func (f *File) Chdir() error
func (f *File) Chmod(mode FileMode) error
func (f *File) Chown(uid, gid int) error
func (f *File) Close() error
func (f *File) Fd() uintptr
func (f *File) Name() string
func (f *File) Read(b []byte) (n int, err error)
func (f *File) ReadAt(b []byte, off int64) (n int, err error)
func (f *File) ReadFrom(r io.Reader) (n int64, err error)
func (f *File) Readdir(n int) ([]FileInfo, error)
func (f *File) Readdirnames(n int) (names []string, err error)
func (f *File) Seek(offset int64, whence int) (ret int64, err error)
func (f *File) SetDeadline(t time.Time) error
func (f *File) SetReadDeadline(t time.Time) error
func (f *File) SetWriteDeadline(t time.Time) error
func (f *File) Stat() (FileInfo, error)
func (f *File) Sync() error
func (f *File) SyscallConn() (syscall.RawConn, error)
func (f *File) Truncate(size int64) error
func (f *File) Write(b []byte) (n int, err error)
func (f *File) WriteAt(b []byte, off int64) (n int, err error)
func (f *File) WriteString(s string) (n int, err error)
```

```go
f, err := os.OpenFile("a.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
if err != nil {
	log.Fatal(err)
}
f.Write([]byte("aaaa"))
f.WriteString("bbbb")
f.Close() //打开成功的文件句柄 不用的时候一定记得关闭
ff, err := os.OpenFile("a.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
b := make([]byte, 1024)
n, err := ff.Read(b)
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(b[:n])) //aaaabbbb
os.Remove("a.txt")  //删除文件
f.Close()
```

# 参考文档

[go pkg os](https://golang.org/pkg/os/)