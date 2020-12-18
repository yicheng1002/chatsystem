package model

import(
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"fmt"
	"go_code/chatroom2/message"
)

//定义一个UserDao的结构体
type UserDao struct{
	pool *redis.Pool  //为什么要把连接池字段声明为指针类型的呢？为什么要把连接池字段设置为私有类型呢
}

//使用工厂模式，创建一个UserDao的实例，
// func NewUserDao(pool *redis.Pool)(userDao *UserDao) {
// 	userDao = &UserDao {
// 		pool : pool,
// }	
// 		return
// }
//先实验一下返回UserDao结构体类型的变量  和 返回UserDao结构体指针类型的变量 的区别
func NewUserDao(pool *redis.Pool)(userDao *UserDao) {
	userDao = &UserDao {
		pool : pool,
}	
		return
}

//我们在启动服务器时，就初始化一个userDao实例
//这里只是声明了一个UserDao结构体类型的变量MyUserDao,但并没有对此变量初始化，在main中初始化吧
var (
	MyUserDao *UserDao  //为什么要声明UserDao指针类型的结构体
)

//开始写UserDao应该提供哪些方法
//调用此方法时，根据传进来的id判断redis库中是否有对应的用户信息
func (this *UserDao)getUserById(conn redis.Conn,id int) (user *message.User,err error) {
	//通过id查redis库中是否有这个用户
	res,err := redis.String(conn.Do("HGet","users",id))
	if err != nil{
		if err == redis.ErrNil{   //ErrNil表示在哈希users中没有找到的id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &message.User{}
	err = json.Unmarshal([]byte(res),user)  //字符串类型的值可以直接强转为切片类型
	if err != nil {
		fmt.Println("哈希users中id对应的用户信息反序列化失败 err：",err)
		return
	}
	return
}

//经过getUserById函数后，得到一个User结构体，现在写方法验证服务器端得到的id和pwd是否和redis中存在的id和pwd一致
//getUserById函数是从redis库中拿到字符串后，把字符串解析成User结构体实例，然后用此实例与服务器从客户端拿到的id、pwd做对比
//1.Login函数完成对用户的验证
//2.如果用户的id和pwd都正确，则返回一个user实例
//3.如果用户的id或pwd有错误，则返回对应的错误信息
func (this *UserDao) Login(userId int,userPwd string) (user *message.User,err error) {
	//从UserDao中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()
	user,err = this.getUserById(conn,userId)
	if err!= nil {
		fmt.Println("redis库中并不存在userId对应的User类型的字符串 err：",err)
		return
	}
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {
	//在model模块取一根连向redis库的连接
	conn := this.pool.Get()
	defer conn.Close()
	_,err = this.getUserById(conn,user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
	}
	//如果err是空的，则说明从redis中取出了user.UserId对应的数据，
	//如果err不是空的，即有错误，没有从redis取出来，这时，则可以正式往redis中注册
	data,err := json.Marshal(user)
	if err != nil {
		fmt.Println("即将注册到redis中的user序列化失败 err：",err)
		return
	}
	_,err = conn.Do("HSet","users",user.UserId,string(data))
	if err != nil {
		fmt.Println("保存注册用户到redis失败 err：",err)
		return
	}
	return
}
