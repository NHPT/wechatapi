package main

import(
    "fmt"
    "github.com/NHPT/wechatapi"
    "github.com/buger/jsonparser"
    "time"
    "strings"
)
func main() {
    // 登录
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err)
        return
    }
    // 获取自己的信息
    myself := wechatapi.GetMySelf()
    if myself.UserName != "" {
        fmt.Println("用户名:",myself.UserName)
        fmt.Println("昵称:",myself.NickName)
        fmt.Println("头像:",myself.HeadImgUrl)
        fmt.Println("性别:",myself.Sex)
        fmt.Println("签名:",myself.Signature)
    }
    // 获取所有联系人
    contactcount, contactlist := wechatapi.GetContactlist()
    fmt.Println("联系人数量:",contactcount)
    fmt.Println("联系人信息:",string(contactlist))
    chatroomcount, chatroomlist := wechatapi.GetChatrooms()
    fmt.Println("群聊数量:",chatroomcount)
    // 解析群聊信息
    jsonparser.ArrayEach(chatroomlist,func(value []byte, dataType jsonparser.ValueType, offset int, err error){
        if err != nil {
            return
        }
        chatroom,_ := jsonparser.GetString(value,"NickName")
        // 输出群昵称
        fmt.Println(chatroom)
    })
    // 获取所有公众号
    offc, offclist :=  wechatapi.GetOfficiallist()
    fmt.Println("公众号数量:",offc)
    fmt.Println("公众号信息:",string(offclist))
    // 根据备注获取联系人信息
    cb := wechatapi.GetContactByRemarkName("团团")
    if len(cb) > 0 {
        fmt.Println("联系人信息:",string(cb))
        uname, _ := jsonparser.GetString(cb,"UserName")
        fmt.Println(uname)
    }
    // 根据备注获取联系人信息，不存在则返回空
    other := wechatapi.GetContactByRemarkName("zh")
    if len(other) > 0 {
        fmt.Println("联系人信息:",string(other))
    }
    // 根据昵称获取联系人信息
    nc := wechatapi.GetContactByNickName("张三")
    if len(nc) > 0 {
        fmt.Println("联系人信息:",string(nc))
    }
    // 获取北京的联系人
    ct,ctlist := wechatapi.GetContacts("北京","","")
    fmt.Println("北京联系人数量:",ct)
    fmt.Println(string(ctlist))
    // 获取西安的男性联系人
    ct2,_ := wechatapi.GetContacts("","西安","1")
    fmt.Println("西安的男性数量:",ct2)
    // 获取所有男性联系人
    ct3,_ := wechatapi.GetContacts("","","1")
    fmt.Println("男性联系人数量:",ct3)
    // 获取所有女性联系人
    ct4,_ := wechatapi.GetContacts("","","2")
    fmt.Println("女性联系人数量:",ct4)
    // 接收消息
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
            // fmt.Println("[+]",msgtime,"收到群[",msg.Srcname,"]成员[",msg.Usernick,"]的消息:",msg.Msg)
            // 细化消息类型 
            if msg.Msgtype == 1 {
                fmt.Println("[+]",msgtime,"收到群[",msg.Srcname,"]成员[",msg.Usernick,"]的文本消息:",msg.Msg)
            }
            if msg.Msgtype == 3 {
                fmt.Println("[+]",msgtime,"收到群[",msg.Srcname,"]成员[",msg.Usernick,"]的图片消息:",msg.Msg)
            }
            if msg.Msgtype == 47 {
                fmt.Println("[+]",msgtime,"收到群[",msg.Srcname,"]成员[",msg.Usernick,"]的表情包消息:",msg.Msg)
            }
            if msg.Msgtype == 49 {
                fmt.Println("[+]",msgtime,"收到群[",msg.Srcname,"]成员[",msg.Usernick,"]转发的公众号消息:",msg.Msg)
            }
            if msg.Msgtype == 10000 {
                fmt.Println("[+]",msgtime,"收到群[",msg.Srcname,"]成员[",msg.Usernick,"]的拍一拍消息:",msg.Msg)
            }
        }
        // 公众号消息
        if msg.Srctype == 2 {
            fmt.Println("[+]",msgtime,"收到公众号[",msg.Srcname,"]的消息:",msg.Msg)
        }
        //联系人消息
        if msg.Srctype == 3 {
            fmt.Println("[+]",msgtime,"收到联系人[",msg.Usernick,"]的消息:",msg.Msg)
            // 自动邀请进群
            if strings.Contains(msg.Msg,"加群") || strings.Contains(msg.Msg,"进群") || strings.Contains(msg.Msg,"入群") || strings.Contains(msg.Msg,"社群"){
                err := wechatapi.JoinChatroom(msg.Fromuser,"Pentest技术交流")
                if err!= nil {
                    fmt.Println("[!] 邀请加入群聊错误")
                }
                _ = wechatapi.SendMsg("期待您的加入！[握手]",msg.Fromuser)
            }
        }
    }
    // 退出登录
    //wechatapi.Logout()
}
