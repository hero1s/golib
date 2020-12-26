package key

import (
	"math/rand"
	"strings"
	"time"
)

type Invite struct {
	base         []byte
	pad          byte // 补位
	baseMap      map[byte]int
	inviteLength int
}

var ShareCode *Invite

func init() {
	base := []byte{'W', 'E', '8', 'A', 'S', '2', 'D', 'Z', 'X', '9', 'C', '7', 'P', '5', 'I', 'K', '3', 'M', 'J', 'U', 'F', 'R', '4', 'V', 'Y', 'L', 'T', 'N', '6', 'B', 'G', 'H'}
	ShareCode = InitInvite(base, 'Q', 6)
}

// 补位码不能在base数组里,要不然无法区分
func InitInvite(base []byte, pad byte, inviteLength int) *Invite {
	for _, v := range base {
		if v == pad {
			panic("补位码不能在base数组里,要不然无法区分")
		}
	}
	invite := &Invite{
		base:         base,
		inviteLength: inviteLength,
		pad:          pad,
	}
	invite.baseMap = make(map[byte]int, len(base))
	for k, v := range base {
		invite.baseMap[v] = k
	}
	return invite
}

// 生成邀请码
func (in *Invite) CreateInviteCode(uid uint64) string {
	var code string
	baseLength := uint64(len(in.base))
	for uid != 0 {
		mod := uid % baseLength
		uid = uid / baseLength
		code = code + string(in.base[mod])
	}
	codeLength := len(code)
	if codeLength < in.inviteLength {
		code = code + string(in.pad)
		for i := 0; i < in.inviteLength-codeLength-1; i++ {
			rand.Seed(time.Now().UnixNano())
			code = code + string(in.base[rand.Intn(int(baseLength))])
		}
	}
	return code
}

// 反解邀请码
func (in *Invite) DecodeInviteCode(code string) uint64 {
	codeLength := len(code)
	if i := strings.Index(code, string(in.pad)); i != -1 {
		codeLength = i
	}
	var r int
	var uid uint64
	baseLength := len(in.base)
	for i := 0; i < codeLength; i++ {
		if string(code[i]) == string(in.pad) {
			continue
		}
		index := in.baseMap[code[i]]
		b := 1
		for j := 0; j < r; j++ {
			b = b * baseLength
		}
		uid = uid + uint64(index*b)
		r++
	}
	return uid
}
