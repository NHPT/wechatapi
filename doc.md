## 目录

- [结构体](https://github.com/NHPT/wechatapi/blob/main/doc.md#%E7%BB%93%E6%9E%84%E4%BD%93)
- [全局变量](https://github.com/NHPT/wechatapi/blob/main/doc.md#%E5%85%A8%E5%B1%80%E5%8F%98%E9%87%8F)
- [函数](https://github.com/NHPT/wechatapi/blob/main/doc.md#%E5%87%BD%E6%95%B0)
    - [无需登录](https://github.com/NHPT/wechatapi/blob/main/doc.md#%E6%97%A0%E9%9C%80%E7%99%BB%E5%BD%95)
        - [PrintQRCode](https://github.com/NHPT/wechatapi/blob/main/doc.md#printqrcodeinfo-string-error)
        - [Login](https://github.com/NHPT/wechatapi/blob/main/doc.md#login-error)
    - [需要登录](https://github.com/NHPT/wechatapi/blob/main/doc.md#%E9%9C%80%E8%A6%81%E7%99%BB%E5%BD%95)
        - [GetMySelf](https://github.com/NHPT/wechatapi/blob/main/doc.md#getmyself-myself)
        - [GetChatrooms](https://github.com/NHPT/wechatapi/blob/main/doc.md#getchatrooms-int64-byte)
        - [GetContactlist](https://github.com/NHPT/wechatapi/blob/main/doc.md#getcontactlist-int64-byte)
        - [GetOfficiallist](https://github.com/NHPT/wechatapi/blob/main/doc.md#getofficiallist-int64-byte)
        - [GetContactByRemarkName](https://github.com/NHPT/wechatapi/blob/main/doc.md#getcontactbyremarknameremarkname-string-byte)
        - [GetContactByNickName](https://github.com/NHPT/wechatapi/blob/main/doc.md#getcontactbynicknamenickname-string-byte)
        - [GetContacts](https://github.com/NHPT/wechatapi/blob/main/doc.md#getcontactsargs-string-int64-byte)
        - [GetMsg](https://github.com/NHPT/wechatapi/blob/main/doc.md#getmsg-msgdata)
        - [SendMsg](https://github.com/NHPT/wechatapi/blob/main/doc.md#sendmsgmsg-string-tousername-string-error)
        - [JoinChatroom](https://github.com/NHPT/wechatapi/blob/main/doc.md#joinchatroominviteuser-string-chatroom-string-error)
        - [RmChatroom](https://github.com/NHPT/wechatapi/blob/main/doc.md#rmchatroomuser-string-chatroom-string-error)
        - [Logout](https://github.com/NHPT/wechatapi/blob/main/doc.md#logout)
- [错误代码](https://github.com/NHPT/wechatapi/blob/main/doc.md#%E9%94%99%E8%AF%AF%E4%BB%A3%E7%A0%81)

wechatapi是基于Web微信开发的适用于个人微信的API接口，该API接口中定义了一些结构体、全局变量和函数。

## 结构体

wechatapi定义了3个结构体：`WechatData`、`MySelf`和`MsgData`结构体。

`WechatData`结构体用来存放函数需要使用的数据，具体包含以下字段：

- `Skey`
- `Sid`
- `Uin`
- `PassTicket`
- `Synckey`
- `DeviceID`
- `MsgID`
- `Webwx_data_ticket`

`MySelf`结构体用来存放自己的信息，具体包含以下字段：

- `UserName`：用户名
- `NickName`：昵称
- `HeadImgUrl`：头像地址
- `Sex`：性别
- `Signature`：签名

`MsgData`结构体用来存放消息，具体包含以下字段：

- `Srctype`：消息来源
    - 0：表示无消息
    - 1：表示群消息
    - 2：表示公众号推送消息
    - 3：表示联系人消息
    - 1101：表示登录过期。
- `Srcname`：来源名称
    - `Srctype=1`时表示群名称
    - `Srctype=2`时表示公众号名称
    - `Srctype=3`时为空。
- `Msgtype`：消息类型
    - 1：表示文本消息
    - 3：表示图片消息
    - 34：表示语音消息
    - 43：表示视频消息
    - 47：表示表情包消息
    - 49：表示文件消息、视频号消息、转发的公众号消息和其它消息
    - 51：表示空消息
    - 10000：表示拍一拍消息
- `Msg`：消息内容
- `Fromuser`：发送者UserName
- `Usernick`：发送者NickName
- `Createtime`：创建时间，即发送时间

## 全局变量

- `webwx`：`WechatData`结构体的实例。登录成功后存储后续函数使用的参数数据。
- `myself` ：`MySelf`结构体的实例。登录成功后存储自己的信息。
- `chatrooms`：`[]byte`类型，登录成功后存储获取到的所有群聊的名称、用户名和成员数量，即：
    - `NickName`
    - `UserName`
    - `MemberCount`
- `contactlist`：`[]byte`类型，登录成功后存储所有联系人信息，包括用户名、昵称、性别、城市等：
    - `UserName`
    - `NickName`
    - `HeadImgUrl`
    - `RemarkName`
    - `Sex`
    - `Signature`
    - `Province`
    - `City`
    - `SnsFlag`
- `officiallist`：`[]byte`类型，登录成功后存储所有公众号信息，包括用户名、昵称、城市等：
    - `UserName`
    - `NickName`
    - `HeadImgUrl`
    - `Signature`
    - `Province`
    - `City`
    - `SnsFlag`

## 函数

### 无需登录

#### PrintQRCode(info string) error

直接在终端输出二维码，依赖`github.com/skip2/go-qrcode`库。

- 参数`info`为要在终端输出的二维码对应的字符串
- 返回错误信息，如果发生错误则返回错误信息，否则返回`nil`

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
)
func main() {
    wechatapi.PrintQRCode("Hello wechatapi")
}
```

#### Login() error

函数获取二维码UUID，输出登录二维码，并检查是否登录成功，如果登录成功则解析返回数据到`webwx`结构体实例中，以备后续函数使用，否则返回错误信息，如果二维码过期则重新输出新的二维码。函数初始化了`webwx`变量的一部分字段。

- 无参数
- 返回错误信息

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println("Login success.")
}
```

### 需要登录

#### GetMySelf() MySelf

获取自己的信息，包括用户名、昵称、签名等信息。

- 无参数
- 返回`MySelf`结构体

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    myself := wechatapi.GetMySelf()
    if myself.UserName != "" {
        fmt.Println("用户名:",myself.UserName)
        fmt.Println("昵称:",myself.NickName)
        fmt.Println("头像:",myself.HeadImgUrl)
        fmt.Println("性别:",myself.Sex)
        fmt.Println("签名:",myself.Signature)
    }
}
```

#### GetChatrooms() (int64, []byte)

获取所有群聊信息。

- 无参数
- 返回所有群聊数量
- 返回所有群聊信息的json数组，包括群聊名称、群聊UserName和群成员数量

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "github.com/buger/jsonparser"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
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
}
```

#### GetContactlist() (int64, []byte)

获取所有联系人信息。

- 无参数
- 返回所有联系人数量
- 返回所有联系人信息的json数组，包含昵称、用户名、性别、城市等

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    contactcount, contactlist := wechatapi.GetContactlist()
    fmt.Println("联系人数量:",contactcount)
    fmt.Println("联系人信息:",string(contactlist))
}
```

#### GetOfficiallist() (int64, []byte)

获取所有公众号信息。

- 无参数
- 返回所有公众号数量
- 返回所有公众号信息的json数组

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    offc, offclist :=  wechatapi.GetOfficiallist()
    fmt.Println("公众号数量:",offc)
    fmt.Println("公众号信息:",string(offclist))
}
```

#### GetContactByRemarkName(remarkname string) []byte

通过联系人备注获取联系人信息。

- 参数`remarkname` 为联系人备注
- 返回联系人信息的json数据

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
    "github.com/buger/jsonparser"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
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
}
```

#### GetContactByNickName(nickname string) []byte

通过联系人昵称获取联系人信息。

- 参数`nickname`为联系人昵称
- 返回联系人信息的json数据

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    nc := wechatapi.GetContactByNickName("张三")
    if len(nc) > 0 {
        fmt.Println("联系人信息:",string(nc))
    }
}
```

#### GetContacts(args ...string) (int64, []byte)

根据条件获取指定省份、城市、性别的联系人。

- 第一个参数为省份
- 第二个参数为城市
- 第三个参数为性别
- 返回联系人信息的json数组

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
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
}
```

#### GetMsg() MsgData

接收联系人、群和公众号等消息。

- 无参数
- 返回`MsgData`结构体

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "time"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
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
        }
    }
}

```

#### SendMsg(msg string, toUserName string) error

发送消息。

- 参数`msg`为`string`类型消息内容
- 参数`toUserName `为`string`类型接收者UserName
- 返回错误信息

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
    "github.com/buger/jsonparser"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    // 根据备注获取联系人信息
    cb := wechatapi.GetContactByRemarkName("团团")
    if len(cb) > 0 {
        fmt.Println("联系人信息:",string(cb))
        uname, _ := jsonparser.GetString(cb,"UserName")
        // 发送消息
        _ = wechatapi.SendMsg("你好我有一个帽衫",uname)
    }
}
```

#### JoinChatroom(inviteuser string, chatroom string) error

邀请联系人加入群聊。

- 参数`inviteuser `为联系人UserName
- 参数`chatroom `为群聊昵称
- 返回错误信息

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
    "github.com/buger/jsonparser"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    // 根据备注获取联系人信息
    cb := wechatapi.GetContactByRemarkName("团团")
    if len(cb) > 0 {
        fmt.Println("联系人信息:",string(cb))
        uname, _ := jsonparser.GetString(cb,"UserName")
        err := wechatapi.JoinChatroom(uname,"Pentest技术交流")
        if err!= nil {
            fmt.Println("[!] 邀请加入群聊错误")
        }
    }
}

```

#### RmChatroom(user string, chatroom string) error

移除群聊中的某用户。

- 参数`user`为待移除的用户
- 参数`chatroom`为群聊昵称
- 返回错误信息

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
    "github.com/buger/jsonparser"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    // 根据备注获取联系人信息
    cb := wechatapi.GetContactByRemarkName("团团")
    if len(cb) > 0 {
        fmt.Println("联系人信息:",string(cb))
        uname, _ := jsonparser.GetString(cb,"UserName")
        err := wechatapi.RmChatroom(uname,"Pentest技术交流")
        if err!= nil {
            fmt.Println("[!] 移出群聊错误")
        }
    }
}

```

#### Logout()

退出微信。

- 无参数
- 无返回

使用示例：

```Go
package main
import(
    "github.com/NHPT/wechatapi"
    "fmt"
)
func main() {
    err := wechatapi.Login()
    if err != nil {
        fmt.Println(err.Error())
    }
    wechatapi.Logout()
  fmt.Println("Wechat logged out.")
}
```

## 错误代码

|**错误代码**|**说明**|
|-|-|
|E00001|网络错误|
|E00002|读取响应错误|
|E00003|登录二维码过期|
|E00004|响应不存在window.redirect_uri|
|E00005|响应不存在window.code|
|E00006|微信账号未实名认证或未添加支付方式|
|E00007|Json数据解析错误|
|E00008|缺少必须参数|
|E00009|Json序列化错误|
|E00010|消息发送失败|
|E00011|群聊不存在|
