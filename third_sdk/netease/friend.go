package netease

import (
	"encoding/json"
	"errors"
	"fmt"
)

//用户关系托管
/*
参数	类型	必须	说明
accid	String	是	加好友发起者accid
faccid	String	是	加好友接收者accid
type	int	是	1直接加好友，2请求加好友，3同意加好友，4拒绝加好友
msg	String	否	加好友对应的请求消息，第三方组装，最长256字符
serverex	String	否	服务器端扩展字段，限制长度256
此字段client端只读，server端读写
*/
func (b *Net) AddFriend(accId, faccId string, typ int64, msg string) error {
	url := "https://api.netease.im/nimserver/friend/add.action"
	data := map[string]interface{}{
		"accid":  accId,
		"faccid": faccId,
		"type":   typ,
		"msg":    msg,
	}
	return handleOnlyCodeResponse(b.postDataHttps(url, data))
}

/*
参数	类型	必须	说明
accid	String	是	发起者accid
faccid	String	是	要修改朋友的accid
alias	String	否	给好友增加备注名，限制长度128，可设置为空字符串
ex	String	否	修改ex字段，限制长度256，可设置为空字符串
serverex	String	否	修改serverex字段，限制长度256，可设置为空字符串
此字段client端只读，server端读写
*/
func (b *Net) UpdateFriend(accId, faccId string, typ int64, msg, ex string) error {
	url := "https://api.netease.im/nimserver/friend/update.action"
	data := map[string]interface{}{
		"accid":  accId,
		"faccid": faccId,
		"type":   typ,
		"msg":    msg,
		"ex":     ex,
	}
	return handleOnlyCodeResponse(b.postDataHttps(url, data))
}

/*
删除好友关系
参数	类型	必须	说明
accid	String	是	发起者accid
faccid	String	是	要删除朋友的accid
isDeleteAlias	Boolean	否	是否需要删除备注信息
默认false:不需要，true:需要
*/
func (b *Net) DeleteFriend(accId, faccId string) error {
	url := "https://api.netease.im/nimserver/friend/delete.action"
	data := map[string]interface{}{
		"accid":  accId,
		"faccid": faccId,
	}
	return handleOnlyCodeResponse(b.postDataHttps(url, data))
}

/*
查询某时间点起到现在有更新的双向好友
accid	String	是	发起者accid
updatetime	Long	是	更新时间戳，接口返回该时间戳之后有更新的好友列表
createtime	Long	否	【Deprecated】定义同updatetime
*/

func (b *Net) GetFriend(accId string, updateTime int64) (r GetFriend, err error) {
	url := "https://api.netease.im/nimserver/friend/get.action"
	//updatetime 要求是微秒
	data := map[string]interface{}{
		"accid":      accId,
		"updatetime": updateTime * 1000,
	}
	rsp, err := b.postDataHttps(url, data)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return r, err
	}
	if r.Code != 200 {
		return r, errors.New(fmt.Sprintf("%v", r.Code))
	}
	return r, nil
}

/*
拉黑/取消拉黑；设置静音/取消静音

参数说明
参数	类型	必须	说明
accid	String	是	用户帐号，最大长度32字符，必须保证一个
APP内唯一
targetAcc	String	是	被加黑或加静音的帐号
relationType	int	是	本次操作的关系类型,1:黑名单操作，2:静音列表操作
value	int	是	操作值，0:取消黑名单或静音，1:加入黑名单或静音

*/

func (b *Net) SpecializeFriend(accid, targetAcc string, relationType, value int64) error {
	url := "https://api.netease.im/nimserver/user/setSpecialRelation.action"
	data := map[string]interface{}{
		"accid":        accid,
		"targetAcc":    targetAcc,
		"relationType": relationType,
		"value":        value,
	}
	return handleOnlyCodeResponse(b.postDataHttps(url, data))
}

/*
查看用户的黑名单和静音列表

参数说明
参数	类型	必须	说明
accid	String	是	用户帐号，最大长度32字符，必须保证一个
APP内唯一
*/

func (b *Net) ListBlockAndMuteFriend(accId string) (r BlockAndMute, err error) {
	url := "https://api.netease.im/nimserver/user/listBlackAndMuteList.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"accid": accId})
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return r, err
	}
	if r.Code != 200 {
		return r, errors.New(fmt.Sprintf("%v", r.Code))
	}
	return r, nil
}
