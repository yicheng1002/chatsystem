package process

import (
	"fmt"
	"go_code/chatroom2/message"
)

//客户端维护一个全局的map类型的变量  用来存放我本人所能看到的在线用户
var onlineUsers map[int]*message.User = make(map[int]*message.User,10)

//在客户端显示当前在线的用户
func outputOnlineUser() {
	//遍历onlineUsers变量，
	fmt.Println("当前在线用户列表：")
	for id,_ := range onlineUsers {
		fmt.Println("在线用户id：",id)
	}
}

//编写一个方法 处理从服务器端返回的notifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user,ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId : notifyUserStatusMes.UserId,
		}
	}
	//当客户端的在线用户列表onlineUsers中有和notifyUserStatusMes.UserId一样的user时，则更新此user的状态
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}

