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
CREATE TABLE `ft_cashflow` (
  `id` int NOT NULL AUTO_INCREMENT,
  `stockid` varchar(8) NOT NULL COMMENT '代码',
  `stockname` varchar(15) NOT NULL COMMENT '名称',
  `reportdate` date DEFAULT NULL COMMENT '报告日期',
  `noticedate` date DEFAULT NULL COMMENT '公告日期',
  `netcashoperate` decimal(18,3) DEFAULT '0.000' COMMENT '经营性现金流',
  `netcashinvest` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '投资性现金流',
  `netcashfinance` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '融资性现金流',
  `netcashflow` decimal(18,3) NOT NULL DEFAULT '0.000' COMMENT '净现金流',
  `netcashflowratio` decimal(20,9) NOT NULL DEFAULT '0.000000000' COMMENT '净现金流同比',
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
			Netcashoperate   float64 `json:"NETCASH_OPERATE"`
			Netcashinvest    float64 `json:"NETCASH_INVEST"`
			Netcashfinance   float64 `json:"NETCASH_FINANCE"`
			Netcashflow      float64 `json:"CCE_ADD"`
			Netcashflowratio float64 `json:"CCE_ADD_RATIO"`
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
	query.Set("reportName", "RPT_DMSK_FN_CASHFLOW")
	query.Set("columns", "ALL")
	//query.Set("columns", "SECURITY_CODE,SECURITY_NAME_ABBR,BASIC_EPS,PARENT_NETPROFIT,TOTAL_OPERATE_INCOME,BPS,PUBLISHNAME,QDATE")

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
		result := db.Exec(` INSERT INTO ft_cashflow (stockid, stockname,reportdate,noticedate,
            netcashoperate,netcashinvest, netcashfinance,netcashflow,netcashflowratio,upttime)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	        ON DUPLICATE KEY UPDATE  stockname=?,noticedate=?, netcashoperate=?,
	        netcashinvest=?, netcashfinance=?,netcashflow=?,netcashflowratio=?,upttime = NOW()`,
			i.StockID, i.StockName, i.Reportdate, i.Noticedate,
			i.Netcashoperate, i.Netcashinvest, i.Netcashfinance, i.Netcashflow, i.Netcashflowratio,
			i.StockName, i.Noticedate,
			i.Netcashoperate, i.Netcashinvest, i.Netcashfinance, i.Netcashflow, i.Netcashflowratio)
		//AI: 修改mysql语句，如果stockid已经存在，不添加数据，只修改字段upttime为现在时间
		if result.Error != nil {
			log.Printf("Failed to store newqtrpt: %v", result.Error)
			return 0, err
		}
	}
	log.Printf("Inserted %d records into database", len(data.Result.Data))
	return len(data.Result.Data), nil
	return 1, nil
}
