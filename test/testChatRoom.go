package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/NHPT/wechatapi"
	"github.com/buger/jsonparser"
)

func main() {
	// 登录
	err := wechatapi.Login()
	fmt.Println(err, "登录信息")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取自己的信息
	myself := wechatapi.GetMySelf()
	fmt.Println(myself, "自己的信息")
	if myself.UserName != "" {
		fmt.Println("用户名:", myself.UserName)
		fmt.Println("昵称:", myself.NickName)
		fmt.Println("头像:", myself.HeadImgUrl)
		fmt.Println("性别:", myself.Sex)
		fmt.Println("签名:", myself.Signature)
	}
	chatroomcount, chatroomlist := wechatapi.GetChatrooms()
	fmt.Println("群聊数量:", chatroomcount)
	// 解析群聊信息
	jsonparser.ArrayEach(chatroomlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		chatroom, _ := jsonparser.GetString(value, "NickName")
		username, _ := jsonparser.GetString(value, "UserName")
		if chatroom == "你的群聊名" {
			fmt.Println(chatroom)
			fmt.Println(username)
			receiveMessages(myself.NickName, username)
			return
		}

	})

	// 退出登录
	//wechatapi.Logout()
}

// 接收消息的函数
func receiveMessages(username string, target string) {
	for {
		msg := wechatapi.GetMsg()
		// 格式化时间
		tm := time.Unix(msg.Createtime, 0)
		msgtime := tm.Format("01-02 15:04:05")
		// 返回1101说明退出登录
		if msg.Srctype == 1101 {
			fmt.Println("[!] 登录失效，请重新登录！")
			break
		}
		// 群消息
		if msg.Srctype == 1 {
			//if msg.Msgtype == 1 {
			//	fmt.Println("[+]", msgtime, "收到群[", msg.Srcname, "]成员[", msg.Usernick, "]的文本消息:", msg.Msg)
			//}
			//if msg.Msgtype == 3 {
			//	fmt.Println("[+]", msgtime, "收到群[", msg.Srcname, "]成员[", msg.Usernick, "]的图片消息:", msg.Msg)
			//}
			//if msg.Msgtype == 47 {
			//	fmt.Println("[+]", msgtime, "收到群[", msg.Srcname, "]成员[", msg.Usernick, "]的表情包消息:", msg.Msg)
			//}
			//if msg.Msgtype == 49 {
			//	fmt.Println("[+]", msgtime, "收到群[", msg.Srcname, "]成员[", msg.Usernick, "]转发的公众号消息:", msg.Msg)
			//}
			//if msg.Msgtype == 10000 {
			//	fmt.Println("[+]", msgtime, "收到群[", msg.Srcname, "]成员[", msg.Usernick, "]的拍一拍消息:", msg.Msg)
			//}
			if msg.Srcname == "你的群聊名" {
				msg.Msg = strings.TrimPrefix(msg.Msg, "<br/>")
				//匹配msg.Msg最开始的几个字符是否为@username
				myname := "@" + username
				if len(msg.Msg) < len(myname) {
					continue
				}
				fmt.Println("@" + username)
				fmt.Println(msg.Msg[:len(myname)])
				fmt.Println(("@" + username) == msg.Msg[:len(myname)])
				if strings.HasPrefix(msg.Msg, "@"+username) {
					fmt.Println("[+]", msgtime, "收到群[", msg.Srcname, "]成员[", msg.Usernick, "]的@消息:", msg.Msg)
					//将msg.Msg去掉最开始的@username
					msg.Msg = strings.TrimPrefix(msg.Msg, "@"+username+" ")
					//将msg.Msg去掉最开始的空格
					msg.Msg = strings.TrimPrefix(msg.Msg, " ")
					sendMessages(target, "@我干嘛")
					fmt.Println(msg.Msg)
				} else {
					fmt.Println(msg.Msg)

				}
			}
		}
		// 公众号消息
		if msg.Srctype == 2 {
			fmt.Println("[+]", msgtime, "收到公众号[", msg.Srcname, "]的消息:", msg.Msg)
		}
		//联系人消息
		if msg.Srctype == 3 {
			fmt.Println("[+]", msgtime, "收到联系人[", msg.Usernick, "]的消息:", msg.Msg)
			if strings.Contains(msg.Msg, "加群") || strings.Contains(msg.Msg, "进群") || strings.Contains(msg.Msg, "入群") || strings.Contains(msg.Msg, "社群") {
				err := wechatapi.JoinChatroom(msg.Fromuser, "Pentest技术交流")
				if err != nil {
					fmt.Println("[!] 邀请加入群聊错误")
				}
				_ = wechatapi.SendMsg("期待您的加入！[握手]", msg.Fromuser)
			}
		}
	}
}

// 发送消息的线程
func sendMessages(username string, message string) {
	// 示例：设定要发送的目标用户和消息内容
	targetUser := username // 替换为实际的用户名或用户ID
	// 发送消息
	err := wechatapi.SendMsg(message, targetUser)
	if err != nil {
		fmt.Println("[!] 发送消息错误:", err)
	} else {
		fmt.Println("[+] 消息已发送给", targetUser, "内容:", message)
	}
}
