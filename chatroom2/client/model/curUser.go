package model

import (

)

//维护当前的链接，保持链接能正常向服务器端发送消息
type CurUser struct {
	Conn net.Conn
	message.User //匿名User结构体类型
}



