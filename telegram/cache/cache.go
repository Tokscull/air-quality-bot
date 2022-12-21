package cache

type Key int64

type Type int

const (
	Message Type = iota
	Location
)

type HandlerKey int

const (
	AirQuality HandlerKey = iota
	Notifications
)

type Data struct {
	Type       Type
	HandlerKey HandlerKey
	ChildInfo  interface{}
}
