# go 命令

go语言本身自带一个命令行工具，这些命令在项目开发中会反复的使用，我们发点时间先了解下

```
$ go
Go is a tool for managing Go source code.

Usage:
	go <command> [arguments]
The commands are:
	bug         start a bug report
	build       compile packages and dependencies
	clean       remove object files and cached files
	doc         show documentation for package or symbol
	env         print Go environment information
	fix         update packages to use new APIs
	fmt         gofmt (reformat) package sources
	generate    generate Go files by processing source
	get         download and install packages and dependencies
	install     compile and install packages and dependencies
	list        list packages or modules
	mod         module maintenance
	run         compile and run Go program
	test        test packages
	tool        run specified go tool
	version     print Go version
	vet         report likely mistakes in packages
```

本小节我们主要关注 go build ，go install，go get，go test，go mod，go run