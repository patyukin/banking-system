package model

import "time"

type SmsHistory struct {
	ID        uint64    `json:"id"`
	Phone     string    `json:"phone"`
	Message   string    `json:"message"`
	SmsID     uint64    `json:"sms_id"`
	Status    string    `json:"status"`
	Report    string    `json:"report"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmailHistory struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	Report    string    `json:"report"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
