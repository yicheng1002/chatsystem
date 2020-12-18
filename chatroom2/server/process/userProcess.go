package process2

import(
	"fmt"
	"net"
	"go_code/chatroom2/message"
	"go_code/chatroom2/server/utils"
	"go_code/chatroom2/server/model"
	"encoding/json"
)

type UserProcess struct{
	Conn net.Conn
	//增加一个字段，表示该连接Conn是哪个用户(客户端)连接过来的
	UserId int
}


//这个函数的目的是服务器端处理客户端发来的注册消息，并返回给客户端对信息的处理结果
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data),&registerMes)
	if err != nil{
		fmt.Println("注册功能：服务器端反序列化客户端发来的数据时失败 json.Unmarshal([]byte(mes.Data),&registerMes) err",err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	var registerResMes message.RegisterResMes
	//从这里开始，服务器端拿着客户端发来的UserId、UserPwd、UserName到redis中验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil{
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		}else{
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}
	}else{
		registerResMes.Code = 200
	}

	data,err := json.Marshal(registerResMes)
	if err != nil{
		fmt.Println("注册功能：服务器端序列化发送到客户端的实际处理消息失败 json.Marshal(registerResMes) err",err)
		return
	}

	//将data的值赋值给resMes
	resMes.Data = string(data)
	data,err = json.Marshal(resMes)
	if err != nil{
		fmt.Println("注册功能：服务器端序列化发送到客户端的处理消息失败 json.Marshal(resMes) err",err)
		return
	}

	//把data(服务器端处理客户端消息的结果)发送给客户端
	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil{
		fmt.Println("注册功能：服务器端发送处理结果到客户端的失败 tf.WritePkg(data) err",err)
		return
	}
	return
}


//这个函数的目的是服务器端处理客户端发来的登录信息，并返回给客户端对信息的处理结果
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data),&loginMes)
	if err != nil{
		fmt.Println("登录功能：服务器端反序列化客户端发来的数据时失败 json.Unmarshal([]byte(mes.Data),&loginMes) err",err)
		return
	}
	//resMes这个变量是服务器端返回给客户端的处理消息
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes
	// if loginMes.Userid == 100 && loginMes.UserPwd == "123456" {
	// 	//合法
	// 	loginResMes.Code = 200
	// }else{
	// 	loginResMes.Code = 500
	// 	loginResMes.Error = "该用户不存在，请注册后再使用"
	// 	//fmt.Println("该用户不存在，请注册后再使用")  //终于找到问题了，我在这里没有设置loginResMes.Error的值，只打印出来了
	// }


	//我们需要将从客户端拿到的Userid、UserPwd到redis中进行验证
	user,err := model.MyUserDao.Login(loginMes.Userid,loginMes.UserPwd)    //这里的Userid、UserPwd字段格式只与LoginMes结构体定义的形式有关
	if err != nil{                                           //在model包中定义的结构体user的tag与保存到redis中的json字符串形式有关
		if err == model.ERROR_USER_NOTEXISTS{
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		}else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		}else{
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误.."
		}
	}else{
		loginResMes.Code = 200

		//这里主要是给UserProcess结构体类型的实例this赋值，然后把这个实例添加到UserMgr结构体实例userMgr中
		//可以这样理解：对于成功连接到服务器端的客户端，要把它添加到在线用户列表的结构体实例userMgr中，
		//用UserProcess结构体类型的变量来表示连接到服务器的客户端
		this.UserId = loginMes.Userid
		fmt.Println("userprocess 的this的值:",this)  //打印一下UserProcess的值，这是调试语句
		userMgr.AddOnlineUser(this)

		// //遍历当前在线用户userMgr，放入到loginResMes.UsersId切片中，返回给客户端
		// for id,_ := range userMgr.onlineUsers{
		// 	loginResMes.UsersId = append(loginResMes.UsersId,id)
		// }
		
		fmt.Println(user,"登录成功")
		
		//当我(loginMes.Userid)上线成功后，通知其他用户我上线了
		this.NotifyOtherOnlineUser(loginMes.Userid)	
	}

	//3.将loginResMes序列化
	data,err := json.Marshal(loginResMes)
	if err != nil{
		fmt.Println("序列化服务器端登录处理返回信息失败 json.Marshal(loginResMes) err",err)
		return
	}
	//fmt.Println("打印data",data)
	//4.将登录返回消息loginResMes也就是data赋值给resMes的Data
	resMes.Data = string(data)

	//5.对resMes序列化
	data,err = json.Marshal(resMes)
	if err != nil{
		fmt.Println("序列化服务器端返回信息失败 json.Marshal(resMes) err",err)
		return
	}

	//6.发送data到客户端
	tf := &utils.Transfer{
		Conn : this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil{
		fmt.Println("服务器端发送数据到客户端失败 writePkg(conn,data) err",err)
		return
	}
	return  //这个地方有必要填么
}

//此方法是在服务器端通知所有在线用户 ：userId表示的用户上线了
//UserMgr是用来处理在线用户的结构体，而UserProcess是处理跟用户相关方法的结构体
func (this *UserProcess) NotifyOtherOnlineUser(userId int)  {
	//遍历userMgr.onlineUsers这个map类型的变量，
	for id,up := range userMgr.onlineUsers {
		if id == userId{
			continue
		}
		//调用方法开始正式 把我（userId）上线的消息发送给 所有的在线用户
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int)  {
	//开始组装NotifyUserStatusMes类型的消息
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data,err := json.Marshal(notifyUserStatusMes)
	if err != nil{
		fmt.Println("将notifyUserStatusMes序列化失败 json.Marshal(notifyUserStatusMes) errr：",err)
		return
	}

	//将序列化的notifyUserStatusMes复制给mes.Data
	mes.Data = string(data)

	//对mes再次进行序列化
	data,err = json.Marshal(mes)
	if err != nil{
		fmt.Println("将mes序列化失败 json.Marshal(mes) errr：",err)
		return
	}

	//把构造好的mes发送给客户端
	tf := &utils.Transfer{
		Conn : this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("往客户端发送 我在线的消息失败 tf.WritePkg(data) err：",err)
		return
	}

}