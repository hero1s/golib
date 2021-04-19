package tcp

import (
	"fmt"
	"github.com/hero1s/golib/log"
	"net"
	"time"
)

//read data from peer and decrypt data, and return data
func ReadPacket(conn net.Conn) (Frame, error) {
	return decodeFrame(conn)
}

//just send data to peer
func WritePacket(conn net.Conn, cmd uint16, encType byte, data []byte) error {
	buf, err := buildFrame(cmd, encType, data)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf)
	return err
}

// 发送消息到指定ip,port的服务器
func SendMsgAndRecvReply(ip string, port int, cmd uint16, msg []byte,encType byte) (bool, Frame) {
	var reply Frame
	tcpIp := fmt.Sprintf("%s:%d", ip, port)
	log.Debugf("sendmsg to:%v",tcpIp)
	conn, err := net.Dial("tcp", tcpIp)
	if err != nil {
		log.Errorf("connect tcp fail:%v,%v", tcpIp, err)
		return false, reply
	}
	defer func() {
		conn.Close()
	}()
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	err = WritePacket(conn, cmd, encType, msg)
	if err != nil {
		log.Debugf("send tcp packet fail:%v,%v", tcpIp, err)
		return false, reply
	}
	reply, err = ReadPacket(conn)
	if err != nil {
		log.Debugf("read tcp packet fail:%v,%v", tcpIp, err)
		return false, reply
	}
	log.Debugf("connect tcp reply:%v,%v", tcpIp, string(reply.Data))
	return true, reply
}

