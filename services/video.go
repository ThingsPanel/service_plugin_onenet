package services

import (
	"encoding/json"
	"plugin_onenet/mqtt"

	"github.com/ThingsPanel/tp-protocol-sdk-go/api"
	"github.com/sirupsen/logrus"
)

func Start() {
	logrus.Info("start video service")
	Deal()

}

// 获取服务接入点列表
func Deal() {
	logrus.Info("get service access point list")
	pointList, err := GetServiceAccessPointList()
	logrus.Info(pointList)
	if err != nil {
		logrus.Error(err)
	}
	if pointList.Code != 200 {
		logrus.Error(pointList.Message)
		return
	}
	if len(pointList.Data) == 0 {
		logrus.Error("no service access point")
		return
	}
	for _, point := range pointList.Data {
		logrus.Info(point)
		DealPoint(point)
	}
}

// 处理接入点
func DealPoint(point api.ServiceAccess) {
	voucherStr := point.Voucher
	voucher := Voucher{}
	// 校验voucherStr是json字符串
	err := json.Unmarshal([]byte(voucherStr), &voucher)
	if err != nil {
		logrus.Error(err)
		return
	}
	devices := point.Devices
	if len(devices) == 0 {
		logrus.Error("no devices")
		return
	}
	for _, device := range devices {
		// 组装待发送数据{key:value}
		data := make(map[string]interface{})
		data["url"] = voucher.Url
		// 发布属性消息
		err := mqtt.PublishAttributes(device.ID, data)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
}
