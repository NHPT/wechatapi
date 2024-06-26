# 项目介绍

本项目是基于Web微信开发的个人微信API接口，支持扫码登录，过期重新登录，获取微信联系人信息，获取微信群及成员信息，接收消息，发送消息，邀请加群，移除群聊和退出登录等功能。可基于此API接口实现AI、WebHook、自动回复、智能客服等和个人微信的对接。

## 登录

- 直接通过控制台输出的二维码扫码登录
- 访问输出的链接扫码登录

## 获取联系人

- 可获取所有联系人信息和数量
- 也可以根据省份、城市、性别等条件获取符合条件的联系人信息和数量。
- 可根据备注获取联系人信息和数量
- 可根据昵称获取联系人信息和数量

## 获取群聊

- 可获取所有群聊数量和名称，以及群成员信息

## 获取公众号

- 可获取所有公众号数量和名称

## 接收消息

- 接收所有消息，返回消息来源、消息类型、发送者、消息内容等
- 根据来源筛选消息：群聊、公众号、联系人
- 根据类型筛选消息：文本、图片、表情包等

|**消息类型**|**识别**|**提取**|
|-|-|-|
|文本|√|√|
|图片|√|√|
|语音|√|×|
|视频|√|×|
|表情包|√|√|
|公众号推文|√|√|
|文件|√|×|
|视频号消息|√|×|
|拍一拍|√|√|


## 发送消息

- 发送文本消息

## 邀请入群

- 发送邀请入群的消息

## 移出群聊

- 将指定用户移出指定群聊

# 安装与使用

## 安装

```Bash
go get github.com/NHPT/wechatapi
```

## 使用示例

示例代码见[test.go](https://github.com/NHPT/wechatapi/blob/main/test/test.go)

![](https://github.com/NHPT/wechatapi/blob/main/imgs/run.png)

![](https://github.com/NHPT/wechatapi/blob/main/imgs/info.png)

## 详细文档

请查阅[文档](https://github.com/NHPT/wechatapi/blob/main/doc.md)

# License

[LICENSE](https://github.com/NHPT/wechatapi/blob/main/LICENSE)
