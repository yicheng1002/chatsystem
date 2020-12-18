package process2

import(
	"fmt"
)
//因为UserMgr结构体的实例在服务器端只有一个，用来维护在线用户列表
//而且在很多地方都要用到（例如客户端？？），所以直接声明为全局变量
var (
	userMgr *UserMgr
)
type UserMgr struct{
	onlineUsers map[int]*UserProcess
}

//对userMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers : make(map[int]*UserProcess,1024),
	}
}

//userMgr变量表示了连接到服务器端的客户端，现在分别完成对userMgr变量的增删改查
//增加表示连接到服务器端的客户端数量增加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}
//删除
func (this *UserMgr) DeleteOnlineUser(userId int)  {
	delete(this.onlineUsers,userId)
}

//返回当前所有在线的用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}
//根据userId，查询userId这个用户是否在线
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess,err error)  {
	up,ok := this.onlineUsers[userId]
	if !ok {//你要查找的这个用户，当前不在线
		err = fmt.Errorf("用户%d不存在",userId)
		return
	}
	return
}


