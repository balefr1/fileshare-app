//Models/UserModel.go
package Models

type User struct {
	Id       uint   `json:"id"`
	Name     string `json:"name" gorm:"not null"`
	LastName string `json:"lastname" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Username string `json:"username" gorm:"unique;not null"`
}

func (b *User) TableName() string {
	return "user"
}
