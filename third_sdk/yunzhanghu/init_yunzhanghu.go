package yunzhanghu

import (
	"github.com/hero1s/golib/log"
	"github.com/hero1s/golib/third_sdk/pay/gopay"
	sdk2 "github.com/hero1s/golib/third_sdk/yunzhanghu/sdk"
)

// 接口配置信息
var (
	BrokerID  = "yiyun73"
	DealerID  = "12345678"
	AppKey    = "xxxxx"
	Des3Key   = "xxxxx"
	Gateway   = "https://api-jiesuan.yunzhanghu.com"
	NotifyURL = ""
)

var YunZhangHuClient *sdk2.Client

func InitYunZhangHu() {
	// 初始化SDK客户端
	log.Infof("初始化云账户客户端:%v \n %v \n %v \n %v \n %v \n %v \n", BrokerID, DealerID, AppKey, Des3Key, Gateway,NotifyURL)
	YunZhangHuClient = sdk2.New(BrokerID, DealerID, Gateway, AppKey, Des3Key)
}

// 1、银行卡下单,响应云账户综合服务平台订单流水号
func CreateBankOrder(cardNo, orderId, realName, idCard string, moneyFee int64, remark, phone string) error {
	// 1、银行卡下单,响应云账户综合服务平台订单流水号
	bankOrderParam := &sdk2.BankOrderParam{
		CardNo:  cardNo,
		PhoneNo: phone,
	}
	bankOrderParam.OrderID = orderId
	bankOrderParam.RealName = realName
	bankOrderParam.IDCard = idCard
	bankOrderParam.Pay = gopay.AliMoneyFeeToString(float64(moneyFee))
	bankOrderParam.PayRemark = remark
	bankOrderParam.NotifyURL = NotifyURL
	ref, err := YunZhangHuClient.CreateBankOrder(bankOrderParam)
	if err != nil {
		return err
	}
	log.Debugf("银行下单信息:%v", ref)
	return nil
}
