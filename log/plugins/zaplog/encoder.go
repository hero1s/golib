package zaplog

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"strings"
)

type customEncoder struct {
	zapcore.Encoder
	name string
}

// 重写 EncodeEntry 方法
func (c *customEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := c.Encoder.EncodeEntry(ent, fields)
	if err != nil {
		return nil, err
	}
	// 获取时间戳
	timestamp := ent.Time.Local().Format("2006-01-02 15:04:05.000")
	// 获取日志级别
	level := ent.Level.CapitalString()
	// 获取日志消息
	message := string(buf.Bytes())
	// 去掉原始日志中的时间戳和日志级别
	// 假设原始日志格式为 "2025-03-24T11:00:35.968+0800	INFO	消息"
	parts := strings.SplitN(message, "\t", 3)
	if len(parts) >= 3 {
		message = parts[2] // 取第三部分作为消息
	} else {
		message = parts[len(parts)-1] // 如果格式不符合预期，取最后一部分作为消息
	}
	// 重新格式化日志
	formatted := fmt.Sprintf("[%s][%s][%s] %s", timestamp, c.name, level, message)
	buf.Reset()
	buf.AppendString(formatted)
	return buf, nil
}

// 重写 Clone 方法
func (c *customEncoder) Clone() zapcore.Encoder {
	return &customEncoder{Encoder: c.Encoder.Clone(), name: c.name}
}
