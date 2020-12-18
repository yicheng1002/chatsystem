package utils

import(
	"fmt"
	"net"
	"go_code/chatroom2/message"
	"encoding/binary"
	"encoding/json"
)

//作为服务器端和客户端交互的结构体，必然有 连接 字段、数据传输的缓存 字段
type Transfer struct{
	Conn net.Conn
	Buf [8086]byte  //数据传输的缓存字段，并不是真实需要处理的数据，与UserProcess结构体做区别，
	                //两个结构体中都没有实际消息的字段，声明的字段都相当于是工具功能
}

//专门写一个函数用来表示 服务器端读取客户端发来信息的过程，得到的结果是一个Message消息类型
//在ReadPkg方法中，是怎么区分先读取客户端发来的长度，又读取客户端发来的信息的呢？？？
//是因为在WritePkg方法中，先发送的是信息的长度，再发送信息本身内容么？
func (this *Transfer) ReadPkg() (mes message.Message,err error)  {
	_,err = this.Conn.Read(this.Buf[:4]) //读取客户端发来的数据，存到长度为4的buf[:4]中
	//n,err := conn.Read(buf) //读取客户端发来的数据，存到buf切片中
	//n,err := conn.Read(buf[:3])
	if  err != nil {   //把n != 4 这个或条件去掉，虽然buf[:4]的长度是4，但服务器端不一定非得读进4个字节
		fmt.Println("读取客户端数据错误 conn.read err：",err)
		return  //当报错了没有添加return结束程序，则会一直在这里循环
					//如果报错了，已经在这里加了return了，为什么退出不了程序呢？？？难道因为这只是一个协程，
					//如果是在主线程的return，那应该就直接退出程序了
	}
	fmt.Println("\n服务器端读取客户端发来的消息，用切片表示的长度是：",this.Buf[:4])

//根据buf[:4]转成一个uint32类型 =>切片转成整形类型uint32到底是一个怎样的转化过程，什么形式转化成什么形式？？？？？
//为了接下来调用conn.Read()方法时，同样要传参一个切片类型，但现在得到的客户端传来消息长度的表现形式是字节切片buf[:4]
//在调用conn.Read()方法时传切片类型时，总不能conn.Read(buf[:buf[:4]])吧，显然不能这样写，
//所以要把用切片表示的长度用整形类型表示
	var pkgLen uint32     
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])

	n,err := this.Conn.Read(this.Buf[:pkgLen])
	if n!=int(pkgLen) || err != nil{
		fmt.Println("读取客户端发来的信息有误 conn.Read(buf[:pkgLen]) err：",err)
		return
	}
	//把存放在切片buf[:pkgLen]中的来自客户端的消息进行反序列化
	err = json.Unmarshal(this.Buf[:pkgLen],&mes)  //把用buf[:pkgLen]字节形式的内容反序列化后存放到message.Message类型的变量mes中
    if err != nil {
		fmt.Println("来自客户端的消息反序列化失败 son.Unmarshal(buf[:pkgLen],&mes) err",err)
		return
	}

	fmt.Println("从客户端读到的数据是：",string(this.Buf[:pkgLen]))   //为了调试加的语句
	return 
}

func (this *Transfer) WritePkg(data []byte) (err error){
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	//binary包中的BigEndian变量是bigEndian类型的，bigEndian类型就是ByteOrder空接口类型
	//调用ByteOrder里面的PutUnit32()方法，把unit32类型表示的长度用byte类型切片表示
	binary.BigEndian.PutUint32(this.Buf[0:4],pkgLen)//把unint32类型的pkgLen用byte类型的切片表示

	//发送长度
	n,err := this.Conn.Write(this.Buf[0:4])//把buf[0:4]切片里面的内容（即消息长度）发送给服务器端
    if n!= 4 ||err != nil{
		fmt.Println("conn.Write(buf[0:4]) fail",err)
		return
	}
	//发送消息本身内容
	n,err = this.Conn.Write(data)//把buf[0:4]切片里面的内容（即消息长度）发送给服务器端
    if n!= int(pkgLen) ||err != nil{
		fmt.Println("conn.Write(buf[0:4]) fail",err)
		return
}
	return
}