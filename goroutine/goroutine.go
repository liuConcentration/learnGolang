package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	//这里循环10个和循环100个1000个我们都不用去管
	//如果把这里改成1000, 就有1000个然在不断的打印,但是不是每个人都有机会在这一毫秒中打印出来,但是很多人都在打印
	//那么这个10和1000有什么关系呢,我们熟悉操作的话,应该知道,我开10个线程没问题,我开100个也还好,但是也差不多了,如果1000个
	//要1000个人并发的去执行一些事情,这就不能通过线程来做到,我们在其他的语音中,要通过异步IO的方式来做到1000个人同时并发的执行
	//但是在go语言当中呢, 我们不用去管,10个也行,100个也行,1000个也行,都可以并发的执行
	//goroutine是一种协程,或者说是和协程是比较像的,协程呢我们叫做是coroutine,这个是其他所有编程语言中都有这个叫法,但是不是所有语言都支持
	//协程是什么
	//协程是一个轻量级的线程,他的作用表面上看和线程差不多,都是并发的执行一些任务的,但是是轻量级的
	//协程为什么是轻量级的
	//非抢占式多任务处理,由协程主动交出控制权,线程就不同了,线城是抢占式的多任务处理,没有主动该控制权, 在任何时候都会被操作系统切换.
	// 哪怕一个语句执行到一半都会被系统从中间掐掉,然后转到其他任务里面去做
	//但是协程就不同了,协程是非抢占式的,什么时候交出控制权呢, 是由我协程内部主动决定的,正式因为这个非抢占式,协程才能做到轻量级
	//抢占式,就需要做到最坏情况,比如我去抢的时候,别人正做到一半,那上下文我们存起来更多的东西.
	//非抢占式,我们 只要集中处理切换的这几个点就可以了,这个对资源的消耗就会小很多
	//协程是 编译器/解释器/虚拟机层面的多任务,他不是错做系统层面的多任务,操作系统还是其中只有线程,没有协程
	//具体在执行上,操作系统有自己的调度器, 我们GO语言有一个调度器会来调度我们的协程,因此多个协程可能在一个或多个线程上运行,这个是有调度器来决定的

	//这个程序表面上看,和抢占式没有什么区别,每个人都会打印,打到一半都会跳出来换成别人,这是 因为这个Printf是一个IO的操作,他在IO的操作里面会有一个协程的切换切换
	//IO操作总会有一个等待的过程
	//我们想办法不让他切换,我们不打印了,我们去开一个数组,在匿名函数中,不断的加自己对应下标的值做++,然后打印数组,这是再运行就会成为一个死循环
	//因为我们之前我们在匿名函数中fmt.Printf("Hello from "+"goroutine %d\n", i)这个是IO操作,IO操作里面是会有协程之间的切换
	//但是这里的a[i]++,我们在改成这个之后,这个把对应的a[i]锁对应的这个内存块的数字去做一个+1,这个a[i]++呢,他的中间没有机会去做协程之间的一个切换
	//那这样的话我们就会被一个协程所抢掉,那这个一个协程如果不主动交出控制权的话,那他就会始终在这个协程里面
	//同样我们的main()也是一个goroutine,其他的goroutine是他来开出来的,他里面的那个time.Sleep(time.Millisecond),因为没有人交出控制权,所以用于sleep不出来
	//这时候我们可以使用runtime.Gosched()函数来手动交出控制权,让别人也有机会去运行,那我们的调度器总有机会调度给我自己,大家有让别人运行之后,才可以一起并发的执行
	//但是一般我们很少会用到runtime.Gosched()函数,这里只是做一个演示,一般我们都要其他机会去进行切换

	var a [10]int //这里改成数组,尝试不等到IO,直接输出所有的值
	for i := 0; i < 10; i++ {
		//匿名函数,不是什么新概念,后面的(i)代表把外面i穿进去,前面的func(i int)代表这个函数的参数要求和类型
		//如果不加go,就是从0-10 反复的去调用这个匿名函数,然后这个函数里面是没有退出条件的,所以就是一个不断打印的死循环
		//加了go之后,就不是不断的调这个函数,而是并发的去执行这个函数,因此我主程序还在跑,然后我并发的去开了一个函数,然后开出来的函数不断的打印hello
		//这就相当于我们开了一个线程,实际上我们开的不是线程,是一个协程,虽然现在看起来差不多,所以主程序会继续往下跑,因此我们就开了10个goroutine,不断地去打印

		//但是这时候如果我们不在主进程中加time.Sleep,运行这个程序的时候,什么都不会打印,而是直接就退出了
		//退出的原因是,因为我们的main()和匿名函数是并发执行的,匿名函数他还来不及打印,我们的main()就从0-10就已经for循环完了,然后我们的main就退出了
		//go语言的程序,一旦main退出了,所有的goroutine就被杀掉了,所以还没来得及打印东西就被杀掉了
		//所以我们需要再main()进程中加一个time.Sleep(time.Millisecond) 让主进程不着急退出

		//这里这个i为什么要单独定义传进去呢,
		//不定义的话就是引用了外面的这个i,不定义的话运行该代码就会提示index out of range,
		//我们在外面命令做了i<10的时候才会穿进去,但是为什么会out of range 呢,如果看不出来也没关系,可以用命令行来做测试
		//go run goroutine.go 运行该文件的时候会提示out of range 的错误
		//通过go run -race goroutine.go检测数据访问的冲突
		//然后会输出 WARING: DATA RACE 在下面有Read ad 0x00000 代表写,Previous write in 0x0000  代表写入 而且还会标识出代码的行号
		//别的地方也没有写入数据 应该错误的就是i造成的,如果想要证明这一点, 我们可以打印i的内存地址来证明
		//不把i穿进去呢,就是函数是编程的概念, 这里这个函数就是一个闭包,他呢就引用了外面的这个i,那for循环外面的i和里面的i是同一个i
		//因此外面的这个i不断的去做加法,最后当外面的这个i跳出以后,调到time.sleep()的时候i已经变成了10,我们的条件是小于10跳进去,最后一次等于10跳出来
		//因此呢,这个最终会是10,当i加到了10以后呢a[i]++就会去引用10,所以就会出错,我们在命令行运行时用-race就是检测出这个错误,就知道我们为什么会出这个错
		//因为我们要把这个i拷一份给这个goroutine,每一个goroutine他都要自己固定下来这个i, 通过传值的形式让他固定下面.
		//函数里面func(i int)我叫i是便于读,值一旦穿进去叫什么都可以后面的()代表的传如的值,把外面的指定变量传入进去
		//最后我们运行,发现没有错误.但是通过go run -race  goroutine.go 发现还是有data race 的waring
		//我们通过行号和内存你地址发现,是一边a[i]++在并发不断写入,而一边fmt.Println(a)在输出,这个给检测了出来,这个问题需要我们使用channel来解决
		go func(i int) {
			for {
				//这里这个i,如果直接用for循环里的i会不安全,所以我们把i通过传值的形式,传进来
				//fmt.Printf("Hello from "+"goroutine %d\n", i)
				//这里改成数组,尝试不等到IO,直接输出所有的值
				a[i]++
				runtime.Gosched() //手动交出控制权,具体为什么可以看var a [10]int上面的注释
			}
		}(i) //这里是传值
	}

	time.Sleep(time.Millisecond)
	//打印++后的数组
	fmt.Println(a)
}

//协程Coroutine
//子程序是协程的一个特例,我们所有的函数调用,都可以看做是一个子程序,所有的这些函数调用都是协程的一个特例
//协程是比子程序更加宽泛的一个概念 - Donnald Knuth说的
//普通函数
//一个main调用函数,dowork,这个dowork做完之后,才会将控制权交还给main函数,然后main函数去执行下一个语句
//协程
//也是main和dowork他不是一个单项的箭头,而是双向的通道,main和dowork之间可以双向的流通,他不知数据,他的控制权也可以双向的流通
//就像我们并发执行两个线程,各做各的,并且两个人可以相互的通信而且呢控制权可以互相的交换给彼此
//那这个main和dowork运行在哪里呢,有可能在同一个线程也可能是多个线程,这个不用我们操心,反正我们就开两个协程让他们通信,我们的调度器很可能把他们放在同一个线程之内
//这样才有可能做到,我们开了1000个协程在我们的机器上运行
//go语言的协程
//go语言的进程开起来之后,会有一个调度器,调度器就负责调度协程,有些协程可能会放在一个线程里面,有些可能是两个放在一个线程里面或者很多个协程放在一个线程里面,这个我们不用管,有调度器来控制
//定义goroutine 只要我们在函数调用前面加上go关键字就能将该函数交给调度器运行,就变成了一个协程
//不需要在定义时区分是否是异步函数
//调度器会在合适的点进行切换,他虽然是一个非抢占式的,但是我们后面还是有一个调度器来进行切换,这些切换的点我们并不能完全的进行控制,这也是传统协程的一点区别
//这一些和传统协程的一点区别,传统的我们都需要在所有的切换的点,显示的写出来,但是goroutine不一样,我们就需要像写普通的函数一样,goroutine调度器会进行切换,但是又和线程的切换不一样
//使用go run -race来检测数据访问冲突的点
//goroutine可能会切换的点
//I/O,select
//channel
//等待锁,goroutine也是有锁的
//函数调用时(有时),函数调用时会有机会切换,但是是否切换呢,由调度器来做决定
//runtime.Gosched() 这个函数是手动切换的一个点,在这一点愿意交出我们的控制权
//上面的几个只是作为参考,不能保证切换,不能保证在其他地方不切换
//虽然这么说,我们的这个goroutine还是非抢占式的,从代码来看他和抢占式的有点像,但是它运行的这个机制还是非抢占式的,比如我们之前不断的做a[i]++的话他就没有机会切换了,进程就死掉了
//我们看看开1000个协程的话,我们的系统到底开了几个线程
//将i<10改成i<1000把a[i]++改成原来的fmt.Printf("Hello from "+"goroutine %d\n", i) 将main函数中的输出删掉,sleep改成1分钟
//通过top命令查看__goroutine进程在th字段也就是线程来看/前面是总线程数/后面是活跃的线程,一般/后面CPU有几个核心就会活跃几个,他觉得没有必要开的超过CPU的核心数,就做了自己的调度
//通过这个我们就可以看到,虽然我们开了1000个gorouteine但是他会映射到我们的几个物理线程上去执行,那么我们后面的调度器会进行调度
