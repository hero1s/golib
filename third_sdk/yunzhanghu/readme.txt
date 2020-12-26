// @Title 云账户提现回调接口
// @Description 云账户提现回调接口
// @Success 0 {string} 状态码
// @router /yunzhanghu-callback [post]
func (b *PublicController) YunZhangHuCallBack() {
	repStr := "fail"
	defer func() {
		b.Ctx.Output.SetStatus(http.StatusOK)
		b.Ctx.ResponseWriter.Write([]byte(repStr))
		b.StopRun()
		return
	}()

	// 4、订单回调,响应云账户综合服务平台订单打款回调订单信息
	data := b.GetString("data")
	mess := b.GetString("mess")
	timestamp := b.GetString("timestamp")
	sign := b.GetString("sign")
	orderDetails, err := yunzhanghu.YunZhangHuClient.OrderCallBack(data, mess, timestamp, sign)
	if err != nil {
		log.Errorf("云账户回调参数解析错误:%v", err)
		notifyYunZhangHuError(err.Error())
		return
	}
	log.Infof("云账户回调结果:%v", orderDetails)

	if orderDetails.Status != "1" {
		log.Errorf(fmt.Sprintf("云账户提现回调状态未打款成功,trade status:%v", orderDetails.Status))
		notifyYunZhangHuError(orderDetails.StatusDetailMessage)
		return
	}
	//----------------------下面是提现成功的情况----------------------
	pay.PayTransferSuccess(orderDetails.OrderID, orderDetails.Ref)
	repStr = "success"
	return
}