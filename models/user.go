package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:char(36);primary_key" json:"id,omitempty"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex:idx_users_username,LENGTH(100);not null" json:"username,omitempty"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex:idx_users_email,LENGTH(255);not null" json:"email,omitempty"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password,omitempty"`
	FullName  string    `gorm:"type:varchar(255);not null" json:"fullName,omitempty"`
	Role      string    `gorm:"type:varchar(50);default:'user'" json:"role,omitempty"`
	CreatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'; ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type CreateUserSchema struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"fullName" validate:"required,min=3,max=255"`
	Role     string `json:"role" validate:"omitempty,oneof=USER ADMIN"`
}

type UpdateUserSchema struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=100"`
	Email    *string `json:"email" validate:"omitempty,email"`
	FullName *string `json:"fullName" validate:"omitempty,min=3,max=255"`
	Role     *string `json:"role" validate:"omitempty,oneof=USER ADMIN"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New().String()
	return nil
}
