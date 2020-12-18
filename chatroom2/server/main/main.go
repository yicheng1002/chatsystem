package main

import(
	"fmt"
	"net"
	"time"
	"go_code/chatroom2/server/model"
)

//服务器端启动一个协程为 服务器端和客户端两者连接的通道传输服务，第一步就先延时关闭这个数据传输通道
func process(conn net.Conn)  {   //net.Conn类型 是 listen.Accept()返回的类型
	//延时关闭服务器端与客户端通信的通道
	defer conn.Close()
	//在这里调用总控方法
	processor := Processor{
		Conn : conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("调用总控程序错误 err:",err)
	}


}

func initUserDao() {  //不需要返回值，只是对model包中声明的全局变量MyUserDao进行初始化
	model.MyUserDao = model.NewUserDao(pool) 
}

func main()  {
	//当服务器启动时，我们就去初始化redis的连接池
	initPool("localhost:6379",16,0,300*time.Second)
	//当有了连接池，我们再初始化一个UserDao实例
	initUserDao()


	fmt.Println("服务器在8889端口监听...")
	//listen,err := net.Listen("tcp","0,0,0,0:8889") //纠错处理，不是0,0,0,0  而是0.0.0.0
	listen,err := net.Listen("tcp","localhost:8889")  //net包的Listen函数的返回值是Listener接口类型，
	if err != nil{      //这个类型的接口中定义了3个方法，然后调用其中的Accpet()方法 得到一个连接到该服务器端口的连接
		fmt.Println("服务器监听8889端口出错net.Listen err:",err)
		return
	}

	//服务端一旦监听8889成功，则开始等待客户端来连接服务器
	for{
		fmt.Println("等待客户端来连接服务器...") 
		conn,err := listen.Accept()  //等待并返回下一个连接到该接口（服务器监听端口）的连接
		if err != nil{               //listener接口类型的变量listen调用Accept()函数返回Conn类型的变量
			fmt.Println("等待客户端来连接服务器listen.Accept err",err)
			return
		}
		fmt.Println("连接到服务器的通道是什么样的输出:",conn)  //调试专用,conn通道的值原来是一个地址
		//客户端一旦连接服务端成功，则启动一个协程和客户端保持通讯，process就是该协程在做的事情
        go process(conn)
	}
}