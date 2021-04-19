package tcp

import (
	"io"
	"net"
)

const (
	HeaderLen   = 5
	MaxDataLen  = 1024 // the max frame data length
	MaxFrameLen = HeaderLen + MaxDataLen
)

var (
	secretKey = []byte("123456")
)

type Frame struct {
	Length  uint16 //表示data的长度
	Cmd     uint16
	EncType byte
	Data    []byte
}

func buildFrame(cmd uint16, encType byte, data []byte) (buf []byte, err error) {
	var encryptData []byte
	var length int
	if encType != 0 {
		encryptData, _ = Encrypt(data, secretKey)
		length = len(encryptData)
	} else {
		length = len(data)
		encryptData = data
	}

	buf = make([]byte, length+HeaderLen)
	s := NewLEStream(buf)
	s.WriteUint16(uint16(length))
	s.WriteUint16(cmd)
	s.WriteByte(encType)
	err = s.WriteBuff(encryptData)
	return
}

func decodeFrame(conn net.Conn) (f Frame, err error) {
	headerBuf := make([]byte, HeaderLen)
	_, err = io.ReadFull(conn, headerBuf)
	if err != nil {
		return
	}
	s := NewLEStream(headerBuf)
	f.Length, _ = s.ReadUint16()
	f.Cmd, _ = s.ReadUint16()
	f.EncType, _ = s.ReadByte()

	data := make([]byte, f.Length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return
	}
	decryptData := data
	if f.EncType != 0 {
		decryptData, _ = Decrypt(data, secretKey)
	}
	ds := NewLEStream(decryptData)
	f.Data, err = ds.ReadBuff(ds.Size())
	return
}
