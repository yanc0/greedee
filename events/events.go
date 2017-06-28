package events

import (
	"time"
	"crypto/sha256"
	"fmt"
	"unicode/utf8"
	"errors"
)

type Event struct {
	SHA256Sum string `gorm:"column:sha256_sum;size:64;primary_key"`
	Name string `gorm:"index"`
	Type string `gorm:"index"`
	Description string
	TTL int64
	AuthUserSource string
	Timestamp time.Time
}

// Generate UUID (SHA256) with event infos
func (e *Event) Gen256Sum(){
	eventBytes := []byte(e.Name + e.Type + e.Description + string(e.TTL) + e.Timestamp.String())
	h := sha256.New()
	h.Write(eventBytes)
	e.SHA256Sum = fmt.Sprintf("%x", h.Sum(nil))
}

// Check if event is well formed
func (e *Event) Check() error {
	if utf8.RuneCountInString(e.SHA256Sum) != 64 {
		return errors.New("Bad SHA256 sum")
	} else if utf8.RuneCountInString(e.Name) <= 0 {
		return errors.New("name cannot be blank")
	}
	return nil
}