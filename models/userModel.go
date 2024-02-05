package models

import (
	"gorm.io/gorm"
	"github.com/gorilla/websocket"
)

type User struct {
	gorm.Model
	Name          string	`json:"name"`
	Email         string	`json:"email"`
	PhoneNumber string		`json:"phone_number"`                                                                                                                            
	Password	  string	`json:"password"`
	Connection    *websocket.Conn `gorm:"-"`
	CurrentAudio  []byte          `gorm:"-"`
	ReceivedAudio chan []byte     `gorm:"-"`
}

type Message struct {
	gorm.Model
	SenderID    uint
	RecipientID uint
	AudioData   []byte
}

type Contact struct {
	gorm.Model
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	PhoneNumber  string `json:"phone_number"`
}