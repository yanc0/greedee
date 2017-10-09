package events

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

type Event struct {
	ID             string    `gorm:"size:64;primary_key";json:"id"`
	Name           string    `gorm:"index"`
	TTL            int64     `gorm:"index"`
	AuthUserSource string    `gorm:"index"`
	CreatedAt      time.Time `gorm:"index"`
	ExpiresAt      time.Time `gorm:"index"`
	Description    string
	Source         string
	Tags           string
	Processed      bool
	Expired        bool
	Status         uint8
}

// Generate UUID (SHA256) with event infos
func (e *Event) Gen256Sum() {
	eventBytes := []byte(e.Name + e.Source + e.Description + string(e.TTL) + e.CreatedAt.String() + string(e.Status))
	h := sha256.New()
	h.Write(eventBytes)
	e.ID = fmt.Sprintf("%x", h.Sum(nil))
}

// Add TTL to timestamp and calculate event expiration date
func (e *Event) GenExpiredAt() {
	if e.IsExpirable() {
		e.ExpiresAt = e.CreatedAt.Add(time.Duration(e.TTL) * time.Second)
	}
}

// IsExpirable return true if event has TTL, success status and CreatedAt
func (e Event) IsExpirable() bool {
	return e.TTL > 0 && e.Status == 0 && !e.CreatedAt.IsZero()
}

// Check if event is well formed
func (e *Event) Check() error {
	if utf8.RuneCountInString(e.ID) != 64 {
		return errors.New("Bad SHA256 sum")
	} else if utf8.RuneCountInString(e.Name) <= 0 {
		return errors.New("name cannot be blank")
	}
	return nil
}

// Fail sets event to fail status
func (e *Event) Fail() {
	e.Status = 1
}

// Success sets event to Success status
func (e *Event) Success() {
	e.Status = 0
}

// isSuccess return true if event is in success status
func (e Event) isSuccess() bool {
	return e.Status == 0
}
