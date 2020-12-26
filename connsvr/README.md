## 长连接服务框架
#### 1:同时支持tcp,websocket连接
#### 2:同时支持protobuff,json格式消息解析
#### 3:支持集群互联(但是基本用不上,5w同时在线以上才考虑集群要不增加复杂性未必能增加稳定性)

#### 说明
##### cluster模块，处理集群互联
##### gate模块，处理前端链接

#### 使用方式
- 1:实现一个gate的Module

  ```
  type Module struct {
  	*gate.Gate
  }
  
  func (m *Module) OnInit() {
  	m.Gate = &gate.Gate{
  		MaxConnNum:      conf.Server.MaxConnNum,
  		PendingWriteNum: conf.PendingWriteNum,
  		MaxMsgLen:       conf.MaxMsgLen,
  		WSAddr:          conf.Server.WSAddr,
  		HTTPTimeout:     conf.HTTPTimeout,
  		TCPAddr:         conf.Server.TCPAddr,
  		LenMsgLen:       conf.LenMsgLen,
  		LittleEndian:    conf.LittleEndian,
  		AgentChanRPC:    game.ChanRPC,
  	}
  
  	switch conf.Encoding {
  	case "json":
  		m.Gate.Processor = msg.JSONProcessor
  	case "protobuf":
  		m.Gate.Processor = msg.ProtobufProcessor
  	default:
  		log.Error("unknown encoding: %v", conf.Encoding)
  	}
  }
  ```
  
- 2:实现一个处理消息的模块(分不分模块看业务需求)
  
  ```
  var (
  	skeleton = base.NewSkeleton()
  	ChanRPC  = skeleton.ChanRPCServer
  )
  
  type Module struct {
  	*module.Skeleton
  }
  
  func (m *Module) OnInit() {
  	m.Skeleton = skeleton
  
  	SubscribeRedisMsg()
  
  }
  
  func (m *Module) OnDestroy() {
  
  }
  ```
- 3:注册路由消息
  ```
  var (
  	JSONProcessor     = json.NewProcessor()
  	ProtobufProcessor = protobuf.NewProcessor()
  )
  
  func init() {
  	JSONProcessor.Register(&C2S_UserHeartBeat{})
  	JSONProcessor.Register(&C2S_UserLogin{})
  }
  
  //心跳包
  type C2S_UserHeartBeat struct {
  	Ptime int64 `json:"ptime" desc:"当前时间"`
  }
  
  //用户登录
  type C2S_UserLogin struct {
  	Token     string `json:"token" desc:"平台token"`
  	LoginPlat int    `json:"login_plat" desc:登录平台`
  }   
  ```
- 4:处理消息
  ```
      func handleMsg(m interface{}, h interface{}) {
      	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
      }
      
      func init() {
      	handleMsg(&msg.C2S_UserHeartBeat{}, handleHeartBeat)
      	handleMsg(&msg.C2S_UserLogin{}, handleUserLogin)
      	handleMsg(&msg.C2S_Message{}, handleMessage)
      	handleMsg(&msg.C2S_UpdateWatchRoom{}, handleUpdateWatchRoom)
      }
      
      func handleHeartBeat(args []interface{}) {
      	//m := args[0].(*msg.C2S_UserHeartBeat)
      	a := args[1].(gate.Agent)
      	//log.Info("接受websocket心跳包:%v message:%v", a.RemoteAddr(), m)
      	info, ok := a.UserData().(UserInfo)
      	if ok && info.Uid != 0 {
      		info.Ptime = time.Now().Unix()
      		info.FlushOnlineTime(false)
      		a.SetUserData(info)
      	}
      	a.WriteMsg(&msg.C2S_UserHeartBeat{Ptime: time.Now().Unix()})
      }
      
      func handleUserLogin(args []interface{}) {
      	m := args[0].(*msg.C2S_UserLogin)
      	a := args[1].(gate.Agent)
      	t, err := utils.DecodeTokenByStr(m.Token)
      	if err != nil {
      		//log.Error("%v 解析token失败:%v,err:%v", m.LoginPlat, m.Token, err)
      		a.Close()
      		return
      	} else {
      		//log.Debug("%v 平台token解析:uid:%v,device:%v,roleid:%v", m.LoginPlat, t.Uid, t.UserData, t.RoleId)
      	}
      	a.SetUserData(UserInfo{Uid: t.Uid, Ptime: 0, LoginPlat: m.LoginPlat, LoginTime: time.Now().Unix(), LastAddTime: time.Now().Unix()})
      	addUserAgent(t.Uid, a, m.LoginPlat)
      
      	//登录返回
      	a.WriteMsg(&msg.S2C_UserLogin{Ret: 1, Msg: "登录成功"})
      }
  
  ```
- 5:启动服务
  ```
  func StartChatServer(end chan bool) bool {
  	if !conf.InitConf() {
  		return false
  	}
  
  	connsvr.RunInside(end,
  		game.Module,
  		gate.Module,
  	)
  	return true
  }  
  ```
  



