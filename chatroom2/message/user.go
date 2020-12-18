package message

import(

)
//定义一个用户结构体
//为了保证序列化和反序列化成功，我们必须保证 用户信息的json字符串的key 和 结构体的字段对应的tag名字一致
type User struct{
	UserId int `json:"userId"`
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
	UserStatus int `json:"userstatus"` //用户状态
}