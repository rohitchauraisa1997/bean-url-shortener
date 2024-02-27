package models

type Admin struct {
	ID       int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Username string `json:"username" gorm:"column:username;unique"`
	Password string `json:"password" gorm:"column:password"`
	Email    string `json:"email" gorm:"column:email;unique"`
}

func (Admin) TableName() string {
	return "Admin"
}
