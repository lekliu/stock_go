package api

import (
	"gorm.io/gorm"
)

//获取最新的股票列表

type Hdl_fetch struct {
	DB *gorm.DB
}

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
