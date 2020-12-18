package process

import(
	"fmt"
	"encoding/json"
	"encoding/binary"
	"net"
	"go_code/chatroom2/message"
	"go_code/chatroom2/client/utils"
	"errors"
	"os"
)
type UserProcess struct{
	//暂时不需要什么字段，Login函数里的通道是通过调用net包的函数获得的，属于输出型（获取型）的参数，
	//当需要输入型的参数时，则在结构体的字段中声明
}

func (this *UserProcess) Register(userid int,userpwd,username string) (err error) {
	//1.连接到服务器
	conn,err := net.Dial("tcp","localhost:8889")
	if err != nil {
		fmt.Println("注册功能：客户端连接服务器端错误 err:",err)
		return
	}
	defer conn.Close()
	//2.开始形成发往服务器端的注册消息
	var mes message.Message
	mes.Type = message.RegisterMesType

	var registerMes message.RegisterMes
	registerMes.User.UserId = userid
	registerMes.User.UserPwd = userpwd
	registerMes.User.UserName = username

	data,err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("注册功能：客户端发送消息registerMes序列化错误 err:",err)
		return
	}
	mes.Data = string(data) //data是切片类型的，这里注意要转为字符串类型

	data,err = json.Marshal(mes)
	if err != nil {
		fmt.Println("注册功能：客户端发送消息mes序列化错误 err:",err)
		return
	}
	
	//客户端正式往服务器端发送消息:
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn : conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册功能：客户端发送消息到服务器端错误 err:",err)
		return
	}


	//以上是把注册消息从客户端发往客户端的部分，接下来要写的是客户端从服务器端接收返回来的消息：
	mes,err = tf.ReadPkg()
	if err != nil {
		fmt.Println("注册功能：客户端读取服务器端发来的消息错误 err:",err)
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data),&registerResMes)
	if err != nil {
		fmt.Println("注册功能：客户端读取服务器端发来的消息mes.Data反序列化错误 err:",err)
		return
	}

	if registerResMes.Code == 200 {
		fmt.Println("注册成功")
		//这里是不是应该调用一下注册成功后应该显示的页面
		os.Exit(0)
	}else{
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}

//写一个供主页面注册功能调用的函数Register

//把Login函数绑定到UserProcess结构体上
func (this *UserProcess) Login(userid int,userpwd string) (err error)  {
	// fmt.Printf("您的用户id是%d，密码是%v",userid,username)
	// return nil
	//1.客户端连接到服务器端
	conn,err := net.Dial("tcp","localhost:8889")
	if err != nil {
		fmt.Printf("登录功能：客户端连接服务器端错误 net.dial err\n")
		return  //在调用net.Dial()时得到了err变量的值，所以可以直接return，login函数可以直接得到返回值err
	}
	//延时关闭
	defer conn.Close()
	//2.准备通过conn发送消息到服务器端
	var mes message.Message
	mes.Type = message.LoginMesType
	
	//3.创建一个登录消息用的结构体LoginMes
	var loginMes message.LoginMes
	loginMes.Userid = userid
	loginMes.UserPwd = userpwd

	//4.将loginMes序列化
	data,err := json.Marshal(loginMes)  //json.Marshal()返回byte类型的切片
	if err != nil{
		fmt.Println("json.Marshal(loginMes) err",err)
		return
	}

	//5.把data的值赋值给mes
	mes.Data = string(data)

	//6.将结构体类型的变量mes序列化
	data,err = json.Marshal(mes)
	if err != nil{
		fmt.Println("json.Marshal(mes) err",err)
		return
	}
	//7.到现在为止，就生成了需要发送给服务器端的message：data
	//先把data的长度发送给服务器端
	//定义一个存放长度的变量
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	//binary包中的BigEndian变量是bigEndian类型的，bigEndian类型就是ByteOrder空接口类型
	//调用ByteOrder里面的PutUnit32()方法，把unit32类型表示的长度用byte类型切片表示
	binary.BigEndian.PutUint32(buf[0:4],pkgLen)//把unint32类型的pkgLen用byte类型的切片表示

	//8.发送长度
	n,err := conn.Write(buf[0:4])//把buf[0:4]切片里面的内容（即消息长度）发送给服务器端
    if n!=4 ||err != nil{
		fmt.Println("conn.Write(buf[0:4]) fail",err)
		return
	}
	fmt.Printf("客户端，发送消息的长度=%d，发送的内容是%v\n",len(data),string(data))

	_,err = conn.Write(data)
	if err != nil{
		fmt.Println("客户端往服务器端发送消息失败 conn.Write(data) err：",err)
		return
	}
	//休眠20秒，为什么要休眠20秒？单纯的等待服务器端返回对发送消息的处理结果么？
	//time.Sleep(20 * time.Second)
	//不休眠了 客户端开始处理服务器端返回的消息
	tf := &utils.Transfer{
		Conn : conn,
	}
	mes,err = tf.ReadPkg()
	fmt.Println("测试客户端有没有拿到服务器端返回来的消息：",mes)  //调试语句

	//将mes.Data的字段的值反序列化成LoginResMes类型
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data),&loginResMes) //把字符串类型的mes.Data反序列化为message.LoginResMes类型的数值loginResMes
	if err != nil {                //在json串反序列化时，第二个参数一定要输入一个变量，意思是把json串格式化为这个变量类型的格式么？？？
		fmt.Println("反序列化mes.Data的字符串 json.Unmarshal([]byte(mes.Data),&loginResMes) err：",err)
	}
    fmt.Println("输出返回信息的code值：",loginResMes.Code)
	if loginResMes.Code == 200 {
		//fmt.Println("登录成功")
		fmt.Println("当前在线用户列表如下：")
		for _,v := range loginResMes.UsersId{   //loginResMes是服务器返回给客户端的消息，
			if v == userid{						//该消息中的切片类型的字段UsersId保存了当前登录到服务器端的客户端
				continue
			}
			fmt.Println("用户id：",v)

			//完成 客户端的onlineUsers变量的初始化
			user := &message.User{
				UserId : v,
				UserStatus : message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")

		//这里我们还需要在客户端启动一个协程？？？？？？  client/process/server.go文件中的serverProcessMes只是定义，在这里才是调用
		//该协程用来和服务器端保持通信，如果服务器端有数据推送给客户端
		//则接收并显示在客户端的终端
		go serverProcessMes(conn)

		for{
			ShowMenu()
		}
		
		//err = nil
	}else if loginResMes.Code == 500{
		loginResMes.Error = "该用户不存在，请注册后再使用"
		//err = error(loginResMes.Error)
		err = errors.New("该用户不存在，请注册后再使用")
	}
	return

}
