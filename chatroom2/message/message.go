package message

const(
	LoginMesType = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
)


//定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct{
	Type string `json:"type"`
	Data string `json:"data"`
}

//还需要写一个注册结构体类型
type LoginMes struct{
	Userid int `json:"userid"`
	UserPwd string `json:"userpwd"`
	UserName string `json:"username"`
}

type LoginResMes struct{
	Code int `json:"code"`
	UsersId []int //增加字段，保存登录成功的用户id，由于登录成功的用户数量是在动态变化的，所以这里声明一个int类型的切片
	Error string `json:"error"`
}

type RegisterMes struct {   
	User User `json:"user"`
}

type RegisterResMes struct {
	Code int `json:"code"`
	Error string `json:"error"`
}

type NotifyUserStatusMes struct{
	UserId int `json:"userId"`  //用户id
	Status int `json:"status"`  //用户的状态
}