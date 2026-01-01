package domain

import (
	"golang.org/x/crypto/bcrypt"
)

// User はユーザーエンティティ
type User struct {
	BaseEntity
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"` // JSONには含めない
	Name     string `json:"name" db:"name"`
}

// SetPassword はパスワードをハッシュ化して設定
func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// CheckPassword はパスワードが正しいかチェック
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

