package apple_pay

import (
	"encoding/json"
)

const (
	ProdVerifyUrl    = "https://buy.itunes.apple.com/verifyReceipt"
	SandBoxVerifyUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
)

//ios7以下版本返回的数据结构
type Response7 struct {
	Status  int64 `json:"status"`
	Receipt struct {
		OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
		PurchaseDateMs          string `json:"purchase_date_ms"`
		UniqueIdentifier        string `json:"unique_identifier"`
		OriginalTransactionId   string `json:"original_transaction_id"`
		Bvrs                    string `json:"bvrs"`
		TransactionId           string `json:"transaction_id"`
		Quantity                string `json:"quantity"`
		UniqueVendorIdentifier  string `json:"unique_vendor_identifier"`
		ItemId                  string `json:"item_id"`
		ProductId               string `json:"product_id"`
		PurchaseDate            string `json:"purchase_date"`
		OriginalPurchaseDate    string `json:"original_purchase_date"`
		OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
		PurchaseDatePst         string `json:"purchase_date_pst"`
		Bid                     string `json:"bid"`
	} `json:"receipt"`
}

//ios7及以上的版本返回的数据结构
type Response struct {
	Status      int64  `json:"status"`
	Environment string `json:"environment"`
	Receipt     struct {
		ReceiptType               string `json:"receipt_type"`
		AdamId                    int64  `json:"adam_id"`
		AppItemId                 int64  `json:"app_item_id"`
		BundleId                  string `json:"bundle_id"`
		ApplicationVersion        string `json:"application_version"`
		DownloadId                int64  `json:"download_id"`
		VersionExternalIdentifier int64  `json:"version_external_identifier"`

		ReceiptCreationDate    string `json:"receipt_creation_date"`
		ReceiptCreationDateMs  string `json:"receipt_creation_date_ms"`
		ReceiptCreationDatePst string `json:"receipt_creation_date_pst"`

		RequestDate    string `json:"request_date"`
		RequestDateMs  string `json:"request_date_ms"`
		RequestDatePst string `json:"request_date_pst"`

		OriginalPurchaseDate       string `json:"original_purchase_date"`
		OriginalPurchaseDateMs     string `json:"original_purchase_date_ms"`
		OriginalPurchaseDatePst    string `json:"original_purchase_date_pst"`
		OriginalApplicationVersion string `json:"original_application_version"`

		InApp []struct {
			Quantity              string `json:"quantity" desc:"商品数量"`
			ProductId             string `json:"product_id" desc:"商品id"`
			TransactionId         string `json:"transaction_id" desc:"交易id"`
			OriginalTransactionId string `json:"original_transaction_id"`

			PurchaseDate    string `json:"purchase_date" desc:"交易日期"`
			PurchaseDateMs  string `json:"purchase_date_ms"`
			PurchaseDatePst string `json:"purchase_date_pst"`

			OriginalPurchaseDate    string `json:"original_purchase_date"`
			OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
			OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`

			IsTrialPeriod string `json:"is_trial_period"`
		} `json:"in_app"`
	}
}

//ios7及以下的版本返回的数据结构
type R struct {
	Receipt struct {
		OriginalPurchaseDatePst   string `json:"original_purchase_date_pst"`
		PurchaseDateMs            string `json:"purchase_date_ms"`
		UniqueIdentifier          string `json:"unique_identifier"`
		OriginalTransactionId     string `json:"original_transaction_id"`
		Bvrs                      string `json:"bvrs"`
		TransactionId             string `json:"transaction_id"`
		Quantity                  string `json:"quantity"`
		UniqueVendorIdentifier    string `json:"unique_vendor_identifier"`
		ItemId                    string `json:"item_id"`
		VersionExternalIdentifier string `json:"version_external_identifier"`
		Bid                       string `json:"bid"`
		IsInIntroOfferPeriod      string `json:"is_in_intro_offer_period"`
		ProductId                 string `json:"product_id"`
		PurchaseDate              string `json:"purchase_date"`
		IsTrialPeriod             string `json:"is_trial_period"`

		PurchaseDatePst        string `json:"purchase_date_pst"`
		OriginalPurchaseDate   string `json:"original_purchase_date"`
		OriginalPurchaseDateMs string `json:"original_purchase_date_ms"`
		Status                 int64  `json:"status"`
	} `json:"receipt"`
	Status int64 `json:"status"`
}

//检验订单是否正确
func Verify(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}
func Verify7(data []byte) (R, error) {
	var r R
	err := json.Unmarshal(data, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}
