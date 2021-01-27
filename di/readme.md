# 由浅入深的讲清楚 Go 依赖注入

> DI (dependency injection)

网络上有很多依赖注入的教程和文章，大部分都解释了什么是DI。能做到浅显易懂的却很少。在 Go 使用 依赖注入的文章就更少了。

本文尝试以实践的角度描述为什么需要依赖注入，在各种场景中使用依赖注入的前后的变化。

## 解决循环依赖

### 编译时循环依赖

使用 Go 永远绕不开的问题是循环依赖。


当存在用户模块和消息模块时：

```go
package cd_user

import (
	cd_message "github.com/goclub/book/di/cyclic_dependency/message"
)

func UserName(userID string) string {
	return "nimoc"
}

func MyMessageList() []string {
	return cd_message.MessageListByUserID("a")
}
```

```go
package cd_message

import (
	cd_user "github.com/goclub/book/di/cyclic_dependency/user"
	"log"
)

func SendMessage(userID string) {
	userName := cd_user.UserName(userID)
	log.Print("SendMessage: Welcome " + userName + "!")
}

func MessageListByUserID(userID string) []string {
	return []string{"Friend request.", "Welcome to join!"}
}
```

可以看到 message 和 user 互相调用了对方的方法

在 main 中使用它们

```go
package main

import (
	cd_message "github.com/goclub/book/di/import cycle/message"
	cd_user "github.com/goclub/book/di/import cycle/user"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		cd_message.SendMessage("a")
		_, err := writer.Write([]byte(`<a href="/message" >message list</a>`)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	http.HandleFunc("/message", func(writer http.ResponseWriter, request *http.Request) {
		messageList := cd_user.MyMessageList()
		data := "message list" + strings.Join(messageList, "\n")
		_, err := writer.Write([]byte(data)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	err := http.ListenAndServe(":1219", nil) ; if err != nil {
		panic(err)
	}
}

```

通过 `go run main.go` 或者 `go build main.go` 会出现错误

```go
package github.com/goclub/book/di/cyclic_dependency
	imports github.com/goclub/book/di/cyclic_dependency/message
	imports github.com/goclub/book/di/cyclic_dependency/user
	imports github.com/goclub/book/di/cyclic_dependency/message: import cycle not allowed
```

[源码](https://github.com/goclub/book/blob/main/di/cyclic_dependency)

> 这个例子中在程序设计上 user 和 message 的耦合不是很恰当，但为了简单的表达出 go 的依赖注入错误，请暂时无视程序设计上的耦合。
> 因为在实际工作中有些模块必须互相依赖

如果开发的是第三方包，则很容易解决依赖问题，只需要将 user 和 message 合并到一个模块即可。
但是在业务代码中多个模块合并成一个模块会导致代码全聚合在一起难以维护。

此时可以通过依赖注入的方式解决问题，核心是使用 Go 的 `interface`

### 使用注入函数解决循环依赖

go 要求开发人员避免循环依赖，我们应当尽量遵循这个要求。但是在业务代码中，循环依赖是无法避免的。

使用注入函数让 user 与 message 直接不直接依赖，由 main 包将 user 和 message 作为参数传递给对方。（注入）

 
```go
package func_user


func UserName(userID string) string {
	return "nimoc"
}

func MyMessageList(MessageListByUserID func(userID string) []string ) []string {
	return MessageListByUserID("a")
}
```

```go
package func_message

import (
	"log"
)

func SendMessage(userID string, UserName func(userID string) string ) {
	userName := UserName(userID)
	log.Print("SendMessage: Welcome " + userName + "!")
}

func MessageListByUserID(userID string) []string {
	return []string{"Friend request.", "Welcome to join!"}
}
```
```go
package main

import (
	func_message "github.com/goclub/book/di/func_user_message/message"
	func_user "github.com/goclub/book/di/func_user_message/user"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		dep := func_user.UserName
		func_message.SendMessage("a", dep)
		_, err := writer.Write([]byte(`<a href="/message" >message list</a>`)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	http.HandleFunc("/message", func(writer http.ResponseWriter, request *http.Request) {
		dep := func_message.MessageListByUserID
		messageList := func_user.MyMessageList(dep)
		data := "message list" + strings.Join(messageList, "\n")
		_, err := writer.Write([]byte(data)) ; if err != nil {
			writer.WriteHeader(500)
			log.Print(err)
		}
	})
	log.Print("http://127.0.0.1:1219")
	err := http.ListenAndServe(":1219", nil) ; if err != nil {
		panic(err)
	}
}
```

[源码](https://github.com/goclub/book/blob/main/di/func_user_message)

此时编译不会出现循环依赖的错误，因为 user.go 和 message.go 文件中均没有依赖对方。
运行时候依赖是由于 main 通过参数传递给 user 和 message 的，此时就避免解决了循环依赖。

## 使用 interface 解决循环依赖


 