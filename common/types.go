// File: common/types.go

package common

type InviteTokenStatus string

const (
	StatusUsed      InviteTokenStatus = "Used"
	StatusAvailable InviteTokenStatus = "Available"
	StatusExpired   InviteTokenStatus = "Expired"
)

func (InviteTokenStatus) Values() (statuses []string) {
	return []string{string(StatusUsed), string(StatusAvailable), string(StatusExpired)}
}
