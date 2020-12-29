package netease

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hero1s/golib/log"
)

//消息功能

/*
给用户或者高级群发送普通消息，包括文本，图片，语音，视频和地理位置，具体消息参考下面描述

参数说明
参数	类型	必须	说明
from	String	是	发送者accid，用户帐号，最大32字符，
必须保证一个APP内唯一
ope	int	是	0：点对点个人消息，1：群消息（高级群），其他返回414
to	String	是	ope==0是表示accid即用户id，ope==1表示tid即群id
type	int	是	0 表示文本消息,
1 表示图片，
2 表示语音，
3 表示视频，
4 表示地理位置信息，
6 表示文件，
100 自定义消息类型（特别注意，对于未对接易盾反垃圾功能的应用，该类型的消息不会提交反垃圾系统检测）
body	String	是	最大长度5000字符，JSON格式。
具体请参考： 消息格式示例
antispam	String	否	对于对接了易盾反垃圾功能的应用，本消息是否需要指定经由易盾检测的内容（antispamCustom）。
true或false, 默认false。
只对消息类型为：100 自定义消息类型 的消息生效。
antispamCustom	String	否	在antispam参数为true时生效。
自定义的反垃圾检测内容, JSON格式，长度限制同body字段，不能超过5000字符，要求antispamCustom格式如下：

{"type":1,"data":"custom content"}

字段说明：
1. type: 1：文本，2：图片。
2. data: 文本内容or图片地址。
option	String	否	发消息时特殊指定的行为选项,JSON格式，可用于指定消息的漫游，存云端历史，发送方多端同步，推送，消息抄送等特殊行为;option中字段不填时表示默认值 ，option示例:

{"push":false,"roam":true,"history":false,"sendersync":true,"route":false,"badge":false,"needPushNick":true}

字段说明：
1. roam: 该消息是否需要漫游，默认true（需要app开通漫游消息功能）；
2. history: 该消息是否存云端历史，默认true；
3. sendersync: 该消息是否需要发送方多端同步，默认true；
4. push: 该消息是否需要APNS推送或安卓系统通知栏推送，默认true；
5. route: 该消息是否需要抄送第三方；默认true (需要app开通消息抄送功能);
6. badge:该消息是否需要计入到未读计数中，默认true;
7. needPushNick: 推送文案是否需要带上昵称，不设置该参数时默认true;
8. persistent: 是否需要存离线消息，不设置该参数时默认true。
pushcontent	String	否	推送文案，android以此为推送显示文案；ios若未填写payload，显示文案以pushcontent为准。超过500字符后，会对文本进行截断。
payload	String	否	推送对应的payload,必须是JSON,不能超过2k字符。
ext	String	否	开发者扩展字段，长度限制1024字符
forcepushlist	String	否	发送群消息时的强推用户列表（云信demo中用于承载被@的成员），格式为JSONArray，如["accid1","accid2"]。若forcepushall为true，则forcepushlist为除发送者外的所有有效群成员
forcepushcontent	String	否	发送群消息时，针对强推列表forcepushlist中的用户，强制推送的内容
forcepushall	String	否	发送群消息时，强推列表是否为群里除发送者外的所有有效成员，true或false，默认为false
bid	String	否	可选，反垃圾业务ID，实现“单条消息配置对应反垃圾”，若不填则使用原来的反垃圾配置
useYidun	int	否	可选，单条消息是否使用易盾反垃圾，可选值为0。
0：（在开通易盾的情况下）不使用易盾反垃圾而是使用通用反垃圾，包括自定义消息。

若不填此字段，即在默认情况下，若应用开通了易盾反垃圾功能，则使用易盾反垃圾来进行垃圾消息的判断
markRead	int	否	可选，群消息是否需要已读业务（仅对群消息有效），0:不需要，1:需要
checkFriend	boolean	否	是否为好友关系才发送消息，默认否
注：使用该参数需要先开通功能服务
*/
func (b *Net) SendMsg(data map[string]interface{}) (s SendMsg, err error) {
	url := "https://api.netease.im/nimserver/msg/sendMsg.action"
	log.Debugf("发送消息,map:%v",data)
	rsp, err := b.postDataHttps(url, data)
	log.Debugf("发送消息,返回值:%v",string(rsp[:]))
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(rsp, &s)
	if err != nil {
		return s, err
	}
	if s.Code != 200 {
		return s, errors.New(fmt.Sprintf("%v", s.Code))
	}

	return s, nil
}

//批量发送点对点普通消息
func (b *Net) SendBatchMsg(data map[string]interface{}) (s SendMsg, err error) {
	url := "https://api.netease.im/nimserver/msg/sendBatchMsg.action"
	log.Debugf("批量发送点对点普通消息,map:%v",data)
	rsp, err := b.postDataHttps(url, data)
	log.Debugf("批量发送点对点普通消息,返回值:%v",string(rsp[:]))
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(rsp, &s)
	if err != nil {
		return s, err
	}
	if s.Code != 200 {
		return s, errors.New(fmt.Sprintf("%v", s.Code))
	}

	return s, nil
}

//发送自定义系统通知
func (b *Net) SendAttachMsg(data map[string]interface{}) error {
	url := "https://api.netease.im/nimserver/msg/sendAttachMsg.action"
	err := handleOnlyCodeResponse(b.postDataHttps(url, data))
	if err == nil {
		log.Debugf("发送自定义系统通知成功,data:%v",data)
	}
	return err
}

//批量发送自定义系统通知
func (b *Net) SendBatchAttachMsg(data map[string]interface{}) error {
	url := "https://api.netease.im/nimserver/msg/sendBatchAttachMsg.action"
	err := handleOnlyCodeResponse(b.postDataHttps(url, data))
	if err == nil {
		log.Debugf("批量发送自定义系统通知成功,data:%v",data)
	}
	return err
}
/*
func (b *Net) SendBatchMsg(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/sendBatchMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) SendBatchAttachMsg(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/sendBatchAttachMsg.action"
	return b.postDataHttps(url, data)
}

func (b *Net) UploadMsg(content, typ, ishttps string) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/upload.action"
	data := map[string]interface{}{
		"content": content,
		"type":    typ,
		"ishttps": ishttps,
	}
	return b.postDataHttps(url, data)
}

func (b *Net) UploadMultiPartMsg(content, typ, ishttps string) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/fileUpload.action"
	data := map[string]interface{}{
		"content": content,
		"type":    typ,
		"ishttps": ishttps,
	}
	return b.postDataHttps(url, data)
}

func (b *Net) RecallMsg(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/recall.action"
	return b.postDataHttps(url, data)
}
*/

//广播消息
/*
参数   	   类型	   必须	  说明
body	   String	是	广播消息内容，最大4096字符
from	   String	否	发送者accid, 用户帐号，最大长度32字符，必须保证一个APP内唯一
isOffline  String	否	是否存离线，true或false，默认false
ttl	       int	    否	存离线状态下的有效期，单位小时，默认7天
targetOs   String	否	目标客户端，默认所有客户端，jsonArray，格式：["ios","aos","pc","web","mac"]
*/

func (b *Net) BroadcastMsg(data map[string]interface{}) (r BroadcastMsg, err error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	rsp, err := b.postDataHttps(url, data)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return r, err
	}
	fmt.Println("----rsp:", string(rsp))
	if r.Code != 200 {
		return r, errors.New(fmt.Sprintf("%v", r.Code))
	}
	return r, nil
}
