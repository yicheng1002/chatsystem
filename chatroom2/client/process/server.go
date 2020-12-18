package process

import(
	"fmt"
	"os"
	"go_code/chatroom2/client/utils"
	"net"
	"go_code/chatroom2/message"
	"encoding/json"
)

func ShowMenu()  {
	fmt.Println("--------恭喜***登录成功--------")
	fmt.Println("--------1.显示在线用户列表--------")
	fmt.Println("--------2.发送消息--------")
	fmt.Println("--------3.信息列表--------")
	fmt.Println("--------4.退出系统--------")
	fmt.Println("请选择（1-4）")
	var key int
	fmt.Scanf("%d\n",&key)
	switch key {
	case 1:
		outputOnlineUser()
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择了退出系统")
		os.Exit(0)
	default :
		fmt.Println("你输入的选项不正确，请重新输入")
	}
}

//用于和服务器端保持通信
func serverProcessMes(conn net.Conn){
	tf := &utils.Transfer{
		Conn : conn,
	}
	for{
		fmt.Println("客户端正在等待读取服务器端发送来的消息")
		mes,err := tf.ReadPkg()
		if err != nil {
			fmt.Println("客户端与服务器端保持通信报错 err：",err)
	}
	//打印一下从服务器端返回的消息
	fmt.Printf("mes=%v\n",mes)

	//如果从服务器端读到消息，则做下一步处理
	switch mes.Type {
		case message.NotifyUserStatusMesType :
			var notifyUserStatusMes message.NotifyUserStatusMes
			err = json.Unmarshal([]byte(mes.Data),&notifyUserStatusMes)
			if err != nil {
				fmt.Println("json.Unmarshal([]byte(mes.Data),&notifyUserStatusMes) err:",err)
				return
			}
			updateUserStatus(&notifyUserStatusMes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")
	}
	}
	
}