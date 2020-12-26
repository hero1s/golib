package publics

import (
	"fmt"
	"github.com/hero1s/gotools/gopay"
	"git.moumentei.com/plat_go/golib/log"
	gopay2 "github.com/iGoogle-ink/gopay"
	"home/models/pay"
	"home/shared"
	"home/shared/db_table"
	"net/http"
	"time"
)

// @Title 微信支付回调接口
// @Description 微信支付回调接口
// @Success 0 {string} 状态码
// @router /wechat-callback [post]
func (b *PublicController) WeChatCallBack() {
	//==异步通知，返回给微信平台的信息==
	rsp := new(gopay2.WeChatNotifyResponse) //回复微信的数据
	rsp.ReturnCode = gopay2.SUCCESS
	rsp.ReturnMsg = gopay2.OK
	defer func() {
		b.Ctx.Output.SetStatus(http.StatusOK)
		b.Ctx.ResponseWriter.Write([]byte(rsp.ToXmlString()))
		b.StopRun()
		return
	}()
	notifyReq, err := gopay2.ParseWeChatNotifyResult(b.Ctx.Request)
	if err != nil {
		log.Errorf("微信回调参数解析失败:%#v", err)
		rsp.ReturnCode = gopay2.FAIL
		rsp.ReturnMsg = gopay2.FAIL
		return
	}
	order, err1 := pay.PayOrderByPayId(notifyReq.OutTradeNo)
	if err1 != nil {
		log.Errorf("微信验证签名前获取订单失败:%#v", notifyReq)
		rsp.ReturnCode = gopay2.FAIL
		rsp.ReturnMsg = gopay2.FAIL
		return
	}
	//验签操作
	tradeType := gopay2.TradeType_App
	if order.PayType == shared.WeChatJSAPI {
		tradeType = gopay2.TradeType_JsApi
	}
	ok, err := pay.GetWeChatClient(order.PayNo, order.PayType, order.PayChannel).VerifyWeChatSign(notifyReq, tradeType)
	if !ok || err != nil {
		log.Errorf("微信支付回调验证失败:%v", err)
		rsp.ReturnCode = gopay2.FAIL
		rsp.ReturnMsg = gopay2.FAIL
		return
	}

	log.Debugf(fmt.Sprintf("微信回调结果:%#v\n", notifyReq))
	if notifyReq.ResultCode != gopay2.SUCCESS { //业务失败
		updateData := map[string]interface{}{
			"status":         shared.OrderPayError,
			"transaction_id": notifyReq.TransactionId,
			"pay_time":       time.Now().Unix(),
		}
		_, err := db_table.TUserPay.UpdateTableRecord(updateData, fmt.Sprintf(`uid='%v'`, notifyReq.OutTradeNo))
		if err != nil {
			log.Errorf(fmt.Sprintf("更新平台订单:%v,微信订单号:%v失败:%v", notifyReq.OutTradeNo, notifyReq.TransactionId, err))
		}
		return
	}
	//----------------------下面是支付成功的情况----------------------
	bRet := pay.PaySuccess(notifyReq.OutTradeNo, notifyReq.TransactionId, int64(notifyReq.TotalFee), notifyReq.Openid,notifyReq.DeviceInfo,false)
	if !bRet {
		log.Errorf("微信回调支付处理失败:%v,%v", notifyReq.OutTradeNo, notifyReq.TransactionId)
		rsp.ReturnCode = gopay2.FAIL
		rsp.ReturnMsg = gopay2.FAIL
		return
	}
	return
}

// @Title 支付宝支付回调接口
// @Description 支付宝支付回调接口
// @Success 0 {string} 状态码
// @router /ali-callback [post]
func (b *PublicController) AliCallBack() {
	repStr := "fail"
	defer func() {
		b.Ctx.Output.SetStatus(http.StatusOK)
		b.Ctx.ResponseWriter.Write([]byte(repStr))
		b.StopRun()
		return
	}()
	notifyReq, err := gopay2.ParseAliPayNotifyResult(b.Ctx.Request)
	if err != nil {
		log.Errorf("支付宝回调解析参数失败%#v:%#v", notifyReq, err)
		return
	}
	order, err1 := pay.PayOrderByPayId(notifyReq.OutTradeNo)
	if err1 != nil {
		log.Errorf("支付宝验证签名前获取订单失败:%#v", notifyReq)
		repStr = "success"
		return
	}
	//验签操作
	ok, err := pay.GetAliClient(order.PayNo, order.PayType, order.PayChannel).VerifyAliPaySign(notifyReq)
	if !ok {
		log.Errorf("支付宝回调参数:%#v 验证失败:%v", notifyReq, err)
		return
	}
	if notifyReq.TradeStatus != "TRADE_SUCCESS" {
		log.Errorf(fmt.Sprintf("支付宝支付app回调通信失败,trade status:%v", notifyReq.TradeStatus))
		return
	}
	//pay不会为nil
	log.Infof("支付宝支付账户:%v,%v", notifyReq.BuyerId, notifyReq.BuyerLogonId)
	//----------------------下面是支付成功的情况----------------------
	totalFee := gopay.AliStringToMoneyFee(notifyReq.TotalAmount)
	bRet := pay.PaySuccess(notifyReq.OutTradeNo, notifyReq.TradeNo, totalFee, notifyReq.BuyerId, notifyReq.BuyerLogonId,true)
	if !bRet {
		log.Errorf("支付宝回调支付处理失败,请手动发货:%v,%v", notifyReq.OutTradeNo, notifyReq.TradeNo)
		return
	}
	repStr = "success"
	return
}
