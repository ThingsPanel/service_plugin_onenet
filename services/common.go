package services

import (
	"fmt"
	httpclient "plugin_onenet/http_client"

	"github.com/ThingsPanel/tp-protocol-sdk-go/api"
	"github.com/sirupsen/logrus"
)

type Voucher struct {
	Url string `json:"url"`
}

// 认证设备并获取设备信息
func AuthDevice(deviceSecret string) (deviceInfo *api.DeviceConfigResponse, err error) {
	voucher := AssembleVoucher(deviceSecret)
	// 读取设备信息
	deviceInfo, err = httpclient.GetDeviceConfig(voucher, "")
	if err != nil {
		// 获取设备信息失败，请检查连接包是否正确
		logrus.Error(err)
		return
	}
	if deviceInfo.Code != 200 {
		err = fmt.Errorf("device auth failed, code: %d, message: %s", deviceInfo.Code, deviceInfo.Message)
		logrus.Error(err)
	}
	return
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
