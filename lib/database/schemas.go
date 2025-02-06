package database

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя в системе
type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"size:50" json:"name"`
	Email    *string   `gorm:"type:text;unique" json:"email"`
	Login    string    `gorm:"size:50;unique;not null" json:"login"`
	Password []byte    `gorm:"type:text;not null" json:"-"` // Пароль будет храниться в виде хеша
	Role     Role      `gorm:"type:user_role;not null" json:"role"`
}

// Lang представляет услугу
type Lang struct {
	ID               uint    `gorm:"primaryKey;autoIncrement"`
	Name             string  `gorm:"size:50;not null"`
	ShortDescription string  `gorm:"size:255;not null"`
	Description      string  `gorm:"type:text;not null"`
	ImgLink          *string `gorm:"size:255"`
	Author           *string `gorm:"size:50"`
	Year             *string `gorm:"size:4"`
	Version          *string `gorm:"size:50"`
	List             JSONB  `gorm:"type:jsonb"`
	Status           bool    `gorm:"default:true;not null"`
}

// Project представляет заявку
// @Description Project represents a project in the system
type Project struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	UserID           uuid.UUID `gorm:"type:uuid;not null"`
	User             *User     `gorm:"foreignKey:UserID" json:"User,omitempty"`
	CreationTime     time.Time `gorm:"default:current_timestamp"`
	FormationTime    *time.Time
	CompletionTime   *time.Time
	Status           Status     `gorm:"type:project_status;not null"`
	ModeratorID      *uuid.UUID `gorm:"type:uuid"`
	Moderator        *User      `gorm:"foreignKey:ModeratorID" json:"Moderator,omitempty"`
	ModeratorComment *string    `gorm:"type:text"`
}

// File представляет файл
type File struct {
	ID        uint    `gorm:"primaryKey"`
	LangID    uint    `gorm:"not null"`
	Lang      *Lang   `gorm:"foreignKey:LangID" json:"Lang,omitempty"`
	ProjectID uint    `gorm:"not null"`
	Code      string  `gorm:"type:text"`
	FileName  *string `gorm:"size:255"`
	FileSize  int64   `gorm:"default:0"`
	Comment   *string `gorm:"type:text"`
	AutoCheck *int
}
