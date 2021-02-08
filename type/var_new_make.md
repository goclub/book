# go 中 var new make 的区别

`var` 用于声明变量。

`new` 分配内存空间，`func new(Type) *Type` 接收 一个类型，返回这个类型的指针，并将指针指向这个类型的零值（zero value）。

`make` 分配内存空间并根据参数初始化 

> 本文主要通过代码示例和原因来解释 var new make 之间的区别。

## var new


通过代码记忆最为合适

[var_new](./var_new/main.go)
```.go
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
		// 或者使用 new （分配空间，并将指针指向零值的空间）
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

```

## var make slice array 

[var_make_slice_array](./var_make_slice_array/main.go)
```.go
package main

import "log"

func main() {
	{
		var names = []string{"a","b","c"}
		log.Print(`var names = []string{"a","b","c"} `, names, len(names)) // [a b c] 3
	}
	log.Print("-----------")
	{
		var names = make([]string, 0)
		// 等同 names := make([]string, 0)
		names = append(names, "a", "b", "c")
		log.Print(`var names = make([]string, 0) `, names, len(names)) // [a b c] 3
	}
	{
		// 可以通过 make 分配空间的同时初始化长度，初始化元素值为类型的空值(zero value)
		var names = make([]string, 2)
		// 等同 names := make([]string, 0)
		names = append(names, "a", "b", "c")
		log.Print(`var names = make([]string, 2) `, names, len(names)) // [  a b c] 5
	}
	log.Print("-----------")
	{
		// [2]string 表示是长度为2由string组成的 array
		var names = [2]string{}
		// 等同: var names = [2]string{}
		log.Print("var names = [2]string{} ", names, len(names)) // [ ] 2
	}
	log.Print("-----------")
	// make([]string, 2, 2)  是 make 最常用的使用场景
	{
		arrayLen := 2
		arrayCap := 2
		var names = make([]string, arrayLen, arrayCap)
		log.Print(`make([]string, 2, 2) 初始化设置数组长度2，容量为2`, names, len(names), cap(names))
		names[0] = "a"
		log.Print(`names[0] = "a" `, names, len(names), cap(names))
	}
	log.Print("-----------")
	{
		arrayLen := 0
		arrayCap := 2
		var names = make([]string, arrayLen, arrayCap)
		log.Print(`make([]string, 0, 2) 初始化设置数组长度0，容量为2`, names, len(names), cap(names))
		names = append(names, "a")
		log.Print(`append`, names, len(names), cap(names))
		names = append(names, "b")
		log.Print(`append`, names, len(names), cap(names))
		names = append(names, "c")
		log.Print(`append（如果数组容量不够，运行 append 会自动扩充容量）`, names, len(names), cap(names))
	}
	log.Print("-----------")
	func() {
		defer func() {
			log.Print(recover())
		}()
		// 大部分场景下 names := make([]string, len, cap) len 和 cap 设置的不同时没有意义的
		names := make([]string, 0, 2)
		// 因为 names[0] = "a" 会panic
		names[0] = "a"
	}()
}

```

## var make map

[var_make_map](./var_make_map/main.go)
```.go
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


```

## make chan

[make_chan](./make_chan/main.go)
```.go
package main

import "log"

func main() {
	{
		var nameCh chan string
		nameCh = make(chan string) // 注释这一行会因为 nameCh 没有分配内存空间导致死锁
		log.Print("nameCh ", nameCh) // 内存地址
		go func() {
			nameCh <- "nimoc"
		}()
		name := <-nameCh
		log.Print(name)
	}
	{
		{
			// 代码可以更简洁一点
			nameCh := make(chan string)
			log.Print("nameCh ", nameCh) // 内存地址
			go func() {
				nameCh <- "nimoc"
			}()
			name := <-nameCh
			log.Print(name)
		}
	}
}

```