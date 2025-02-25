package data

import "gorm.io/gorm"

type AuthData struct {
	AuthID      string `gorm:"unique;not null;type:uuid;uniqueIndex"`
	Username    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	PhoneNumber string `gorm:"unique;not null"`
}

type Preferences struct {
	FavoriteColor *string
}

type User struct {
	gorm.Model
	AuthData
	Preferences
}
