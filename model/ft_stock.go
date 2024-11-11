package model

type FT_Stock struct {
	ID        int    `json:"id" gorm:"primaryKey" gorm:"column:id"`
	StockID   string `json:"stockid" gorm:"column:stockid"`
	StockName string `json:"stockname" gorm:"column:stockname"`
}

func (FT_Stock) TableName() string {
	return "ft_stock"
}
