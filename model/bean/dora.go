package bean

type Dora struct {
	ID   uint64 `gorm:"column:id" json:"id"`
	Mark string `gorm:"column:mark" json:"mark"`
}
