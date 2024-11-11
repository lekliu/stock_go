package stock

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

/*
CREATE TABLE `ft_income` (
  `id` int NOT NULL AUTO_INCREMENT,
  `stockid` varchar(8) NOT NULL COMMENT '代码',
  `stockname` varchar(15) NOT NULL COMMENT '名称',
  `reportdate` date DEFAULT NULL COMMENT '报告日期',
  `noticedate` date DEFAULT NULL COMMENT '公告日期',
  `netprofit` decimal(18,3) DEFAULT '0.000' COMMENT '净利润',
  `netprofitratio` decimal(20,9) NOT NULL DEFAULT '0.000000000' COMMENT '净利润同比',
  `totaloperateincome` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '营业总收入',
  `totaloperatecost` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '营业总支出',
  `operatecost` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '营业支出',
  `saleexpense` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '销售费用',
  `manageexpense` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '管理费用',
  `financeexpense` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '财务费用',
  `operateprofit` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '营业利润',
  `totalprofit` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '利润总额',
  `totaloperateincomeratio` decimal(20,9) NOT NULL DEFAULT '0.000000000' COMMENT '营业总收入同比',
  `upttime` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `stockid_reportdate` (`stockid`,`reportdate`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb3;
*/

type ResponseData struct {
	Result struct {
		Data []struct {
			StockID                 string  `json:"SECURITY_CODE"`
			StockName               string  `json:"SECURITY_NAME_ABBR"`
			Reportdate              string  `json:"REPORT_DATE"`
			Noticedate              string  `json:"NOTICE_DATE"`
			Netprofit               float32 `json:"PARENT_NETPROFIT"`
			Netprofitratio          float32 `json:"PARENT_NETPROFIT_RATIO"`
			Totaloperateincome      float32 `json:"TOTAL_OPERATE_INCOME"`
			Totaloperateincomeratio float32 `json:"TOI_RATIO"`
			Totaloperatecost        float32 `json:"TOTAL_OPERATE_COST"`
			Operatecost             float32 `json:"OPERATE_COST"`
			Saleexpense             float32 `json:"SALE_EXPENSE"`
			Manageexpense           float32 `json:"MANAGE_EXPENSE"`
			Financeexpense          float32 `json:"FINANCE_EXPENSE"`
			Operateprofit           float32 `json:"OPERATE_PROFIT"`
			Totalprofit             float32 `json:"TOTAL_PROFIT"`
		} `json:"data"`
		Pages int `json:"pages"`
	} `json:"result"`
}

// FetchAndStore 函数返回插入的数据数量
func FetchAndStore(db *gorm.DB, baseURL string, page, pageSize int, rptdate string) (int, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("Failed to parse URL: %v", err)
		return 0, err
	}

	query := u.Query()
	//query.Set("callback", "mydata")
	query.Set("sortColumns", "NOTICE_DATE,SECURITY_CODE")
	query.Set("sortTypes", "-1,-1")
	query.Set("pageNumber", strconv.Itoa(page))
	query.Set("pageSize", strconv.Itoa(pageSize))
	query.Set("reportName", "RPT_DMSK_FN_INCOME")
	query.Set("columns", "ALL")

	var filter = "(SECURITY_TYPE_CODE+in+(%22058001001%22%2C%22058001008%22))(TRADE_MARKET_CODE!%3D%22069001017%22)(REPORT_DATE%3D%27" + rptdate + "%27)"
	//query.Set("filter", filter)

	u.RawQuery = query.Encode()
	//fmt.Println(u.String() + "&filter=" + filter)
	resp, err := http.Get(u.String() + "&filter=" + filter)
	//url := "https://datacenter-web.eastmoney.com/api/data/v1/get?callback=mydata&sortColumns=UPDATE_DATE%2CSECURITY_CODE&sortTypes=-1%2C-1&pageSize=50&pageNumber=1&reportName=RPT_LICO_FN_CPD&columns=ALL&filter=(SECURITY_TYPE_CODE+in+(%22058001001%22%2C%22058001008%22))(TRADE_MARKET_CODE!%3D%22069001017%22)(REPORTDATE%3D%272024-09-30%27)"

	if err != nil {
		log.Printf("Failed to fetch data from URL %s: %v", u.String(), err)
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d while fetching data", resp.StatusCode)
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应内容并去除多余字符
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return 0, err
	}

	//fmt.Println(string(body))
	var data ResponseData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, err
	}
	if page == 1 {
		fmt.Println(data.Result)
	}
	for _, i := range data.Result.Data {
		result := db.Exec(` INSERT INTO ft_income (stockid,stockname,reportdate,noticedate,
            netprofit,totaloperateincome, totaloperatecost, operatecost,saleexpense,manageexpense,
            financeexpense, operateprofit, totalprofit, totaloperateincomeratio, netprofitratio, 
            upttime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW() )
	        ON DUPLICATE KEY UPDATE  stockname=?,noticedate=?,
            netprofit=?,totaloperateincome=?, totaloperatecost=?, operatecost=?,saleexpense=?,manageexpense=?,
            financeexpense=?, operateprofit=?, totalprofit=?, totaloperateincomeratio=?, netprofitratio=?,upttime = NOW()`,
			i.StockID, i.StockName, i.Reportdate, i.Noticedate,
			i.Netprofit, i.Totaloperateincome, i.Totaloperatecost, i.Operatecost, i.Saleexpense, i.Manageexpense,
			i.Financeexpense, i.Operateprofit, i.Totalprofit, i.Totaloperateincomeratio, i.Netprofitratio,
			i.StockName, i.Noticedate,
			i.Netprofit, i.Totaloperateincome, i.Totaloperatecost, i.Operatecost, i.Saleexpense, i.Manageexpense,
			i.Financeexpense, i.Operateprofit, i.Totalprofit, i.Totaloperateincomeratio, i.Netprofitratio)
		//AI: 修改mysql语句，如果stockid已经存在，不添加数据，只修改字段upttime为现在时间
		if result.Error != nil {
			log.Printf("Failed to store newqtrpt: %v", result.Error)
			return 0, err
		}
	}
	log.Printf("Inserted %d records into database", len(data.Result.Data))
	return len(data.Result.Data), nil
	//return 1, nil
}
