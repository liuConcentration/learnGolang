package main

import "fmt"

//函数是编程 VS函数指针
/**
函数是一等公民: 参数,变量,返回值都可以是函数
高阶函数
函数->闭包
 */

//通过闭包来做累加,该函数返回一个闭包
func adder() func(int) int {
	//闭包特征之一,,自由变量
	sum := 0
	//闭包特征之一,这里的这个v是一个参数,也是一个局部变量
	return func(v int) int {
		//闭包特征之一,sum不是在函数体里面定义的,他是这个函数所处的这个环境,sum是外面的
		sum += v
		return sum
	}
}
func main() {

	a := adder()
	for i := 0; i < 10; i++ {
		fmt.Printf("0 + 1 + ...... + %d = %d\n", i, a(i))
	}
}
