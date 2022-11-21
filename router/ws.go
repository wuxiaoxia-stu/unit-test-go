package router

// ws消息体格式定义
type Message struct {
	Id      string // 消息uuid
	From    int    // 消息发送者
	To      int    // 消息接受者（用户id， 房间id）
	Type    int    // 1：文本消息  2：单图消息   3：语音消息  4：图文消息
	Content string // 消息内容
	Time    int
}

//  接口定义
//  1.登录
//  2.连接
//  3.发送消息
//  4.加入房间
//  5.退出房间
//  6.群发
//  7.广播
