# http请求参数校验

在web服务开发的过程中，我们经常需要对用户传递的参数进行验证。验证的代码很容易写得冗长，而且比较丑陋。本文我们介绍一个第三放库[validator](https://github.com/go-playground/validator)专门解决这个问题。

## 验证单一变量
```golang
import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)
func validateVariable() {
	myEmail := "someone.gmail.com"
	errs := validate.Var(myEmail, "required,email")
	if errs != nil {
		fmt.Println(errs)
		return
    }
    //这边验证通过 写逻辑代码
}
var validate *validator.Validate
func main() {
	validate = validator.New()
	validateVariable()
}
```
运行一下
```bash
$ go run main.go
Key: '' Error:Field validation for '' failed on the 'email' tag
```

## 验证结构体
```golang
package main

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	Name      	   string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=150"`  //大于等于0 小于等于150
	Email          string     `validate:"required,email"`
}

var validate *validator.Validate
func main() {
	validate = validator.New()
	validateStruct()
}

func validateStruct() {
	user := &User{
		Name:      		"wida",
		Age:            165,
		Email:          "someone.gmail.com",
	}

	err := validate.Struct(user)
	fmt.Println(err)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}
		for _, err := range err.(validator.ValidationErrors) {  //变量所有参数错误
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField()) 
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		return
	}
}
```

运行一下
```bash
$ go run main.go
Key: 'User.Age' Error:Field validation for 'Age' failed on the 'lte' tag
Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag
User.Age
Age
User.Age
Age
lte
lte
uint8
uint8
165
150

User.Email
Email
User.Email
Email
email
email
string
string
someone.gmail.com
```
运行的结果显示，`age`和`email`验证通不过。