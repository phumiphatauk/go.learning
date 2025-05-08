package models

import (
	"time"
)

type Logger struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"` // Primary key
	Time         time.Time `gorm:"not null" json:"time"`
	RemoteIP     string    `gorm:"size:45;not null" json:"remote_ip"` // Supports IPv6
	Host         string    `gorm:"size:255;not null" json:"host"`
	Method       string    `gorm:"size:10;not null" json:"method"` // GET, POST, etc.
	URI          string    `gorm:"size:2048;not null" json:"uri"`  // Supports long URIs
	UserAgent    string    `gorm:"size:255" json:"user_agent"`
	Status       int       `gorm:"not null" json:"status"`
	Error        string    `gorm:"size:255" json:"error"`
	Latency      int64     `gorm:"not null" json:"latency"` // Microseconds
	LatencyHuman string    `gorm:"size:50" json:"latency_human"`
	BytesIn      int64     `gorm:"not null" json:"bytes_in"`
	BytesOut     int64     `gorm:"not null" json:"bytes_out"`
}
