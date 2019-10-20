# Go http client

## http client 
http client是一个十分高频使用的组件，特别是近几年基于http协议restful api的盛行，很多语言或者第三库都在设计和打造好用，而且稳定的http client，比如java的apache httpclient，android的okhttp，c语言的libcurl等。


## go 的http client
go标准库实现一个叫`http.DefaultClient`的http client。 本文主要介绍我们比较常用场景下如何使用`http.DefaultClient`。
我们先实现一个go的web服务。

```golang
var addr = ":8999"
func main()  {
		http.HandleFunc("/", func(write http.ResponseWriter, r* http.Request){
			r.ParseForm()
			fmt.Printf("method:%s,a:%s,c:%s",r.Method,r.FormValue("a"),r.FormValue("c"))
		})
		http.HandleFunc("/json", func(write http.ResponseWriter, r* http.Request){
			b,_:= ioutil.ReadAll(r.Body)
			fmt.Println(string(b))
		})
		http.HandleFunc("/file", func(write http.ResponseWriter, r* http.Request){
			file, header, err :=r.FormFile("uploadfile")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			nameParts := strings.Split(header.Filename, ".")
			ext := nameParts[1]
			savedPath := nameParts[0] + "_up."+ext
			f, err := os.OpenFile(savedPath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				panic(err)
			}

		})
		http.ListenAndServe(addr,nil)
}
```

然后运行

```bash
$ go run main.go
```

## Get
go的http client get请求相对简单，我们看下代码
```golang
func get()  {
	resp,err := http.Get("http://localhost"+addr+"/?a=b&c=d")
	if err !=nil {
		log.Fatal(err)
	}
	defer resp.Body.Close() //这边需要把body close
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}
``` 
需要注意的是 `http.DefaultClient` 实现了连接池，`resp.Body.Close()` 代表这次请求已经处理完了，连接会重新放到池子里。

## Post

### x-www-form-urlencoded
```golang
func post()  {
	resp,err := http.Post("http://localhost"+addr, "application/x-www-form-urlencoded",
		strings.NewReader("a=b&c=d"))
	if err !=nil {
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}
```

### form-data
```golang
func postform()  {
	resp,err := http.PostForm("http://localhost"+addr, url.Values{"a": {"b"}, "c": {"d"}})
	if err !=nil {
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}
```

### body

```golang
func postjson()  {
	jsonStr :=[]byte(`{{"a":"b"},{"c":"d"}}`)
	req, err := http.NewRequest("POST", "http://localhost"+addr+"/json", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}
```

### 文件上传

```golang
func fileupload()  {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile("uploadfile", "test.txt") //第一个字段名，第二个是参数名
	if err != nil {
		log.Fatalf("Create form file failed: %s\n", err)
	}
	srcFile, err := os.Open("test.txt")
	if err != nil {
		log.Fatalf("%Open source file failed: s\n", err)
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Fatalf("Write to form file falied: %s\n", err)
	}
	writer.Close()
	resp,err := http.Post("http://localhost"+addr+"/file", writer.FormDataContentType(), buf)
	if err !=nil {
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return
	}
	fmt.Println(string(body))
}
```
## go第三方http client
标准库的http client功能基本上能满足我们的日常开发需求，当然还有第三方package 可能封装得可能更优雅。下面提供三个比较高人气的go的第三方http client实现。
- [beego-httplib](https://github.com/astaxie/beego/tree/develop/httplib)
- [gorequest](https://github.com/parnurzeal/gorequest)
- [resty](https://github.com/go-resty/resty)