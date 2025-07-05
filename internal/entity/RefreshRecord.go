package entity

import (
	"github.com/google/uuid"
)

type RefreshRecord struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	RefreshHash []byte    `gorm:"unique;column:refresh_hash"`
	UserAgent   string    `gorm:"column:user_agent"`
	ClientIp    string    `gorm:"column:client_ip"`
	AccessJTI   string    `gorm:"column:access_jti"`
}
