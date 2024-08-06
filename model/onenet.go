package model

type DeviceItem struct {
	DeviceNumber string `json:"device_number"`
	DeviceName   string `json:"device_name"`
	Description  string `json:"description"`
}

// 事件/命令
type EventInfo struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}
