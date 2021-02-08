package main

import "log"

func main () {
	{
		var data map[string]int
		log.Print(`var data map[string]int`, data)
	}
	func() {
		defer func() {
			log.Print(recover())
		}()
		var data map[string]int
		// 声明了但是没有分配空间，会导致赋值时 panic
		data["age"] = 1
	}()
	{
		// 初始化并赋值
		var data = map[string]int{}
		log.Print(`var data = map[string]int{} `, data) // map[]
		data["name"] = 1
		log.Print(`var data = map[string]int{} `, data) // map[name:1]
	}
	{
		// 通过 make 分配空间
		var data = make(map[string]int)
		log.Print(`var data = make(map[string]int) `, data) // map[]
		data["name"] = 1
		log.Print(`var data = make(map[string]int) `, data) // map[name:1]
	}
}

