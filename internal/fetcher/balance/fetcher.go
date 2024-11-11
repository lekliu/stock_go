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
CREATE TABLE `ft_balance` (
  `id` int NOT NULL AUTO_INCREMENT,
  `stockid` varchar(8) NOT NULL COMMENT '代码',
  `stockname` varchar(15) NOT NULL COMMENT '名称',
  `reportdate` date DEFAULT NULL COMMENT '报告日期',
  `noticedate` date DEFAULT NULL COMMENT '公告日期',
  `totalassets` decimal(18,3) DEFAULT '0.000' COMMENT '总资产',
  `totalassetsratio` decimal(15,9) NOT NULL DEFAULT '0.000000000' COMMENT '总资产同比',
  `monetaryfunds` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '货币资金',
  `accountsreceive` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '应收账款',
  `inventory` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '存货',
  `totalliabilities` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '总负债',
  `totalliabratio` decimal(15,9) NOT NULL DEFAULT '0.000000000' COMMENT '总负债同比',
  `accountspayable` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '应付账款',
  `totalequity` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '股东权益合计',
  `upttime` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `stockid_reportdate` (`stockid`,`reportdate`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb3;
*/

type ResponseData struct {
	Result struct {
		Data []struct {
			StockID          string  `json:"SECURITY_CODE"`
			StockName        string  `json:"SECURITY_NAME_ABBR"`
			Reportdate       string  `json:"REPORT_DATE"`
			Noticedate       string  `json:"NOTICE_DATE"`
			Totalassets      float32 `json:"TOTAL_ASSETS"`
			Totalassetsratio float32 `json:"TOTAL_ASSETS_RATIO"`
			Monetaryfunds    float32 `json:"MONETARYFUNDS"`
			Accountsreceive  float32 `json:"ACCOUNTS_RECE"`
			Inventory        float32 `json:"INVENTORY"`
			Totalliabilities float32 `json:"TOTAL_LIABILITIES"`
			Totalliabratio   float32 `json:"TOTAL_LIAB_RATIO"`
			Accountspayable  float32 `json:"ACCOUNTS_PAYABLE"`
			Totalequity      float32 `json:"TOTAL_EQUITY"`
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
	query.Set("reportName", "RPT_DMSK_FN_BALANCE")
	query.Set("columns", "ALL")

	var filter = "(SECURITY_TYPE_CODE+in+(%22058001001%22%2C%22058001008%22))(TRADE_MARKET_CODE!%3D%22069001017%22)(REPORT_DATE%3D%27" + rptdate + "%27)"
	//query.Set("filter", filter)

	u.RawQuery = query.Encode()
	//fmt.Println(u.String() + "&filter=" + filter)
	resp, err := http.Get(u.String() + "&filter=" + filter)

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
		result := db.Exec(` INSERT INTO ft_balance (stockid, stockname, reportdate, noticedate, 
            totalassets, monetaryfunds, accountsreceive, inventory, totalliabilities, 
            accountspayable, totalequity, totalassetsratio, totalliabratio, upttime) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	        ON DUPLICATE KEY UPDATE stockname=?, noticedate=?, 
            totalassets=?, monetaryfunds=?, accountsreceive=?, inventory=?, totalliabilities=?, 
            accountspayable=?, totalequity=?, totalassetsratio=?, totalliabratio=?,upttime = NOW()`,
			i.StockID, i.StockName, i.Reportdate, i.Noticedate,
			i.Totalassets, i.Monetaryfunds, i.Accountsreceive, i.Inventory, i.Totalliabilities,
			i.Accountspayable, i.Totalequity, i.Totalassetsratio, i.Totalliabratio,
			i.StockName, i.Noticedate,
			i.Totalassets, i.Monetaryfunds, i.Accountsreceive, i.Inventory, i.Totalliabilities,
			i.Accountspayable, i.Totalequity, i.Totalassetsratio, i.Totalliabratio)
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
