package models

import (
	"time"

	"gorm.io/datatypes"
)

type Client struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName string         `json:"first_name" gorm:"type:varchar(100);not null;index"`
	LastName  string         `json:"last_name" gorm:"type:varchar(100);not null;index"`
	BirthDate time.Time      `json:"birth_date" gorm:"type:date;not null;index"`
	CoreData  datatypes.JSON `json:"core_data" gorm:"type:jsonb"`
	Features  datatypes.JSON `json:"features" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Client) TableName() string {
	return "clients"
}

func (c *Client) IsValid() bool {
	return c.FirstName != "" && c.LastName != "" && !c.BirthDate.IsZero()
}
