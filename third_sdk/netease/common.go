package netease

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hero1s/golib/log"
)

var (
	Code200 = errors.New("200")
	Code403 = errors.New("403")
	Code414 = errors.New("414")
	Code416 = errors.New("416")
	Code431 = errors.New("431")
	Code500 = errors.New("500")
)

type QueryUserRoomIdsType struct {
	Code uint64                   `json:"code"`
	Desc QueryUserRoomIdsDescType `json:"desc"`
}
type QueryUserRoomIdsDescType struct {
	Roomids []string `json:"roomids"`
}
type QueryMembersType struct {
	Code uint64               `json:"code"`
	Desc QueryMembersDescType `json:"desc"`
}
type QueryMembersDescType struct {
	Data []QueryMembersDataType `json:"data"`
}
type QueryMembersDataType struct {
	Roomid     uint64 `json:"roomid"`
	Ext        string `json:"ext"`
	Accid      string `json:"accid"`
	Nick       string `json:"nick"`
	Typet      string `json:"type"` //COMMON:普通成员(固定成员)；CREATOR:聊天室创建者；MANAGER:聊天室管理员；TEMPORARY:临时用户(非聊天室固定成员)；ANONYMOUS:匿名用户(未注册账号)；LIMITED:受限用户(黑名单+禁言)
	OnlineStat bool   `json:"onlineStat"`
}

type CreateUserID struct {
	Code int64 `json:"code"`
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
		Name  string `json:"name"`
		Icon  string `json:"icon"`
	} `json:"info"`
}

type MembersByPageType struct {
	Code uint64                `json:"code"`
	Desc MembersByPageDescType `json:"Desc"`
}
type MembersByPageDescType struct {
	Data []MembersByPageDataType `json:"data"`
}
type MembersByPageDataType struct {
	Roomid       uint64 `json:"roomid"`
	Accid        string `json:"accid"`
	Nick         string `json:"nick"`
	Avator       string `json:"avator"`
	Ext          string `json:"ext"`
	Typet        string `json:"type"`
	Level        uint64 `json:"level"`
	OnlineStat   bool   `json:"onlineStat"`
	EnterTime    uint64 `json:"enterTime"`
	Blacklisted  bool   `json:"blacklisted"`
	Muted        bool   `json:"muted"`
	TempMuted    bool   `json:"tempMuted"`
	TempMuteTtl  bool   `json:"tempMuteTtl"`
	IsRobot      bool   `json:"isRobot"`
	RobotExpirAt uint64 `json:"robotExpirAt"`
}
type MuteRoomType struct {
	Code int64            `json:"code"`
	Desc MuteRoomDescType `json:"desc"`
}
type MuteRoomDescType struct {
	success bool `json:"success"`
}
type SendMsgType struct {
	Code int64           `json:"code"`
	Desc SendMsgDescType `json:"desc"`
}
type TemporaryMuteType struct {
	Code int64                 `json:"code"`
	Desc TemporaryMuteDescType `json:"desc"`
}
type TemporaryMuteDescType struct {
	MuteDuration uint64 `json:"mute_duration"`
}

type RemoveRobotType struct {
	Code int64               `json:"code"`
	Desc RemoveRobotDescType `json:"desc"`
}
type RemoveRobotDescType struct {
	FailAccids    string `json:"failAccids"`
	SuccessAccids string `json:"successAccids"`
}

type AddRobotType struct {
	Code int64            `json:"code"`
	Desc AddRobotDescType `json:"desc"`
}
type AddRobotDescType struct {
	FailAccids    string `json:"failAccids"`
	SuccessAccids string `json:"successAccids"`
	oldAccids     string `json:"oldAccids"`
}

type SendMsgDescType struct {
	Time             string `json:"time"`
	FromAvator       string `json:"fromAvator"`
	MsgidClient      string `json:"msgid_client"`
	FromClientType   string `json:"fromClientType"`
	Attach           string `json:"attach"`
	RoomId           string `json:"roomId"`
	FromAccount      string `json:"fromAccount"`
	FromNick         string `json:"fromNick"`
	Type             string `json:"type"`
	Ext              string `json:"ext"`
	HighPriorityFlag uint64 `json:"highPriorityFlag"`
	MsgAbandonFlag   uint64 `json:"msgAbandonFlag"`
}

type TopNType struct {
	Code int64          `json:"code"`
	Data []TopNDataType `json:"data"`
}
type TopNDataType struct {
	ActiveNums uint64 `json:"activeNums"` // 该聊天室内的活跃数
	Datetime   uint64 `json:"datetime"`   // 统计时间点，单位秒，按天统计的是当天的0点整点；按小时统计的是指定小时的整点
	EnterNums  uint64 `json:"enterNums"`  // 进入人次数量
	Msgs       uint64 `json:"msgs"`       // 聊天室内发生的消息数
	Period     uint64 `json:"period"`     // 统计周期，HOUR表示按小时统计；DAY表示按天统计
	RoomId     uint64 `json:"Room_id"`    // 聊天室ID号
}

type CreateChatRoomType struct {
	Code     int64        `json:"code"`
	Chatroom ChatRoomType `json:"chatroom"`
}

//设置聊天室内用户角色返回值
//返回的type字段可能为：
//LIMITED,          //受限用户,黑名单+禁言
//COMMON,           //普通固定成员
//CREATOR,          //创建者
//MANAGER,          //管理员
//TEMPORARY,        //临时用户,非固定成员
type SetMemberRoleType struct {
	Code int64 `json:"code"`
	/*
		Desc struct {
			Roomid uint64 `json:"roomid"`
			Level uint64 `json:"level"`
			Accid string	`json:"accid"`
			Type string `json:"type"`
		} `json:"desc"`*/
}

type RequestAddr struct {
	Code int64    `json:"code"`
	Addr []string `json:"addr"`
}

type ChatRoomType struct {
	Roomid          uint64 `json:"roomid"`
	Valid           bool   `json:"valid"`
	Announcement    string `json:"announcement"`
	Name            string `json:"name"`
	Broadcasturl    string `json:"broadcasturl"`
	Ext             string `json:"ext"`
	Creator         string `json:"creator"`
	Onlineusercount uint64 `json:"onlineusercount"`
}

type GetBatchChatRoomType struct {
	Code         int64          `json:"code"`
	NoExistRooms []uint64       `json:"noExistRooms"`
	FailRooms    []uint64       `json:"failRooms"`
	SuccRooms    []ChatRoomType `json:"succRooms"`
}

type RefreshToken struct {
	Code int64 `json:"code"`
	Info struct {
		Token string `json:"token"`
		Accid string `json:"accid"`
	} `json:"info"`
}

type Uinfos struct {
	Accid  string `json:"accid"`
	Name   string `json:"name"`
	Props  string `json:"props"`
	Icon   string `json:"icon"`
	Sign   string `json:"sign"`
	Email  string `json:"email"`
	Birth  string `json:"birth"`
	Mobile string `json:"mobile"`
	Gender int64  `json:"gender"`
	Ex     string `json:"ex"`
}

type GetUserInfo struct {
	Code   int64    `json:"code"`
	Uinfos []Uinfos `json:"uinfos"`
}

type Friends struct {
	CreateTime  int64  `json:"create_time"` //返回的是微秒
	Bidirection bool   `json:"bidirection"`
	Faccid      string `json:"faccid"`
	Alias       string `json:"alias"`
}

type GetFriend struct {
	Code    int64     `json:"code"`
	Size    int64     `json:"size"`
	Friends []Friends `json:"friends"`
}

type BlockAndMute struct {
	Code      int64    `json:"code"`
	MuteList  []string `json:"mutelist"`
	BlackList []string `json:"blacklist"`
}

type SendMsg struct {
	Code int64 `json:"code"`
	Data struct {
		MsgId    int64 `json:"msgid"`
		AntiSpam bool  `json:"antispam"`
	} `json:"data"`
}

type OnlyCode struct {
	Code int64 `json:"code"`
}
type QueueListType struct {
	Code int64 `json:"code"`
	Desc struct {
		List []map[string]string `json:"list"`
	} `json:"desc"`
}
type BroadcastMsg struct {
	Code int64 `json:"code"`
	Msg  struct {
		ExpireTime  int64    `json:"expireTime"`
		Body        string   `json:"body"`
		CreateTime  int64    `json:"createTime"`
		IsOffline   bool     `json:"isOffline"`
		BroadcastId int64    `json:"broadcastId"`
		TargetOs    []string `json:"targetOs"`
	} `json:"msg"`
}
type CodeDescType struct {
	Code int64  `json:"code"`
	Desc string `json:"desc"`
}

func handleOnlyCodeResponse(rsp []byte, err error) error {
	if err != nil {
		return err
	}
	var r CodeDescType
	err = json.Unmarshal(rsp, &r)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		log.Error("---------IM error msg:", r.Desc)
		return errors.New(fmt.Sprintf("%v", r.Code))
	}
	return nil
}
