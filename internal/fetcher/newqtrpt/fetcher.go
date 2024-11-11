package stock

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"stock/utils/dateutil"
	"strconv"
	"strings"
)

/*
每股收益 (Earnings Per Share)简写：EPS
净利润 (Net Profit)简写：NP
营业总收入 (Total Operating Revenue)简写：TOR
每股净资产 (Net Asset Value Per Share) 简写：NAVPS
*/

type ResponseData struct {
	Result struct {
		Data []struct {
			StockID   string  `json:"SECURITY_CODE"`
			StockName string  `json:"SECURITY_NAME_ABBR"`
			EPS       float32 `json:"BASIC_EPS"`
			NP        float32 `json:"PARENT_NETPROFIT"`
			TOR       float32 `json:"TOTAL_OPERATE_INCOME"`
			NAVPS     float32 `json:"BPS"`
			Industry  string  `json:"PUBLISHNAME"`
			Qdate     string  `json:"QDATE"`
		} `json:"data"`
		Pages int `json:"pages"`
	} `json:"result"`
}

// FetchAndStore 函数返回插入的数据数量
func FetchAndStore(db *gorm.DB, baseURL string, page, pageSize int) (int, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("Failed to parse URL: %v", err)
		return 0, err
	}

	var lastQt1 = dateutil.GetLastQt1()
	query := u.Query()
	query.Set("callback", "mydata")
	query.Set("sortColumns", "UPDATE_DATE,SECURITY_CODE")
	query.Set("sortTypes", "-1,-1")
	query.Set("pageNumber", strconv.Itoa(page))
	query.Set("pageSize", strconv.Itoa(pageSize))
	query.Set("reportName", "RPT_LICO_FN_CPD")
	//query.Set("columns", "ALL")
	query.Set("columns", "SECURITY_CODE,SECURITY_NAME_ABBR,BASIC_EPS,PARENT_NETPROFIT,TOTAL_OPERATE_INCOME,BPS,PUBLISHNAME,QDATE")

	var filter = "(SECURITY_TYPE_CODE+in+(%22058001001%22%2C%22058001008%22))(TRADE_MARKET_CODE!%3D%22069001017%22)(REPORTDATE%3D%27" + lastQt1 + "%27)"
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

	// 去除前面的“mydata(”和末尾的“);”
	trimmedBody := strings.TrimPrefix(string(body), "mydata(")
	trimmedBody = strings.TrimSuffix(trimmedBody, ");")

	var data ResponseData
	if err := json.Unmarshal([]byte(trimmedBody), &data); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, err
	}
	//fmt.Println(data.Result)

	for _, i := range data.Result.Data {
		result := db.Exec(` INSERT INTO ft_newqtrpt 
 					(stockid,stockname,EPS,NP,TOR,NAVPS,industry,qdate,upttime)
	                VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	                 ON DUPLICATE KEY UPDATE  stockname=?,EPS=?,NP=?,TOR=?,NAVPS=?,industry=?,qdate=?,upttime = NOW()`,
			i.StockID, i.StockName, i.EPS, i.NP, i.TOR, i.NAVPS, i.Industry, i.Qdate,
			i.StockName, i.EPS, i.NP, i.TOR, i.NAVPS, i.Industry, i.Qdate)
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
