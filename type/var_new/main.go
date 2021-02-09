package main

import "log"

func main () {
	{
		var i int
		i++
		log.Print("var i int ",i)
	}
	{
		func() {
			defer func() {
				log.Print(recover())
			}()
			var i *int
			// panic: runtime error: invalid memory address or nil pointer dereference
			// 通过 var 声明的 *int 是空指针 递增会报错
			*i++
			log.Print("var i *int ",*i)
		}()
	}
	{
		var v int
		// 通过 = & 可以避免
		var i *int = &v
		*i++
		log.Print("var i *int = &v ",*i)
	}
	{
		// 或者使用 new （分配空间，并将指针指向零值）
		var i *int = new(int)
		*i++
		log.Print("var i *int = new(int)",*i)
	}
	{
		// 可以进一步简写
		i := new(int)
		*i++
		log.Print("i := new(int)",*i)
	}
}
