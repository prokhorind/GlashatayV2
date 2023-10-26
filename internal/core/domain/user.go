package domain

type User struct {
	ID     int64
	Active bool
}

func (m *User) TableName() string {
	return "users"
}
