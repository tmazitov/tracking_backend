package bl

type UserRole int

const (
	Base   UserRole = 1
	Manage UserRole = 2
	Admin  UserRole = 3
)

type GetUser struct {
	TelNumber string   `json:"telNumber"`
	ShortName string   `json:"shotName"`
	RoleID    UserRole `json:"roleId"`
}
