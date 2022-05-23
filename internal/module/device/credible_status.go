package device

const (
	CredibleStatusDeviceNotClockIn = 0 // CredibleStatusDeviceNotPicketed 设备未纠察
	CredibleStatusDeviceFound      = 1 // CredibleStatusDeviceNotPicketed 设备纠查存在
	CredibleStatusDeviceNotFound   = 2 // CredibleStatusDevicePicketedNotFound 设备纠查不存在
	CredibleStatusDeviceError      = 3 //CredibleStatusDevicePicketedError 设备纠查异常
)

func GetCredibleStatus(exists bool) int {
	if exists {
		return CredibleStatusDeviceFound
	}

	return CredibleStatusDeviceNotFound
}
