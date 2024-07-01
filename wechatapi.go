package wechatapi

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "html"
    "regexp"
	"time"
    "github.com/skip2/go-qrcode"
    "github.com/buger/jsonparser"
    "strings"
    "encoding/xml"
    "encoding/json"
    "math/rand"
    "strconv"
)

type WechatData struct {
	Skey     string `xml:"skey"`
	Sid    string `xml:"wxsid"`
	Uin    string `xml:"wxuin"`
	PassTicket string `xml:"pass_ticket"`
    Synckey  []byte
    DeviceID string
    MsgID    string
    Webwx_data_ticket string
}

type MySelf struct {
    UserName   string `json:"UserName"`
    NickName   string `json:"NickName"`
    HeadImgUrl string `json:"HeadImgUrl"`
    Sex        int64 `json:"Sex"`
    Signature  string `json:"Signature"`
}

type MsgData struct {
    Srctype int64 `json:"Srctype"`
    Srcname string `json:"Srcname"`
    Msgtype int64 `json:"Msgtype"`
    Msg string `json:"Msg"`
    Fromuser string `json:"Fromuser"`
    Usernick string `json:"Usernick"`
    Createtime int64 `json:"Createtime"`
}

var (
    webwx WechatData
    myself MySelf
    chatrooms []byte
    contactlist []byte
	officiallist []byte
    wxurl = "https://wx.qq.com"
    webpushurl = "https://webpush.wx.qq.com"
    extspam = "Go8FCIkFEokFCggwMDAwMDAwMRAGGvAESySibk50w5Wb3uTl2c2h64jVVrV7"+
    "gNs06GFlWplHQbY/5FfiO++1yH4ykCyNPWKXmco+wfQzK5R98D3so7rJ5LmGFvBLjGceleySrc3SOf2Pc1gVeh"+
    "zJgODeS0lDL3/I/0S2SSE98YgKleq6Uqx6ndTy9yaL9qFxJL7eiA/R3SEfTaW1SBoSITIu+EEkXff+Pv8NHOk7"+
    "N57rcGk1w0ZzRrQDkXTOXFN2iHYIzAAZPIOY45Lsh+A4slpgnDiaOvRtlQYCt97nmPLuTipOJ8Qc5pM7ZsOsAP"+
    "PrCQL7nK0I7aPrFDF0q4ziUUKettzW8MrAaiVfmbD1/VkmLNVqqZVvBCtRblXb5FHmtS8FxnqCzYP4WFvz3T0T"+
    "crOqwLX1M/DQvcHaGGw0B0y4bZMs7lVScGBFxMj3vbFi2SRKbKhaitxHfYHAOAa0X7/MSS0RNAjdwoyGHeOepX"+
    "OKY+h3iHeqCvgOH6LOifdHf/1aaZNwSkGotYnYScW8Yx63LnSwba7+hESrtPa/huRmB9KWvMCKbDThL/nne14h"+
    "nL277EDCSocPu3rOSYjuB9gKSOdVmWsj9Dxb/iZIe+S6AiG29Esm+/eUacSba0k8wn5HhHg9d4tIcixrxveflc"+
    "8vi2/wNQGVFNsGO6tB5WF0xf/plngOvQ1/ivGV/C1Qpdhzznh0ExAVJ6dwzNg7qIEBaw+BzTJTUuRcPk92Sn6Q"+
    "Dn2Pu3mpONaEumacjW4w6ipPnPw+g2TfywJjeEcpSZaP4Q3YV5HG8D6UjWA4GSkBKculWpdCMadx0usMomsSS/"+
    "74QgpYqcPkmamB4nVv1JxczYITIqItIKjD35IGKAUwAA=="
)

/** 获取自己的信息
 * @return MySelf 返回自己的信息
 */
func GetMySelf() MySelf {
    return myself
}

/** 获取所有群聊信息
 * @return int64 群聊数量
 * @return []string 返回所有群聊
 */
func GetChatrooms() (int64, []byte) {
    var count int64 = 0
    jsonparser.ArrayEach(chatrooms, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }
        count++
    })
    return count, chatrooms
}

/** 根据群昵称获取群UserName
 * @param  string 群昵称
 * @return string 群UserName
 */
func getChatroomName(nickname string) string {
    name := ""
    jsonparser.ArrayEach(chatrooms, func(value []byte, dataType jsonparser.ValueType, offset int, err error){
        if err != nil {
            return
        }
        nick, _ := jsonparser.GetString(value,"NickName")
        if nickname == nick {
            name, _ = jsonparser.GetString(value,"UserName")
        }
    })
    return name
}

/** 获取所有联系人信息
 * @return int64 联系人数量
 * @return []byte 返回所有联系人
 */
func GetContactlist() (int64, []byte) {
    var count int64 = 0
    jsonparser.ArrayEach(contactlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }
        count++
    })
    return count, contactlist
}

/** 获取所有公众号信息
 * @return int64 公众号数量
 * @return []byte 返回所有公众号
 */
func GetOfficiallist() (int64, []byte) {
    var count int64 = 0
    jsonparser.ArrayEach(officiallist, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }
        count++
    })
    return count, officiallist
}

/**
 * 直接在终端输出二维码
 * 使用github.com/skip2/go-qrcode库
 * @param info 要显示的字符串
 * @return error 错误信息
 */
func PrintQRCode(info string) error {
	// 生成QR码
	qr, err := qrcode.New(info, qrcode.Highest)
	if err != nil {
		return err
	}
	// 获取 QR 码图像
	size := 5 
	img := qr.Image(size)

	// 遍历图像的每个像素点
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 获取像素的颜色值
			r, g, b, _ := img.At(x, y).RGBA()
			_, _, _, alpha := img.At(x, y).RGBA()
			// 假设黑色像素为 QR 码的黑色部分，白色或透明像素为白色部分
			if alpha == 0xFFFF && (r == 0 && g == 0 && b == 0) {
				fmt.Print("██")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
    return nil
}

/**
 * 获取微信登录二维码UUID值，用于后续登录
 * @return string 返回UUID
 */
func getQRCodeuuid() string {
    // 请求微信登录页面获取响应中的值
    resp, err := http.Get(wxurl + "/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage%3Fmod%3Ddesktop&fun=new")
    if err != nil {
        return ""
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return ""
    }
    re := regexp.MustCompile(`window\.QRLogin\.uuid\s*=\s*"([^"]+)"`)
    uuid := re.FindStringSubmatch(string(body))
    if uuid == nil {
		return ""
    }
    return uuid[1]
}

/**
 * 检查是否登录，如果登录则返回 redirect_uri，否则返回错误信息
 * @param uuid 登录二维码UUID
 * @return string 返回 redirect_uri
 * @return error 错误信息
 */
func checkLogin(uuid string) (string, error) {
    p := time.Now().Unix()
    r := int(p * 11)
    loginURL := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=%s&tip=0&r=%d&_=%d", uuid, r, p)
    resp, err := http.Get(loginURL)
    if err != nil {
        return "", fmt.Errorf("E00001: Please check your network connection.")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("E00002: Failed to read response body.")
    }
    codeRegex := regexp.MustCompile(`window\.code=(\d+);`)
    redirectURIRegex := regexp.MustCompile(`window\.redirect_uri\s*=\s*"([^"]+)";`)
    codeMatch := codeRegex.FindStringSubmatch(string(body))
    if codeMatch != nil {
        // 二维码过期
        if codeMatch[1] == "400" {
            return "400", fmt.Errorf("E00003: The QR code has expired.")
        }
        // 登录成功
        if codeMatch[1] == "200" {
            redirectURIMatch := redirectURIRegex.FindStringSubmatch(string(body))
            if redirectURIMatch == nil {
                return "", fmt.Errorf("E00004: window.redirect_uri not found.")
            }
            // 返回提取到的 redirect_uri
            return redirectURIMatch[1], nil
        }
        return "", fmt.Errorf("")
    }else {
        return "", fmt.Errorf("E00005: window.code is not exist")
    }
}

/**
 * 输出登录二维码，并检查登录结果，登录成功后初始化数据
 * @return error 错误信息
 */
func Login() error {
	uuid := getQRCodeuuid()
	if uuid == "" {
        return fmt.Errorf("[!] Cound not find uuid.")
	} else {
		fmt.Println("[*] 请扫描二维码登录微信")
        fmt.Println("[*] https://api.pwmqr.com/qrcode/create/?url=" + "https://login.weixin.qq.com/l/" + uuid)
        fmt.Println("[*] https://api.isoyu.com/qr/?m=1&e=L&p=20&url=" + "https://login.weixin.qq.com/l/" + uuid)
        PrintQRCode("https://login.weixin.qq.com/l/" + uuid)
        // 轮询检查是否已经扫码登录
		for {
			redirectURI, err := checkLogin(uuid)
			if err != nil {
                if redirectURI == "400" {
                    uuid = getQRCodeuuid()
                    fmt.Println("[*] 二维码过期，请重新扫描二维码登录微信")
                    fmt.Println("[*] https://api.pwmqr.com/qrcode/create/?url=" + "https://login.weixin.qq.com/l/" + uuid)
                    fmt.Println("[*] https://api.isoyu.com/qr/?m=1&e=L&p=20&url=" + "https://login.weixin.qq.com/l/" + uuid)
                    PrintQRCode("https://login.weixin.qq.com/l/" + uuid)
                    continue
                }
                if strings.Contains(err.Error(), "E00001") {
                    fmt.Println("[!] " + err.Error())
                    return err
                }
			} else {
                fmt.Println("[+] 登录成功!")
                client := &http.Client{
                    CheckRedirect: func(req *http.Request, via []*http.Request) error {
                        // 返回错误以阻止自动跟随重定向
                        return http.ErrUseLastResponse
                    },
                }
                req, _ := http.NewRequest("GET",redirectURI + "&fun=new&version=v2&mod=desktop",nil)
                //fmt.Println(redirectURI + "&fun=new&version=v2&mod=desktop")
                req.Header.Add("Referer","https://wx.qq.com/?target=t")
                req.Header.Add("client-version","2.0.0")
                req.Header.Add("extspam",extspam)
                resp, err := client.Do(req)
                if err != nil {
                    fmt.Println("[!] E00001: Please check your network connection.")
                    return err
                }
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                    fmt.Println("[!] E00002: Failed to read response body.")
                }
                cookies := resp.Header["Set-Cookie"]
                webwx_data_ticket := ""
                for _, cookie := range cookies {
                    if strings.HasPrefix(cookie, "webwx_data_ticket="){
                        parts := strings.Split(cookie, ";")
                        for _, part := range parts {
                            part = strings.TrimSpace(part)
                            if strings.HasPrefix(part, "webwx_data_ticket="){
                                webwx_data_ticket = part[len("webwx_data_ticket="):]
                                break
                            }
                        }
                        if webwx_data_ticket != "" {
                            webwx.Webwx_data_ticket = webwx_data_ticket
                            break
                        }
                    }
                }
                err = xml.Unmarshal([]byte(body), &webwx)
                if err != nil {
                    fmt.Println("[!] E00006: You need real name authentication and WeChat payment activation.")
                }
                fmt.Println("[*] 正在初始化数据...")
                webInit()
                fmt.Println("[*] 初始化数据成功!")
                return nil
			}
			time.Sleep(2 * time.Second)
		}
	}
}

/**
 * 获取webInit数据，初始化数据
 */
func webInit() {
    r := time.Now().Unix()
    rand.Seed(time.Now().UnixNano())
    const count = 15
	var randomString string
	for i := 0; i < count; i++ {
		randomDigit := rand.Intn(10)
		randomString += fmt.Sprint(randomDigit)
	}
	webwx.DeviceID = "e" + randomString
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxinit?r=%d&pass_ticket=%s", r, webwx.PassTicket)
    data := fmt.Sprintf(`{"BaseRequest":{"Uin":"%s","Sid":"%s","Skey":"%s","DeviceID":"%s"}}`, webwx.Uin, webwx.Sid, webwx.Skey, webwx.DeviceID)
    resp, err := http.Post(url, "application/json;charset=utf-8", strings.NewReader(data))
    if err != nil {
        fmt.Println("[!] E00001: Please check your network connection.")
        return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("[!] E00002: Failed to read response body.")
        return
    }
    SyncKey, _, _, err := jsonparser.Get(body,"SyncKey")
    if err != nil {
        fmt.Println("[!] E00007: %s",err)
        return
    }
    webwx.Synckey = SyncKey
    err = setMySelf(body)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("[*] 正在获取联系人，请稍后...")
    err = getAllContact()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("[*] 获取联系人成功，正在获取微信群，请稍后...")
    _ = getStatusNotify()
    getChatRoomFromWebwxsync()
    fmt.Println("[*] 获取微信群成功")
}

/**
 * 获取自己的信息，并存入myself结构体实例
 * @param jsondata []byte 服务器返回的JSON数据
 * @return error 错误信息
 */
func setMySelf(jsondata []byte) error {
    username, err := jsonparser.GetString(jsondata, "User", "UserName")
    if err != nil {
        return fmt.Errorf("E00007: %s",err)
    }
    myself.UserName = username
    nickname, err := jsonparser.GetString(jsondata, "User", "NickName")
    if err != nil {
        return fmt.Errorf("E00007: %s",err)
    }
    myself.NickName = nickname
    headimgurl, err := jsonparser.GetString(jsondata, "User", "HeadImgUrl")
    if err != nil {
        return fmt.Errorf("E00007: %s",err)
    }
    myself.HeadImgUrl = headimgurl
    if value, err := jsonparser.GetInt(jsondata, "User", "Sex"); err == nil {
        myself.Sex = value
    }
    signature, err := jsonparser.GetString(jsondata, "User", "Signature")
    if err != nil {
        return fmt.Errorf("E00007: %s",err)
    }
    myself.Signature = signature
    return nil
}

/**
 * 获取微信状态通知
 * @return error 错误信息
 */
func getStatusNotify() error {
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxstatusnotify?pass_ticket=%s",webwx.PassTicket)
    data := fmt.Sprintf(`{"BaseRequest":{"Uin":%s,"Sid":"%s","Skey":"%s","DeviceID":"%s"},"Code":3,"FromUserName":"%s","ToUserName":"%s","ClientMsgId":"%d"}`,webwx.Uin,webwx.Sid,webwx.Skey,webwx.DeviceID,myself.UserName,myself.UserName,time.Now().Unix())
    resp, err := http.Post(url, "application/json;charset=utf-8", strings.NewReader(data))
    if err != nil {
        return fmt.Errorf("E00001: Please check your network connection.")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("E00002: Failed to read response body.")
    }
    msgID, err := jsonparser.GetString(body,"MsgID")
    if err != nil {
        return fmt.Errorf("E00007: %s",err)
    }
    webwx.MsgID = msgID
    return nil
}

/**
 * 获取SyncKey
 * @param SyncKey []byte Json数据
 * @return string 返回格式化后的synckey数据
 * @return error 错误信息
 */
func getSyncKey(SyncKey []byte) (string, error) {
    SyncKeyList, _, _, err := jsonparser.Get(SyncKey,"List")
    if err != nil {
        return "", fmt.Errorf("E00007: %s",err)
    }
    // Initialize an empty string for the result
    result := ""

    // Iterate over the array elements
    jsonparser.ArrayEach(SyncKeyList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }

        // Use jsonparser to get the "Key" and "Val" values
        key, _ := jsonparser.GetInt(value, "Key")
        val, _ := jsonparser.GetInt(value, "Val")

        // Format the current element as "key_val" and add it to the result string
        if len(result) > 0 {
            result += "%7C" // Add separator (URL-encoded pipe symbol)
        }
        result += fmt.Sprintf("%d_%d", key, val)
    })
    return result, nil
}

/** 通过联系人备注名获取联系人信息
 * @param remarkname string 联系人备注名
 * @return []byte 返回联系人信息
 */
func GetContactByRemarkName(remarkname string) []byte {
    var contact []byte
    jsonparser.ArrayEach(contactlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }
        remarkName, _ := jsonparser.GetString(value, "RemarkName")
        if remarkName == remarkname {
            contact = value
        }
    })
    return contact
}

/** 通过联系人昵称获取联系人信息
 * @param nickname string 联系人昵称
 * @return []byte 返回联系人信息
 */
func GetContactByNickName(nickname string) []byte {
    var contact []byte
    jsonparser.ArrayEach(contactlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }
        nickName, _ := jsonparser.GetString(value, "NickName")
        if nickName == nickname {
            contact = value
        }
    })
    return contact
}

/** 获取指定省份、城市、性别的联系人
 * @param args ...string 省份、城市、性别
 * @return int64 符合条件的联系人数量
 * @return []byte 符合条件的联系人信息
 */
func GetContacts(args ...string) (int64, []byte) {
    if len(args) != 3 {
        fmt.Println("[!] E00008: 3 parameters are required.")
        return 0,[]byte{}
    }
    var count int64
    var contactarray []map[string]interface{}
    var contactslist []byte
	jsonparser.ArrayEach(contactlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        contact := make(map[string]interface{})
        _ = json.Unmarshal(value, &contact)
        if args[0] != "" {
            Province, _ := jsonparser.GetString(value, "Province")
            if args[0] == Province {
                contactarray = append(contactarray,contact)
                count ++
            }
        }
        if args[1] != "" {
            City, _ := jsonparser.GetString(value, "City")
            if args[1] == City {
                contactarray = append(contactarray,contact)
                count ++
            }
        }
        if args[2] != "" {
            Sexint, _ := strconv.ParseInt(args[2], 10, 64)
            Sex, _ := jsonparser.GetInt(value, "Sex")
            if Sexint == Sex {
                contactarray = append(contactarray,contact)
                count ++
            }
        }
        
    })
    contactslist,_ = json.Marshal(contactarray)
    return count, contactslist
}
/**
 * 获取所有联系人，包括公众号
 * @return error 错误信息
 */
func getAllContact() error {
    var Seq, MemberCount int64 = 0,0
    var memberArray []map[string]interface{}
    for {
        url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxgetcontact?pass_ticket=%s&r=%d&seq=%d&skey=%s",webwx.PassTicket,time.Now().Unix(),Seq,webwx.Skey)
        req, err := http.NewRequest("GET",url,nil)
        if err != nil {
            return err
        }
        req.Header.Add("Cookie","wxuin="+webwx.Uin+";wxsid="+webwx.Sid)
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            return fmt.Errorf("E00001: Please check your network connection.")
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return fmt.Errorf("E00002: Failed to read response body.")
        }
        if value,err := jsonparser.GetInt(body, "Seq"); err == nil {
            Seq = value
        }
        if value,err := jsonparser.GetInt(body, "MemberCount"); err == nil {
            MemberCount += value
        }
        memberLists, _, _, err := jsonparser.Get(body,"MemberList")
        if err != nil {
            return fmt.Errorf("E00007: %s",err)
        }
        var temp []map[string]interface{}
        err = json.Unmarshal(memberLists, &temp)
        if err != nil {
            return err
        }
        memberArray = append(memberArray, temp...)
        if Seq == 0 {
            break
        }
    }
    memberList, err := json.Marshal(memberArray)
	if err != nil {
		return fmt.Errorf("E00009: error marshalling merged JSON array: %v", err)
	}
    // 保留指定字段
	var contactarray []map[string]interface{}
    var officialarray []map[string]interface{}
	jsonparser.ArrayEach(memberList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		contact := make(map[string] interface{})
		username, _ := jsonparser.GetString(value, "UserName")
		nickname, _ := jsonparser.GetString(value, "NickName")
		headimgurl, _ := jsonparser.GetString(value, "HeadImgUrl")
		remarkname, _ := jsonparser.GetString(value, "RemarkName")
		sex, _ := jsonparser.GetInt(value, "Sex")
		signature, _ := jsonparser.GetString(value, "Signature")
		province, _ := jsonparser.GetString(value, "Province")
		city, _ := jsonparser.GetString(value, "City")
		snsflag, _ := jsonparser.GetInt(value, "SnsFlag")
		contact["UserName"] = username
		contact["NickName"] = nickname
		contact["HeadImgUrl"] = headimgurl
		contact["RemarkName"] = remarkname
		contact["Sex"] = sex
		contact["Signature"] = signature
		contact["Province"] = province
		contact["City"] = city
		contact["SnsFlag"] = snsflag
		if snsflag == 0 {
            officialarray = append(officialarray,contact)
			officiallist,_ = json.Marshal(officialarray)
		} else {
            contactarray = append(contactarray,contact)
			contactlist,_ = json.Marshal(contactarray)
		}
	})
    return nil
}


/**
 * 批量获取联系人
 * @param username ...string 用户名
 */
func getBatchContact(username ...string) []byte {
    count := int64(len(username))
    users := make([]map[string]string, count)
	for i, user := range username {
		users[i] = map[string]string{
			"UserName":        user,
			"EncryChatRoomId": "",
		}
	}
	userList, err := json.Marshal(users)
	if err != nil {
        return []byte{}
	}

	r := time.Now().Unix()
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?type=ex&r=%d&pass_ticket=%s",r, webwx.PassTicket)
    data := fmt.Sprintf(
        `{
            "BaseRequest":{
                "Uin":%s,
                "Sid":"%s",
                "Skey":"%s",
                "DeviceID":"%s"
            },
            "Count":%d,
            "List":%s
        }`,
        webwx.Uin,webwx.Sid,webwx.Skey,webwx.DeviceID,count,string(userList))
    req, err := http.NewRequest("POST",url,strings.NewReader(data))
    if err != nil {
        return []byte{}
    }
    req.Header.Add("Cookie","wxuin="+webwx.Uin+";wxsid="+webwx.Sid)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return []byte{}
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []byte{}
    }
    contactList, _, _, err := jsonparser.Get(body,"ContactList")
    if err != nil {
        return []byte{}
    }
    return contactList
}


/**
 * 检查消息
 * @return int64 返回retcode的值
 * @return int64 返回selector的值，为0表示无消息
 */
func syncCheck() (int64, int64) {
    r := time.Now().Unix()
    synckey, _ := getSyncKey(webwx.Synckey)
    url := fmt.Sprintf(webpushurl + "/cgi-bin/mmwebwx-bin/synccheck?r=%d&skey=%s&sid=%s&uin=%s&deviceid=%s&synckey=%s&_=%d",r,webwx.Skey,webwx.Sid,webwx.Uin,webwx.DeviceID,synckey,r)
    req, err := http.NewRequest("GET",url,nil)
    if err != nil {
        return 0, 0
    }
    req.Header.Add("Cookie","wxuin="+webwx.Uin+";webwx_data_ticket="+webwx.Webwx_data_ticket)
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return 0, 0
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    re1 := regexp.MustCompile(`retcode:"([^"]+)"`)
    re2 := regexp.MustCompile(`selector:"([^"]+)"`)
    retcode := re1.FindStringSubmatch(string(body))
    if len(retcode) < 2 {
        return 0, 0
    }
    retcodeInt, _ := strconv.ParseInt(retcode[1], 10, 64)
    selector := re2.FindStringSubmatch(string(body))
    if len(selector) < 2 {
        return 0, 0
    }
    selectorInt, _ := strconv.ParseInt(selector[1], 10, 64)
    //fmt.Println(retcodeInt, selectorInt)
    return retcodeInt, selectorInt
}

/** 检查用户名是否在联系人列表中
 *  @param slice []byte 联系人列表
 *  @param str string 要检查的用户名
 *  @return bool 如果用户名在列表中，返回 true；否则返回 false
 */
func contains(slice []byte, str string) bool {
    res := false
    jsonparser.ArrayEach(slice, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }
        username, _ := jsonparser.GetString(value, "UserName")
        if username == str {
            res = true
        }
    })
	return res
}

/** 获取Webwxsync数据
 * @return []byte webwxsync的响应body
 * @return error 错误信息
 */
func getWebwxsync() ([]byte, error) {
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxsync?sid=%s&skey=%s&pass_ticket=%s",webwx.Sid, webwx.Skey, webwx.PassTicket)
    data := fmt.Sprintf(
        `{
            "BaseRequest":{
                "Uin":%s,
                "Sid":"%s",
                "Skey":"%s",
                "DeviceID":"%s"
            },
            "SyncKey":%s,
            "rr":%s
        }`,
        webwx.Uin,webwx.Sid,webwx.Skey,webwx.DeviceID,webwx.Synckey,webwx.Uin)
    req, err := http.NewRequest("POST",url,strings.NewReader(data))
    if err != nil {
        return []byte{}, err
    }
    req.Header.Add("Cookie","wxuin="+webwx.Uin+";wxsid="+webwx.Sid)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return []byte{}, fmt.Errorf("E00001: Please check your network connection.")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []byte{}, fmt.Errorf("E00002: Failed to read response body.")
    }
    // 更新Synckey
    SyncKey, _, _, err := jsonparser.Get(body,"SyncKey")
    if err != nil {
        return []byte{}, fmt.Errorf("E00007: %s",err)
    }
    webwx.Synckey = SyncKey
    return body, nil
}

/** 获取Webwxsync中的群聊
 */
func getChatRoomFromWebwxsync() {
    StatusNotifyUserName := ""
    for i := 0; i < 20; i++ {
        webwxsync, err := getWebwxsync()
        if err != nil {
            continue
        }
        AddMsgList, _, _, err := jsonparser.Get(webwxsync, "AddMsgList")
        if err != nil {
            continue
        }
        // 遍历 AddMsgList 数组中的每个元素
        jsonparser.ArrayEach(AddMsgList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
            if err != nil {
                return
            }
            StatusNotifyUserName, _ = jsonparser.GetString(value, "StatusNotifyUserName")
        })
        if StatusNotifyUserName != "" {
            break
        }
    }
    if StatusNotifyUserName == "" {
        fmt.Println("[!] 未获取到群聊")
        return
    }
    var chatroomarray []map[string]interface{}
    parts := strings.Split(StatusNotifyUserName, ",")   
    for _, part := range parts {
        trimmedPart := strings.TrimSpace(part)
        if strings.HasPrefix(trimmedPart, "@@") && !contains(chatrooms, trimmedPart) {
            ContactList := getBatchContact(trimmedPart)
            tempnick := ""
            var tempcount int64
            jsonparser.ArrayEach(ContactList, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
                if err != nil {
                    return
                }
                tempnick,_ = jsonparser.GetString(value,"NickName")
                tempcount,_ = jsonparser.GetInt(value,"MemberCount")
            })
            data := make(map[string] interface{})
            data["UserName"] = trimmedPart
            data["NickName"] = tempnick
            data["MemberCount"] = tempcount
            chatroomarray = append(chatroomarray, data)
        }
    }
    chatrooms,_ = json.Marshal(chatroomarray)
}

/**
 * 获取消息
 */
func GetMsg() MsgData {
    var createtime,srctype,msgtype int64
    fromuser := ""
    srcname := ""
    srcuser := ""
    msg := ""

    for {
        retcode, selector := syncCheck()
        if retcode != 0 {
            return MsgData{Srctype: retcode,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
        }
        if selector == 0 {
            continue
        } else {
            break
        }
    }
    body, err := getWebwxsync()
    if err != nil {
        fmt.Println("[!]",err)
        return MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
    }
    // 获取消息
    count, _ := jsonparser.GetInt(body,"AddMsgCount")
    if count == 0 {
        return MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
    }
    msglist, _, _, err := jsonparser.Get(body,"AddMsgList")
    jsonparser.ArrayEach(msglist, func(value []byte, dataType jsonparser.ValueType, offset int, err error){
        if err != nil {
            return
        }
        fromuser, _ = jsonparser.GetString(value,"FromUserName")
        if fromuser == myself.UserName {
            return
        }
        createtime, _ = jsonparser.GetInt(value,"CreateTime")
        msgtype,_ = jsonparser.GetInt(value,"MsgType")
        if msgtype == 51 {
            return
        }
        content,_ := jsonparser.GetString(value,"Content")
        // 群聊消息
        if strings.HasPrefix(fromuser, "@@") {
            srctype = 1
            user := strings.Split(content, ":")
            if len(user) < 2 {
                srctype = 0
                return
            } 
            msg = user[1]
            chatroominfo := getBatchContact(fromuser)
            jsonparser.ArrayEach(chatroominfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
                if err != nil {
                    return
                }
                srcname,_ = jsonparser.GetString(value,"NickName")
                memberlist,_ ,_ ,err := jsonparser.Get(value,"MemberList")
                if err != nil {
                    fmt.Println(err)
                }
                jsonparser.ArrayEach(memberlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
                    if err != nil {
                        return
                    }
                    membername,_ := jsonparser.GetString(value,"UserName")
                    if membername == user[0] {
                        srcuser,_ = jsonparser.GetString(value,"NickName")
                    }
                })
            })
            switch msgtype {
                case 1:
                    //fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的文本消息:",msg)
                    return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                case 3:
                    msgid,_ := jsonparser.GetString(value,"MsgId")
                    msg = wxurl + "/cgi-bin/mmwebwx-bin/webwxgetmsgimg?&MsgID=" + msgid + "&skey=" + webwx.Skey
                    //fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的图片消息:",msg)
                    return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                //case 34:
                //    fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的语音消息暂不支持获取")
                //case 43:
                //    fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的视频消息暂不支持获取")
                case 47:
                    re := regexp.MustCompile(`(?i)cdnurl\s*=\s*"(.*?)"`)
                    match := re.FindStringSubmatch(content)
                    if len(match) > 1 {
                        msg = html.UnescapeString(html.UnescapeString(match[1]))
                        //fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的表情包消息:",msg)
                        return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                    }
                    srctype = 0
                case 49:
                    appmsgtype,_ := jsonparser.GetInt(value,"AppMsgType")
                    switch appmsgtype {
                        case 5:
                            title,_ := jsonparser.GetString(value,"FileName")
                            msg,_ = jsonparser.GetString(value,"Url")
                            msg = title + html.UnescapeString(msg)
                            //fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的公众号消息:",title + msg)
                            return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                        //case 6:
                        //    fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的文件消息，暂不支持获取")
                        //case 51:
                        //    fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的视频号消息，暂不支持获取")
                        default:
                            //fmt.Println("[+]",msgtime,"群 [",srcname,"] 成员 [",srcuser,"] 的其它类型消息，暂不支持获取")
                            return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                        }
                case 10000:
                    return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                default:
                    //fmt.Println("[+] 收到其他消息:",string(content))
                    return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
            }
        } else {
            srcuser = getNicknameByUsername(fromuser,"contactlist")
            // 公众号消息
            if srcuser == "" {
                srctype = 2
                srcname = getNicknameByUsername(fromuser,"officiallist")
                title,_ := jsonparser.GetString(value,"FileName")
                msg_url,_ := jsonparser.GetString(value,"Url")
                msg_url = html.UnescapeString(msg_url)
                msg = title + msg_url
                //fmt.Println("[+]",msgtime,"公众号 [",srcname,"] 的消息:",msg)
                return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
            } else {// 联系人消息
                srctype = 3
                remarkname := getRemarkNameByUsername(fromuser)
                if remarkname != "" {
                    srcuser = remarkname
                }
                //content,_ := jsonparser.GetString(value,"Content")
                //msgtype,_ := jsonparser.GetInt(value,"MsgType")
                switch msgtype {
                    case 1:
                        msg = content
                        //fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的文本消息:",content)
                        return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                    case 3:
                        msgid,_ := jsonparser.GetString(value,"MsgId")
                        msg = wxurl + "/cgi-bin/mmwebwx-bin/webwxgetmsgimg?&MsgID=" + msgid + "&skey=" + webwx.Skey
                        //fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的图片消息:",msg)
                        return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                    //case 34:
                    //    fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的语音消息暂不支持获取")
                    //case 43:
                    //    fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的视频消息暂不支持获取")
                    case 47:
                        re := regexp.MustCompile(`(?i)cdnurl\s*=\s*"(.*?)"`)
                        match := re.FindStringSubmatch(content)
                        if len(match) > 1 {
                            msg = html.UnescapeString(html.UnescapeString(match[1]))
                            //fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的表情包消息:",msg)
                            return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                        }
                        srctype = 0
                    case 49:
                        appmsgtype,_ := jsonparser.GetInt(value,"AppMsgType")
                        switch appmsgtype {
                            case 5:
                                title,_ := jsonparser.GetString(value,"FileName")
                                msg,_ = jsonparser.GetString(value,"Url")
                                msg = title + html.UnescapeString(msg)
                                //fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的公众号消息:",title + msg)
                                return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: title + msg,Createtime: createtime}
                            //case 6:
                            //    fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的文件消息，暂不支持获取")
                            //case 51:
                            //    fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的视频号消息，暂不支持获取")
                            default:
                                //fmt.Println("[+]",msgtime,"[ ",srcuser,"] 的其它类型消息，暂不支持获取")
                                return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                            }
                    case 10000:
                        return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                    default:
                        //fmt.Println("[+] 收到其他消息:",string(content),msgtype)
                        return //MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
                }
            }
        }
    })
    return MsgData{Srctype: srctype,Srcname: srcname,Msgtype: msgtype,Fromuser: fromuser,Usernick: srcuser,Msg: msg,Createtime: createtime}
}

/** 根据用户名获取备注
 * @param username string 用户名
 * @return string 备注
 */
func getRemarkNameByUsername(username string) string {
    remarkname := ""
    jsonparser.ArrayEach(contactlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error){
        if err != nil {
            return
        }
        uName, _ := jsonparser.GetString(value,"UserName")
        if username == uName {
            remarkname, _ = jsonparser.GetString(value,"RemarkName")
        }
    })
    return remarkname
}

/** 根据用户名获取昵称
 * @param username string 用户名
 * @param source string 数据来源，contactlist或officiallist结构体实例
 * @return string 昵称
 */
func getNicknameByUsername(username string, source string) string {
    nickname := ""
    if source == "contactlist" {
        jsonparser.ArrayEach(contactlist, func(value []byte, dataType jsonparser.ValueType, offset int, err error){
            if err != nil {
                return
            }
            uName, _ := jsonparser.GetString(value,"UserName")
            if username == uName {
                nick, _ := jsonparser.GetString(value,"NickName")
                nickname = nick
            }
        })
    }
    if source == "officiallist" {
        jsonparser.ArrayEach(officiallist, func(value []byte, dataType jsonparser.ValueType, offset int, err error){
            if err != nil {
                return
            }
            uName, _ := jsonparser.GetString(value,"UserName")
            if username == uName {
                nick, _ := jsonparser.GetString(value,"NickName")
                nickname = nick
            }
        })
    }
    return nickname
}

/**
 * 发送消息
 * @param msg 消息内容
 * @param toUserName 接收者
 * @return error 错误信息
 */
func SendMsg(msg string, toUserName string) error {
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxsendmsg?pass_ticket=%s",webwx.PassTicket)
    localid := time.Now().UnixNano()
    data := fmt.Sprintf(
        `{
            "BaseRequest":{
                "Uin":%s,
                "Sid":"%s",
                "Skey":"%s",
                "DeviceID":"%s"
            },
            "Msg":{
                "Type":1,
                "Content":"%s",
                "FromUserName":"%s",
                "ToUserName":"%s",
                "LocalID":"%s",
                "ClientMsgId":"%s"
            },
            "Scene":0
        }`,
        webwx.Uin,webwx.Sid,webwx.Skey,webwx.DeviceID,msg,myself.UserName,toUserName,localid,localid)
    req, err := http.NewRequest("POST",url,strings.NewReader(data))
    if err != nil {
        return err
    }
    req.Header.Add("Cookie","wxuin="+webwx.Uin+";wxsid="+webwx.Sid)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("E00001: Please check your network connection.")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("E00002: Failed to read response body.")
    }
    ret,_ := jsonparser.GetInt(body,"BaseResponse","Ret")
    if ret != 0 {
        return fmt.Errorf("E00010: Msg send failed.")
    }
    return nil
}

/**
 * 邀请联系人加入群聊
 * @param inviteuser 被邀请联系人UserName
 * @param chatroom 群聊名称
 * @return error
 */
func JoinChatroom(inviteuser string, chatroom string) error {
    chatroomname := getChatroomName(chatroom)
    if chatroomname == "" {
        return fmt.Errorf("E00011: Chatroom does not exist.")
    }
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=invitemember&pass_ticket=%s",webwx.PassTicket)
    data := fmt.Sprintf(
    `{
        "InviteMemberList":"%s",
        "ChatRoomName":"%s",
        "BaseRequest":{
            "Uin":%s,
            "Sid":"%s",
            "Skey":"%s",
            "DeviceID":"%s"
        }
    }`,
    inviteuser,chatroomname,webwx.Uin,webwx.Sid,webwx.Skey,webwx.DeviceID)
    req, err := http.NewRequest("POST",url,strings.NewReader(data))
    if err != nil {
        return err
    }
    req.Header.Add("Cookie","wxuin="+webwx.Uin+";wxsid="+webwx.Sid)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("E00001: Please check your network connection.")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("E00002: Failed to read response body.")
    }
    ret,_ := jsonparser.GetInt(body,"BaseResponse","Ret")
    if ret != 0 {
        return fmt.Errorf("E00010: Msg send failed.")
    }
    return nil
}

/**
 * 移除群聊中的某用户
 * @param user 待移除用户UserName
 * @param chatroom  群聊
 * @return error 错误信息
 */
func RmChatroom(user string, chatroom string) error {
    chatroomname := getChatroomName(chatroom)
    if chatroomname == "" {
        return fmt.Errorf("E00011: Chatroom does not exist.")
    }
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxupdatechatroom?fun=delmember&pass_ticket=%s",webwx.PassTicket)
    data := fmt.Sprintf(
    `{
        "DelMemberList":"%s",
        "ChatRoomName":"%s",
        "BaseRequest":{
            "Uin":%s,
            "Sid":"%s",
            "Skey":"%s",
            "DeviceID":"%s"
        }
    }`,
    user,chatroomname,webwx.Uin,webwx.Sid,webwx.Skey,webwx.DeviceID)
    req, err := http.NewRequest("POST",url,strings.NewReader(data))
    if err != nil {
        return err
    }
    req.Header.Add("Cookie","wxuin="+webwx.Uin+";wxsid="+webwx.Sid)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("E00001: Please check your network connection.")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("E00002: Failed to read response body.")
    }
    ret,_ := jsonparser.GetInt(body,"BaseResponse","Ret")
    if ret != 1 {
        return fmt.Errorf("E00010: Msg send failed.")
    }
    return nil
}

/**
 * 退出微信
 */
func Logout() {
    url := fmt.Sprintf(wxurl + "/cgi-bin/mmwebwx-bin/webwxlogout?redirect=1&type=0&skey=%s",webwx.Skey)
    data := fmt.Sprintf(`sid=%s&uin=%s`,webwx.Sid,webwx.Uin)
    req, err := http.NewRequest("POST",url,strings.NewReader(data))
    if err != nil {
        fmt.Println(err.Error())
    }
    req.Header.Add("Cookie","wxuin="+webwx.Uin+";wxsid="+webwx.Sid+";webwx_data_ticket="+webwx.Webwx_data_ticket)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("[!] E00001: Please check your network connection.")
    }
    defer resp.Body.Close()
    fmt.Println("[+] 已退出微信.")
}
 
/**
 * 获取ChatSet
 * @param jsondata []byte 从webInit获取的Json数据
 * @return string 返回ChatSet
 * @return error 错误信息

 func GetChatSet(jsondata []byte) (string, error) {
    ChatSet, err := jsonparser.GetString(jsondata,"ChatSet")
    if err != nil {
        return "", fmt.Errorf("E00008: %s",err)
    }
    return ChatSet, nil
}
*/


/**
 * 获取SKey
 * @param jsondata []byte 从webInit获取的Json数据
 * @return string 返回SKey
 * @return error 错误信息

 func GetSKey(jsondata []byte) (string, error) {
    Skey, err := jsonparser.GetString(jsondata,"SKey")
    if err != nil {
        return "", fmt.Errorf("E00008: %s",err)
    }
    return Skey, nil
}
*/
