package client

type ActionButtonType byte

const (
	UnspecifiedAction ActionButtonType = iota
	View
	HTTP
	Broadcast
)

type ActionButton interface {
	actionType() ActionButtonType
}
