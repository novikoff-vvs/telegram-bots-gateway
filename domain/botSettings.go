package domain

import "gorm.io/gorm"

type BotSettings struct {
	gorm.Model
	ApiKey      string `json:"api_key"`
	AdminID     string `json:"admin_id"`
	ServiceUrl  string `json:"service_url"`
	ServiceType string `json:"service_type"`
}
