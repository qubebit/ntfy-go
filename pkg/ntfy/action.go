package ntfy

const (
	UnspecifiedAction ActionButtonType = iota
	View
	HTTP
	Broadcast
)

type (
	ActionButtonType byte

	ActionButton interface {
		actionType() ActionButtonType
	}
)
