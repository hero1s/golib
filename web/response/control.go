package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hero1s/golib/log"
	"net/http"
)

// 结构体
type Control struct {
}

// 路由接口
type Router interface {
	InitRoutes(r *gin.Engine)
}

// 响应成功
func (c *Control) Success(context *gin.Context, message string) {
	context.JSON(http.StatusOK, JsonObject{
		Code: "0",
		Msg:  message,
		Data: nil,
	})
	context.Abort()
	return
}

// 返回状态和结果
func (c *Control) SuccessContent(context *gin.Context, message string, data interface{}) {
	context.JSON(http.StatusOK, JsonObject{
		Code: "0",
		Msg:  message,
		Data: data,
	})
	context.Abort()
	return
}

// 无权限的访问
func (c *Control) RefusedError(context *gin.Context) {
	context.JSON(http.StatusUnauthorized, JsonObject{
		Code: "401",
		Msg:  "Unauthorized",
		Data: nil,
	})
	context.Abort()
	return
}

// 无权限的访问
func (c *Control) BindingError(context *gin.Context) {
	context.JSON(http.StatusOK, JsonObject{
		Code: "404",
		Msg:  "Binding Error",
		Data: nil,
	})
	context.Abort()
	return
}

// 服务器内部错误
func (c *Control) InternalError(context *gin.Context, err string) {
	context.JSON(http.StatusInternalServerError, JsonObject{
		Code: "500",
		Msg:  err,
		Data: nil,
	})
	context.Abort()
	return
}

// 通用的统一返回结果
func (c *Control) ReturnResult(context *gin.Context, code string, message string, data interface{}) {
	context.JSON(http.StatusOK, JsonObject{
		Code: code,
		Msg:  message,
		Data: data,
	})
	context.Abort()
	return
}

func (c *Control) File(context *gin.Context, filePath, fileName string) {
	context.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	context.File(filePath)
}

func (c *Control) TryBindParam(context *gin.Context, param any) bool {
	if err := context.Bind(param); err != nil {
		log.Errorf("try bind param error:%v", err)
		c.BindingError(context)
		return false
	}
	return true
}
