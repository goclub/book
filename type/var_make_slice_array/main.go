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
