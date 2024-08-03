package httpclient

import (
	"fmt"
	"log"
	"time"

	tpprotocolsdkgo "github.com/ThingsPanel/tp-protocol-sdk-go"
	"github.com/ThingsPanel/tp-protocol-sdk-go/api"
	"github.com/spf13/viper"
)

var client *tpprotocolsdkgo.Client

func Init() {
	addr := viper.GetString("thingspanel.address")
	log.Println("创建http客户端:", addr)
	client = tpprotocolsdkgo.NewClient(addr)
	go ServiceHeartbeat()
}

func GetDeviceConfig(deviceNumber string) (*api.DeviceConfigResponse, error) {
	deviceConfigReq := api.DeviceConfigRequest{
		DeviceNumber: deviceNumber,
	}
	response, err := client.API.GetDeviceConfig(deviceConfigReq)
	if err != nil {
		errMsg := fmt.Sprintf("获取设备配置失败 (请求参数： %+v): %v", deviceConfigReq, err)
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return response, nil
}

func ServiceHeartbeat() {
	for {
		err := reportHeartbeat()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(50 * time.Second)
	}
}

// 这里需要改为自己的服务
func reportHeartbeat() error {
	sid := viper.GetString("server.identifier")
	serviceHeartbeatReq := api.HeartbeatRequest{
		ServiceIdentifier: sid,
	}
	response, err := client.API.Heartbeat(serviceHeartbeatReq)
	if err != nil {
		return fmt.Errorf("服务心跳上报失败 (请求参数：%+v): %v", serviceHeartbeatReq, err)
	}
	if response.Code != 200 {
		return fmt.Errorf("服务心跳上报失败 (请求参数：%+v): %v", serviceHeartbeatReq, response.Message)
	}
	return nil
}

// 获取服务接入点列表
func GetServiceAccessPointList() (*api.ServiceAccessListResponseData, error) {
	ServiceAccessListReq := api.ServiceAccessListRequest{
		ServiceIdentifier: viper.GetString("server.identifier"),
	}

	response, err := client.API.GetServiceAccessList(ServiceAccessListReq)
	if err != nil {
		return nil, fmt.Errorf("获取服务接入点列表失败: %v", err)
	}
	return response, nil
}

// 获取服务接入点
func GetServiceAccessPoint(serviceAccessPointID string) (*api.ServiceAccessResponseData, error) {
	ServiceAccessPointReq := api.ServiceAccessRequest{
		ServiceAccessID: serviceAccessPointID,
	}

	response, err := client.API.GetServiceAccess(ServiceAccessPointReq)
	if err != nil {
		return nil, fmt.Errorf("获取服务接入点失败: %v", err)
	}
	return response, nil
}
