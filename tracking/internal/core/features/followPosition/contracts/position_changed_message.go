package contracts

type Position struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type PositionChangedMessage struct {
	MessageId string   `json:"message_id"`
	DeviceId  string   `json:"device_id"`
	Position  Position `json:"position"`
	Timestamp int64    `json:"timestamp"`
}