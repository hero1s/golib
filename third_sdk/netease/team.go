package netease

const (
	CreateGroupURL   = "https://api.netease.im/nimserver/team/create.action"
	AddIntoGroupURL  = "https://api.netease.im/nimserver/team/add.action"
	KickFromGroupURL = "https://api.netease.im/nimserver/team/kick.action"
	RemoveGroupURL   = "https://api.netease.im/nimserver/team/remove.action"
)

/*
//群组功能
func (b *Net) CreateGroup(data map[string]interface{}) (map[string]interface{}, error) {
	return b.postDataHttps(CreateGroupURL, data)
}
func (b *Net) AddIntoGroup(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) KickFromGroup(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) RemoveGroup(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}

func (b *Net) UpdateGroup(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) QueryGroup(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}

func (b *Net) QueryDetailGroup(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) GetMarkReadInfoGroup(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}

func (b *Net) ChangeGroupOwner(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) AddGroupManager(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) RemoveGroupManger(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) JoinTeams(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) UpdateGroupNick(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) MuteTeam(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) MuteTeamList(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
func (b *Net) LeaveTeam(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}

func (b *Net) MuteTeamListAll(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}

func (b *Net) ListTeamMute(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.netease.im/nimserver/msg/broadcastMsg.action"
	return b.postDataHttps(url, data)
}
*/
