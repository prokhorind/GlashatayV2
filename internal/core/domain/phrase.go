package domain

type Phrase struct {
	Id     int64
	UserID int
	Phrase string
	Type   string
	Public bool
}

func (m *Phrase) TableName() string {
	return "phrases"
}
