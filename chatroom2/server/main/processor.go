package main

import(
	"fmt"
	"net"
	"go_code/chatroom2/server/utils"
	"go_code/chatroom2/message"
	"go_code/chatroom2/server/process"
	"io"
)

type Processor struct{
	Conn net.Conn
}

func (this *Processor) process2() (err error){
	for{
		tf := &utils.Transfer{  //为什么要这样写呢？用 tf去调用绑定在 Transfer类型的方法ReadPkg
			Conn : this.Conn,    //因为ReadPkg()是绑定在Transfer类型上的方法，必须要用这个类型的变量去调用这个方法才能使用
		}

		//没有给tf指定Buf字段值，输出一下tf的值，看是怎么现实的：
		//如果要输出结构体类型的值，以下3种方法都可以，
		//总结：fmt.Printf()可直接在参数内指定输出格式，例如%v，%+v
		//fmt.Println()不能在参数内直接指定输出格式，但可利用fmt.Sprintf()函数返回格式化的值
		fmt.Printf("tf的值是%+v",*tf)
		//②fmt.Println("tf的值是：",*tf)
		//③fmt.Println("tf的值是：",fmt.Sprintf("%v",*tf))
		//上面三行输出的tf的值是：
		//{Conn:0xc04206c040 Buf:[0 0 0 0 0 0 0 0 0 0 0......0 0 0]}
		//由于在初始化utils.Transfer类型的值tf时，没有指定Buf字段的值，则根据[8096]byte类型默认初始值为8096个0

		mes,err := tf.ReadPkg()
		if err != nil{
			if err == io.EOF{
				fmt.Println("读不到客户端的内容了，客户端退出，服务器端也退出..")
				return err
			}else{
				fmt.Println("服务器端读取客户端数据有误 readPkg(conn) err",err)
				return err
			}			
		}

		//从这里根据客户端和服务器端之间的数据传输通道读到的消息的类型不同，服务器端对消息做不同的处理
		err = this.serverProcessMes(&mes)
		if err != nil {
			fmt.Println("服务器端处理客户端发来的消息失败 serverProcessMes(conn,&mes) err：",err)
			return err
		}
		//fmt.Println("从客户端读到的消息mes是：",mes)
	}
}
	
func (this *Processor) serverProcessMes(mes *message.Message)(err error)  {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		up := &process2.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up := &process2.UserProcess{
			Conn : this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}