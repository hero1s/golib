package netease

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/utils/uuid"
	"strings"
)

//查询用户创建的开启状态聊天室列表
func (b *Net)QueryUserRoomIds(accid string) ([]string,error) {
	var data []string
	url := "https://api.netease.im/nimserver/chatroom/queryUserRoomIds.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"creator":accid})
	log.Debugf("批量获取在线成员信息,post返回%v",string(rsp[:]))
	if err != nil {
		return data,err
	}
	var r QueryUserRoomIdsType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return data,err
	}
	if r.Code != 200 {
		log.Errorf("批量获取在线成员信息,返回:%#v",r)
		return data,err
	}
	return r.Desc.Roomids,nil
}

//变更聊天室内的角色信息
/*
roomid 	long	是	聊天室id
accid 	String 	是	需要变更角色信息的accid
上面的当成单个参数
下面的放到map里面
save 	boolean 	否	变更的信息是否需要持久化，默认false，仅对聊天室固定成员生效
needNotify 	boolean 	否	是否需要做通知
notifyExt 	String 	否	通知的内容，长度限制2048
nick 	String 	否	聊天室室内的角色信息：昵称，不超过64个字符
avator 	String 	否	聊天室室内的角色信息：头像
ext 	String 	否	聊天室室内的角色信息：开发者扩展字段
*/
func (b *Net)UpdateMyRoomRole(roomid uint64,accid string,pm map[string]interface{}) (error) {
	//请注意看这里,要更新的key都在里面,千万不要写其它key,要更新哪个写哪个,不更新不要写
	keys := map[string]bool{
		"save":true,
		"needNotify":true,
		"notifyExt":true,
		"nick":true,
		"avator":true,
		"ext":true,
	}
	for key,_ := range pm{
		if _,ok := keys[key];ok == false {
			log.Errorf("变更聊天室内的角色信息,key:%v 不存在", key)
			return errors.New("key:" + key + "不存在")
		}
	}
	pm["roomid"] = roomid
	pm["accid"] = accid
	url := "https://api.netease.im/nimserver/chatroom/updateMyRoomRole.action"
	rsp, err := b.postDataHttps(url, pm)
	log.Debugf("变更聊天室内的角色信息,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r CodeDescType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Errorf("变更聊天室内的角色信息,返回:%#v",r)
		return errors.New(r.Desc)
	}
	return nil
}

//批量获取在线成员信息
/*
roomid 	long	是	聊天室id
accids 	JSONArray 	是	["abc","def"], 账号列表，最多200条
 */
func (b *Net) QueryMembers(roomid uint64,accids []string) ([]QueryMembersDataType, error) {
	var data []QueryMembersDataType
	ids_string := strings.Replace(strings.Trim(fmt.Sprint(accids), "[]"), " ", "\",\"", -1) //切片拼接字符串
	ids_string = "[\"" + ids_string + "\"]"
	log.Debugf("批量获取在线成员信息,ids:%v",ids_string)
	url := "https://api.netease.im/nimserver/chatroom/queryMembers.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid,"accids":ids_string})
	log.Debugf("批量获取在线成员信息,post返回%v",string(rsp[:]))
	if err != nil {
		return data,err
	}
	var r QueryMembersType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return data,err
	}
	if r.Code != 200 {
		log.Errorf("批量获取在线成员信息,返回:%#v",r)
		return data,err
	}
	return r.Desc.Data,nil
}
//分页获取成员列表
/*
roomid 	long	是	聊天室id
type 	int	是	需要查询的成员类型,0:固定成员;1:非固定成员;2:仅返回在线的固定成员
endtime 	long	是	单位毫秒，按时间倒序最后一个成员的时间戳,0表示系统当前时间
limit 	long	是	返回条数，<=100
 */
func (b *Net)MembersByPage(roomid uint64,typet uint64,endtime uint64,limit uint64) ([]MembersByPageDataType,error) {
	var rlt []MembersByPageDataType
	url := "https://api.netease.im/nimserver/chatroom/membersByPage.action"
	data := map[string]interface{}{"roomid":roomid,"type":typet,"endtime":endtime,"limit":limit}

	rsp, err := b.postDataHttps(url, data)
	log.Debugf("分页获取成员列表,post返回%v",string(rsp[:]))
	if err != nil {
		return rlt,err
	}
	var r MembersByPageType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return rlt,err
	}
	if r.Code != 200 {
		log.Errorf("分页获取成员列表,返回:%#v",r)
		return rlt,err
	}
	return r.Desc.Data,nil
}

type ChannelInfoType struct {
	Cid uint64 `json:"cid"`
	Cname string `json:"cname"`
	Accid string `json:"accid"`
	Total int64 `json:"total"`
	Mode int64 `json:"mode"`
	Status int64 `json:"status"` // 房间状态【1：初始状态，2：进行中，3：正常结束，4：异常结束】
	Createtime int64 `json:"createtime"`
	Destroytime int64 `json:"destroytime"`
}
//获得音视频房间信息
/*
channelid 	long	是	聊天室id
*/
func (b *Net)GetChannelid(channelid uint64) (ChannelInfoType,error) {
	var rlt ChannelInfoType
	url := fmt.Sprintf("https://roomserver-dev.netease.im/v1/api/rooms/%v",channelid)

	data := map[string]interface{}{}

	rsp, err := b.getDataHttps(url, data)
	log.Debugf("获得音视频房间信息,get返回%v",string(rsp[:]))
	if err != nil {
		return rlt,err
	}
	err = json.Unmarshal(rsp, &rlt)
	if err != nil {
		return rlt,err
	}
	return rlt,nil
}

//查询聊天室统计指标TopN
/*
topn 	int 	否	topn值，可选值 1~500，默认值100
timestamp 	long	否	需要查询的指标所在的时间坐标点，不提供则默认当前时间，单位秒/毫秒皆可
period 	String 	否	统计周期，可选值包括 hour/day, 默认hour
orderby 	String 	否	取排序值,可选值 active/enter/message,分别表示按日活排序，进入人次排序和消息数排序， 默认active
 */
func (b *Net)TopN(pm map[string]interface{}) ([]TopNDataType,error) {
	//请注意看这里,要更新的key都在里面,千万不要写其它key,要更新哪个写哪个,不更新不要写
	var rlt []TopNDataType
	keys := map[string]bool{
		"topn":true,
		"timestamp":true,
		"period":true,
		"orderby":true,
	}
	for key,_ := range pm{
		if _,ok := keys[key];ok == false {
			log.Errorf("查询聊天室统计指标TopN,key:%v 不存在", key)
			return rlt,errors.New("key:" + key + "不存在")
		}
	}

	url := "https://api.netease.im/nimserver/stats/chatroom/topn.action"
	rsp, err := b.postDataHttps(url, pm)
	log.Debugf("查询聊天室统计指标TopN,post返回%v",string(rsp[:]))
	if err != nil {
		return rlt,err
	}
	var r TopNType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return rlt,err
	}
	if r.Code != 200 {
		log.Errorf("创建聊天室失败,返回:%#v",r)
		return rlt,err
	}
	return r.Data,nil
}
// 将聊天室整体禁言
/*
roomid 	long	是	聊天室id
operator 	String 	是	操作者accid，必须是管理员或创建者
mute 	String 	是	true或false
needNotify 	String 	否	true或false，默认true
notifyExt 	String 	否	通知扩展字段
*/
func (b *Net)MuteRoom(roomid uint64,operator string,mute bool,needNotify bool,notifyExt string) error {
	url := "https://api.netease.im/nimserver/chatroom/muteRoom.action"
	data := map[string]interface{}{"roomid":roomid,"operator":operator,"mute":mute}
	if needNotify {
		data["needNotify"] = needNotify
	}
	if notifyExt != "" {
		data["notifyExt"] = notifyExt
	}

	rsp, err := b.postDataHttps(url, data)
	log.Debugf("将聊天室整体禁言,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r MuteRoomType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Errorf("将聊天室整体禁言,返回:%#v",r)
		return err
	}
	return nil
}
//设置临时禁言状态
/*
roomid 	long	是	聊天室id
operator 	String 	是 	操作者accid,必须是管理员或创建者
target 	String 	是 	被禁言的目标账号accid
muteDuration 	long	是 	0:解除禁言;>0设置禁言的秒数，不能超过2592000秒(30天)
needNotify 	String 	否 	操作完成后是否需要发广播，true或false，默认true
notifyExt 	String 	否 	通知广播事件中的扩展字段，长度限制2048字符
*/
func (b *Net)TemporaryMute(roomid uint64,operator string,target string,muteDuration uint64,needNotify bool,notifyExt string) error {
	url := "https://api.netease.im/nimserver/chatroom/temporaryMute.action"
	data := map[string]interface{}{"roomid":roomid,"operator":operator,"target":target,"muteDuration":muteDuration}
	if needNotify{
		data["needNotify"] = needNotify
	}
	if notifyExt != "" {
		data["notifyExt"] = notifyExt
	}

	rsp, err := b.postDataHttps(url, data)
	log.Debugf("设置临时禁言状态,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r TemporaryMuteType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Errorf("设置临时禁言状态,返回:%#v",r)
		return err
	}
	return nil
}

//从聊天室内删除机器人
/*
roomid 	long	是	聊天室id
accids 	JSONArray 	是 	机器人账号accid列表，必须是有效账号，账号数量上限100个*/
func (b *Net) RemoveRobot(roomid uint64,accids []string) error {
	ids_string := strings.Replace(strings.Trim(fmt.Sprint(accids), "[]"), " ", "\",\"", -1) //切片拼接字符串
	ids_string = "[\"" + ids_string + "\"]"
	log.Debugf("从聊天室内删除机器人,ids:%v",ids_string)
	url := "https://api.netease.im/nimserver/chatroom/removeRobot.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid,"accids":ids_string})
	log.Debugf("从聊天室内删除机器人,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r RemoveRobotType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Errorf("从聊天室内删除机器人,返回:%#v",r)
		return err
	}
	return nil
}

//往聊天室里加机器人,机器人过期时间为24小时。
/*
roomid 	long	是	聊天室id
accids 	JSONArray 	是 	机器人账号accid(string)列表，必须是有效账号，账号数量上限100个
roleExt 	String 	否	机器人信息扩展字段，请使用json格式，长度4096字符
notifyExt 	String 	否	机器人进入聊天室通知的扩展字段，请使用json格式，长度2048字符*/
func (b *Net) AddRobot(pm map[string]interface{}) error {
	//请注意看这里,要更新的key都在里面,千万不要写其它key,要更新哪个写哪个,不更新不要写
	keys := map[string]bool{
		"roomid":true,
		"accids":true,
		"roleExt":true,
		"notifyExt":true,
	}
	for key,_ := range pm{
		if _,ok := keys[key];ok == false{
			log.Errorf("往聊天室里加机器人,key:%v 不存在",key)
			return errors.New("key:"+ key+"不存在")
		}
	}
	ids_string := strings.Replace(strings.Trim(fmt.Sprint(pm["accids"]), "[]"), " ", "\",\"", -1) //切片拼接字符串
	ids_string = "[\"" + ids_string + "\"]"
	pm["accids"] = ids_string
	log.Debugf("往聊天室里加机器人,ids:%v",ids_string)
	url := "https://api.netease.im/nimserver/chatroom/addRobot.action"
	rsp, err := b.postDataHttps(url, pm)
	log.Debugf("往聊天室里加机器人,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r AddRobotType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Errorf("往聊天室里加机器人,返回:%#v",r)
		return err
	}
	return nil
}

//发送聊天室消息
/*  这里是pm的字段key值
roomid 		long	是	聊天室id
fromAccid 	String 	是	消息发出者的账号accid
msgType 	int 	是	消息类型：
						0: 表示文本消息，
						1: 表示图片，
						2: 表示语音，
						3: 表示视频，
						4: 表示地理位置信息，
						6: 表示文件，
						10: 表示Tips消息，
						100: 自定义消息类型（特别注意，对于未对接易盾反垃圾功能的应用，该类型的消息不会提交反垃圾系统检测）
resendFlag 	int 	否	重发消息标记，0：非重发消息，1：重发消息，如重发消息会按照msgid检查去重逻辑
attach 		String 	否 	消息内容，格式同消息格式示例中的body字段,长度限制4096字符
	ext 	String 	否	消息扩展字段，内容可自定义，请使用JSON格式，长度限制4096字符
antispam 	String 	否	对于对接了易盾反垃圾功能的应用，本消息是否需要指定经由易盾检测的内容（antispamCustom）。
						true或false, 默认false。
						只对消息类型为：100 自定义消息类型 的消息生效。
antispamCustom 	String	否	在antispam参数为true时生效。
							自定义的反垃圾检测内容, JSON格式，长度限制同body字段，不能超过5000字符，要求antispamCustom格式如下：
							{"type":1,"data":"custom content"}
							字段说明：
							1. type: 1：文本，2：图片。
							2. data: 文本内容or图片地址。
skipHistory 	int 	否	是否跳过存储云端历史，0：不跳过，即存历史消息；1：跳过，即不存云端历史；默认0
bid 			String 	否	可选，反垃圾业务ID，实现“单条消息配置对应反垃圾”，若不填则使用原来的反垃圾配置
highPriority 	Boolean 否	可选，true表示是高优先级消息，云信会优先保障投递这部分消息；false表示低优先级消息。默认false。
							强烈建议应用恰当选择参数，以便在必要时，优先保障应用内的高优先级消息的投递。若全部设置为高优先级，则等于没有设置。
useYidun 		int 	否	可选，单条消息是否使用易盾反垃圾，可选值为0。
							0：（在开通易盾的情况下）不使用易盾反垃圾而是使用通用反垃圾，包括自定义消息。
							若不填此字段，即在默认情况下，若应用开通了易盾反垃圾功能，则使用易盾反垃圾来进行垃圾消息的判断
needHighPriorityMsgResend 	Boolean 	否	可选，true表示会重发消息，false表示不会重发消息。默认true。注:若设置为true，
											用户离开聊天室之后重新加入聊天室，在有效期内还是会收到发送的这条消息，目前有效期默认30s。在没有配置highPriority时needHighPriorityMsgResend不生效。
abandonRatio 	int 	否	可选，消息丢弃的概率。取值范围[0-9999]；
							其中0代表不丢弃消息，9999代表99.99%的概率丢弃消息，默认不丢弃；
							注意如果填写了此参数，highPriority参数则会无效；
							此参数可用于流控特定业务类型的消息。*/
func (b *Net) SendMsgChatRoom(pm map[string]interface{}) error  {
	//请注意看这里,要更新的key都在里面,千万不要写其它key,要更新哪个写哪个,不更新不要写
	keys := map[string]bool{
		"roomid":true,
		"msgId":true,
		"fromAccid":true,
		"msgType":true,
		"resendFlag":true,
		"attach":true,
		"ext":true,
		"antispam":true,
		"antispamCustom":true,
		"skipHistory":true,
		"bid":true,
		"highPriority":true,
		"useYidun":true,
		"needHighPriorityMsgResend":true,
		"abandonRatio":true,
	}
	for key,_ := range pm{
		if _,ok := keys[key];ok == false{
			log.Errorf("发送聊天室消息,key:%v 不存在",key)
			return errors.New("key:"+ key+"不存在")
		}
	}

	url := "https://api.netease.im/nimserver/chatroom/sendMsg.action"
	pm["msgId"] = uuid.GenStringUUID()
	rsp, err := b.postDataHttps(url, pm)
	log.Debugf("发送聊天室消息,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r SendMsgType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		log.Errorf("发送聊天室消息,post返回%v",string(rsp[:]))
		return err
	}
	if r.Code != 200 {
		log.Errorf("发送聊天室消息,返回码:%#v,返回:%v",r,string(rsp[:]))
		return err
	}
	/*
	str,_ := json.Marshal(pm)
	db_data := map[string]interface{}{
		"roomid":pm["roomid"],
		"content":string(str[:]),
		"ptime":time.Now().Unix(),
		"msgid":pm["msgId"],
	}
	db_table.TSendChatLog.NewTableRecord(db_data)*/
	return nil
}

//请求聊天室地址
//roomid 	long	是	聊天室id
//accid 	String 	是 	进入聊天室的账号
//clienttype 	int 	否	1:weblink（客户端为web端时使用）; 2:commonlink（客户端为非web端时使用）;3:wechatlink(微信小程序使用), 默认1
//clientip 	String	否	客户端ip，传此参数时，会根据用户ip所在地区，返回合适的地址
func (b *Net) RequestAddr(roomid uint64,accid string,clienttype int,clientip string) ([]string,error) {
	var addrs []string
	url := "https://api.netease.im/nimserver/chatroom/requestAddr.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{
		"roomid":roomid,"accid":accid,"clienttype":clienttype/*,"clientip":clientip上线之后把这段打开,内网测试取不到ip*/})
	log.Debugf("请求聊天室地址,post返回%v",string(rsp[:]))
	if err != nil {
		return addrs,err
	}
	var r RequestAddr
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return addrs,err
	}
	if r.Code != 200 {
		log.Errorf("请求聊天室地址,返回:%#v",r)
		return addrs,err
	}
	return r.Addr,nil
}

//设置聊天室内用户角色
//roomid 	long	是	聊天室id
//operator 	String 	是 	操作者账号accid
//target 	String 	是	被操作者账号accid
//opt 	int 	是	操作：
//1: 设置为管理员，operator必须是创建者
//2:设置普通等级用户，operator必须是创建者或管理员
//-1:设为黑名单用户，operator必须是创建者或管理员
//-2:设为禁言用户，operator必须是创建者或管理员
//optvalue 	String 	是	true或false，true:设置；false:取消设置；
//执行“取消”设置后，若成员非禁言且非黑名单，则变成游客
//notifyExt 	String 	否	通知扩展字段，长度限制2048，请使用json格式
func (b *Net) SetMemberRole(roomid uint64,operator string,target string,opt int64,optvalue bool,notifyExt string) error {
	url := "https://api.netease.im/nimserver/chatroom/setMemberRole.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{
		"roomid":roomid,"operator":operator,"target":target,"opt":opt,"optvalue":optvalue,"notifyExt":notifyExt})
	log.Debugf("设置聊天室内用户角色,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r SetMemberRoleType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Errorf("设置聊天室成员,返回码:%#v,返回值:%v",r,string(rsp[:]))
		return errors.New(string(rsp[:]))
	}
	return nil
}

func (b *Net) CreateChatRoom(accId string,roomName string) (uint64, error) {
	url := "https://api.netease.im/nimserver/chatroom/create.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"creator": accId,"name":roomName})
	log.Debugf("创建聊天室,post返回%v",string(rsp[:]))
	if err != nil {
		return 0, err
	}
	var r CreateChatRoomType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return 0, err
	}

	if r.Code != 200 {
		log.Errorf("创建聊天室失败,返回:%#v",r)
		return 0,err
	}
	return r.Chatroom.Roomid, nil
}

/**
@param roomid int 聊天室id,注意是云信那边的
@param needOnlineUserCount bool 是否需要返回在线人数,true或者false,默认false
 */
func (b *Net) GetChatRoom(roomid uint64,needOnlineUserCount bool)(ChatRoomType,error)  {
	var r CreateChatRoomType
	var rlt ChatRoomType
	url := "https://api.netease.im/nimserver/chatroom/get.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid": roomid,"needOnlineUserCount":needOnlineUserCount})
	log.Debugf("获得聊天室信息,post返回%v",string(rsp[:]))
	if err != nil {
		return rlt, err
	}
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return rlt, err
	}
	if r.Code != 200 {
		log.Errorf("创建聊天室失败,返回:%#v",r)
		return rlt,err
	}
	return r.Chatroom,nil
}

//批量获得聊天室信息
func (b *Net) GetBatchChatRoom(ids []uint64,needOnlineUserCount bool)([]ChatRoomType,error)  {
	var data []ChatRoomType
	url := "https://api.netease.im/nimserver/chatroom/getBatch.action"
	if len(ids) > 20 {
		log.Error("批量获取房间信息最大限制id数为20,ids:%v",ids)
		return data,errors.New("房间id大于20")
	}
	ids_string := strings.Replace(strings.Trim(fmt.Sprint(ids), "[]"), " ", ",", -1) //切片拼接字符串
	ids_string = "[" + ids_string + "]"
	log.Debugf("批量获取房间信息,ids:%v",ids_string)
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomids": ids_string,"needOnlineUserCount":needOnlineUserCount})
	log.Debug("指获得聊天室信息,post返回%v",string(rsp[:]))
	if err != nil {
		return data, err
	}
	var r GetBatchChatRoomType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return data, err
	}
	if r.Code != 200 {
		log.Error("创建聊天室失败,返回:%#v",r)
		return data,err
	}
	for i := 0;i< len(r.SuccRooms);i++{
		data = append(data,r.SuccRooms[i])
	}
	return data,nil
}


//修改关闭聊天室状态
//@param roomid 网易云信的聊天室id int64
//@param operator 网易云信的创建者id string
//@param valid true或false，false:关闭聊天室；true:打开聊天室 设置为false之后,返回值有个muted(静音)参数,有需要请取出
func (b *Net) ToggleCloseStat(roomid uint64,operator string,valid bool) error  {
	url := "https://api.netease.im/nimserver/chatroom/toggleCloseStat.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid,"operator":operator,"valid":valid})
	log.Debug("修改聊天室开/关闭状态,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}

	type Result struct {
		Code int `json:"code"`
		Desc interface{} `json:"desc"`
	}

	ret := Result{}
	err = json.Unmarshal(rsp, &ret)
	if err != nil {
		return err
	}

	if ret.Code != 200 {
		log.Error("修改聊天室开/关闭状态,返回:%#v", ret)
		return errors.New("修改聊天室开/关闭状态："+ fmt.Sprintf("%v", ret.Desc))
	}

	return nil
}

type  toggleCloseStatType struct {
	Code int64 `json:"code"`
	Desc  ChatRoomType `json:"desc"`
}

type QueueInitType struct {
	Code int64 `json:"code"`
}
//初始化队列
//roomid 	long	是	聊天室id
//sizeLimit 	long	是	队列长度限制，0~1000
func (b *Net) QueueInit(roomid uint64,sizeLimit uint64) error {
	url := "https://api.netease.im/nimserver/chatroom/queueInit.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid,"sizeLimit":sizeLimit})
	log.Debug("初始化队列,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r QueueInitType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Error("初始化队列,返回:%#v",r)
		return err
	}
	return nil
}

//往聊天室有序队列中新加或更新元素
/*
roomid 	long	是	聊天室id
key 	String 	是 	elementKey,新元素的UniqKey,长度限制128字符
value 	String 	是 	elementValue,新元素内容，长度限制4096字符
operator 	String 	否 	提交这个新元素的操作者accid，默认为该聊天室的创建者，若operator对应的帐号不存在，会返回404错误。
若指定的operator不在线，则添加元素成功后的通知事件中的操作者默认为聊天室的创建者；若指定的operator在线，则通知事件的操作者为operator。
transient 	String 	否 	这个新元素的提交者operator的所有聊天室连接在从该聊天室掉线或者离开该聊天室的时候，提交的元素是否需要删除。
true：需要删除；false：不需要删除。默认false。
当指定该参数为true时，若operator当前不在该聊天室内，则会返回403错误。*/
func (b *Net) QueueOffer(roomid uint64,key string,value string,operator string,transient bool) error {
	url := "https://api.netease.im/nimserver/chatroom/queueOffer.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid,
		"key":key,"value":value,"operator":operator,"transient":transient})
	log.Debug("往聊天室有序队列中新加或更新元素,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r OnlyCode
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Error("往聊天室有序队列中新加或更新元素,返回:%#v",r)
		return err
	}
	return nil
}
/*
批量更新聊天室队列元素
roomid 	long	是	聊天室id
operator 	String 	是	操作者accid,必须是管理员或创建者
elements 	String 	是	更新的key-value对，最大200个，示例：{"k1":"v1","k2":"v2"}
needNotify 	boolean 	否	true或false,是否需要发送更新通知事件，默认true
notifyExt 	String 	否	通知事件扩展字段，长度限制2048
*/
func (b *Net) QueueBatchUpdateElements(roomid uint64,elements string,operator string,needNotify bool ,notifyExt string) error {
	url := "https://api.netease.im/nimserver/chatroom/queueBatchUpdateElements.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid,
		"elements":elements,"needNotify":needNotify,"operator":operator,"notifyExt":notifyExt})
	log.Debug("批量更新聊天室队列元素,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r OnlyCode
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Debug("批量更新聊天室队列元素,返回:%#v",r)
		return err
	}
	return nil
}

//删除聊天室队列
func (b *Net)QueueDrop(roomid uint64) error {
	url := "https://api.netease.im/nimserver/chatroom/queueDrop.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid})
	log.Debug("删除聊天室队列,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r OnlyCode
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Error("删除聊天室队列,返回:%#v",r)
		return err
	}
	return nil
}

//删除聊天室队列
func (b *Net)QueuePoll(roomid uint64,key string) error {
	url := "https://api.netease.im/nimserver/chatroom/queuePoll.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid,"key":key})
	log.Debug("删除聊天室队列元素,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r OnlyCode
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		//log.Error("删除聊天室队列元素,返回:%#v",r)
		return errors.New(fmt.Sprintf("删除聊天室队列元素,返回:%#v",r))
	}
	return nil
}

//list聊天室队列
func (b *Net)QueueList(roomid uint64) QueueListType {
	var r QueueListType
	url := "https://api.netease.im/nimserver/chatroom/queueList.action"
	rsp, err := b.postDataHttps(url, map[string]interface{}{"roomid":roomid})
	log.Error("排序列出队列中所有元素,post返回%v",string(rsp[:]))
	if err != nil {
		return r
	}

	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return r
	}
	if r.Code != 200 {
		log.Error("list聊天室队列元素,返回:%#v",r)
		return r
	}
	return r
}
//更新聊天室信息
func (b *Net) UpdateChatRoom(roomid uint64,pm map[string]interface{}) error  {
	//请注意看这里,要更新的key都在里面,千万不要写其它key,要更新哪个写哪个,不更新不要写
	keys := map[string]bool{
		"roomid":true,//房间id uint64
		"name":true,//聊天室名称,string
		"announcement":true,//公告,string
		"broadcasturl":true,//直播
		"ext":true,//扩展字符,string
		"needNotify":true,//true或false,是否需要发送更新通知事件，默认true
		"notifyExt":true,//通知事件扩展字段，长度限制2048
		"queuelevel":true,//队列管理权限：0:所有人都有权限变更队列，1:只有主播管理员才能操作变更
	}
	for key,_ := range pm{
		if _,ok := keys[key];ok == false{
			log.Error("更新房间信息,key:%v 不存在",key)
			return errors.New("key:"+ key+"不存在")
		}
	}

	url := "https://api.netease.im/nimserver/chatroom/update.action"
	pm["roomid"] = roomid
	rsp, err := b.postDataHttps(url, pm)
	log.Debug("更新聊天室信息,post返回%v",string(rsp[:]))
	if err != nil {
		return err
	}
	var r CreateChatRoomType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Error("创建聊天室失败,返回:%#v",r)
		return err
	}
	return nil
}





