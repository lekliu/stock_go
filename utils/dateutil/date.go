package dateutil

import (
	"fmt"
	"time"
)

func GetLastQt1() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	var lastQtDate string

	if month >= 1 && month <= 3 {
		lastQtDate = fmt.Sprintf("%d-12-31", year-1)
	} else if month >= 4 && month <= 6 {
		lastQtDate = fmt.Sprintf("%d-03-31", year)
	} else if month >= 7 && month <= 9 {
		lastQtDate = fmt.Sprintf("%d-06-30", year)
	} else if month >= 10 && month <= 12 {
		lastQtDate = fmt.Sprintf("%d-09-30", year)
	}

	return lastQtDate
}

func GetLastQt4() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	lastYear := year - 1
	var lastQtDates string

	if month >= 1 && month <= 3 {
		lastQtDates = fmt.Sprintf("%d-12-31%%2C%d-09-30%%2C%d-06-30%%2C%d-03-31", lastYear, lastYear, lastYear, lastYear)
	} else if month >= 4 && month <= 6 {
		lastQtDates = fmt.Sprintf("%d-03-31%%2C%d-12-31%%2C%d-09-30%%2C%d-06-30", year, lastYear, lastYear, lastYear)
	} else if month >= 7 && month <= 9 {
		lastQtDates = fmt.Sprintf("%d-06-30%%2C%d-03-31%%2C%d-12-31%%2C%d-09-30", year, year, lastYear, lastYear)
	} else if month >= 10 && month <= 12 {
		lastQtDates = fmt.Sprintf("%d-09-30%%2C%d-06-30%%2C%d-03-31%%2C%d-12-31", year, year, year, lastYear)
	}

	return lastQtDates
}

func GetLastQt5() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	lastYear := year - 1
	last2Year := year - 2

	// 构建日期字符串，使用 %27 表示单引号，%2C 表示逗号
	Qt3 := fmt.Sprintf("%%27%d-09-30%%27", year)
	Qt2 := fmt.Sprintf("%%27%d-06-30%%27", year)
	Qt1 := fmt.Sprintf("%%27%d-03-31%%27", year)
	Qt14 := fmt.Sprintf("%%27%d-12-31%%27", lastYear)
	Qt13 := fmt.Sprintf("%%27%d-09-30%%27", lastYear)
	Qt12 := fmt.Sprintf("%%27%d-06-30%%27", lastYear)
	Qt11 := fmt.Sprintf("%%27%d-03-31%%27", lastYear)
	Qt24 := fmt.Sprintf("%%27%d-12-31%%27", last2Year)

	var lastQtDates string

	// 根据当前月份返回对应的最近五个季度的日期
	if month >= 1 && month <= 3 {
		lastQtDates = fmt.Sprintf("%s%%2C%s%%2C%s%%2C%s%%2C%s", Qt14, Qt13, Qt12, Qt11, Qt24)
	} else if month >= 4 && month <= 6 {
		lastQtDates = fmt.Sprintf("%s%%2C%s%%2C%s%%2C%s%%2C%s", Qt1, Qt14, Qt13, Qt12, Qt11)
	} else if month >= 7 && month <= 9 {
		lastQtDates = fmt.Sprintf("%s%%2C%s%%2C%s%%2C%s%%2C%s", Qt2, Qt1, Qt14, Qt13, Qt12)
	} else if month >= 10 && month <= 12 {
		lastQtDates = fmt.Sprintf("%s%%2C%s%%2C%s%%2C%s%%2C%s", Qt3, Qt2, Qt1, Qt14, Qt13)
	}

	return lastQtDates
}
