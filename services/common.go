package services

import (
	"fmt"
	httpclient "plugin_onenet/http_client"

	"github.com/ThingsPanel/tp-protocol-sdk-go/api"
)

type Voucher struct {
	Url string `json:"url"`
}

// 凭证信息组装
func AssembleVoucher(deviceSecret string) (voucher string) {
	return fmt.Sprintf(`{"UID":"%s"}`, deviceSecret)
}

// 获取服务接入点列表并处理信息
func GetServiceAccessPointList() (*api.ServiceAccessListResponseData, error) {
	// 读取服务接入点列表
	return httpclient.GetServiceAccessPointList()
}
