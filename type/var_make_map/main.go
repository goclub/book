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
		// 只声明了变量
		var data map[string]int
		// 声明了但是没有分配空间，会导致赋值时 panic
		data["age"] = 1
	}()
	{
		// 声明变量并初始化赋值
		var data = map[string]int{}
		log.Print(`var data = map[string]int{} `, data) // map[]
		data["name"] = 1
		log.Print(`var data = map[string]int{} `, data) // map[name:1]
	}
	{
		// 通过 make 分配空间 也可以避免panic
		var data = make(map[string]int)
		log.Print(`var data = make(map[string]int) `, data) // map[]
		data["name"] = 1
		log.Print(`var data = make(map[string]int) `, data) // map[name:1]
	}
}

