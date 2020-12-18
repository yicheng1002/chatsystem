package main

import(
	"fmt"
	"os"
	"go_code/chatroom2/client/process"
)

func main()  {
	var userid int
	var username string
	var userpwd  string
	//key用来接收用户输入的功能选项
	var key int
	//loop用来存储是否还继续显示系统菜单
	//var loop = true
for true {
	fmt.Println("----------欢迎登录多人聊天系统----------")
	fmt.Println("----------1 登录聊天系统")
	fmt.Println("----------2 注册用户")
	fmt.Println("----------3 退出系统")
	fmt.Println("请选择(1-3):")
	fmt.Println("--------------------")
	fmt.Scanf("%d\n",&key)   //在输入的时候为什么要控制格式，比如此处是%d
	//fmt.Scanf("请输入：","%d\n",&key)  在fmt.Scanf()函数中这样加入汉字的写法是不行的
	//fmt.Scanf()函数的第一个参数用来 设置后面输入的各种形式，所以当把第一个参数设置成 “请输入：”时，
	//没有指定后面即将输入参数的格式，所以导致在键盘（标准输入）输入内容时，不会得到正常的输入

	switch key {
	case 1:
		fmt.Println("登录聊天系统")
		fmt.Println("请输入用户id：")
		fmt.Scanf("%d\n",&userid)
		fmt.Println("请输入用户密码：")
		fmt.Scanf("%v\n",&userpwd)

		up := &process.UserProcess{}
		up.Login(userid,userpwd)
		//err := up.Login(userid,username)
	//调用登录处理函数
	//	loop = false
	case 2:
		fmt.Println("注册用户")
		fmt.Println("请输入用户id：")
		fmt.Scanf("%d\n",&userid)
		fmt.Println("请输入用户密码：")
		fmt.Scanf("%s\n",&userpwd)
		fmt.Println("请输入用户昵称：")
		fmt.Scanf("%s\n",&username)

		up := &process.UserProcess{}
		up.Register(userid,userpwd,username)
	case 3:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("你的输入有误，请重新输入")
	} 

}
	
}