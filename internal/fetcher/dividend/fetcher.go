package stock

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"stock/utils/HttpClient"
	"strconv"
)

/*
所属板块 - sector
主营业务 - mainbusiness
总股本 - totalshares
流通股本 - outstandingshares
*/

type ResponseData struct {
	Result struct {
		Date []struct {
			StockID      string  `json:"SECURITY_CODE"`
			StockName    string  `json:"SECURITY_NAME_ABBR"`
			Dividenddate string  `json:"EX_DIVIDEND_DATE"`
			Bonus        float64 `json:"PRETAX_BONUS_RMB"`
			Ratio        float64 `json:"BONUS_IT_RATIO"`
		} `json:"data"`
		Pages int `json:"pages"`
	} `json:"result"`
}

// FetchAndStore 函数返回插入的数据数量
func FetchAndStore(db *gorm.DB, baseURL string, page, pageSize int, stockid string) (int, error) {
	client, err := HttpClient.GetHttpClientHandle()
	if err != nil {
		log.Printf("Get Http Client Handle : %v", err)
		return 0, err
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("Failed to parse URL: %v", err)
		return 0, err
	}

	query := u.Query()
	query.Set("sortColumns", "REPORT_DATE")
	query.Set("sortTypes", strconv.Itoa(-1))
	query.Set("pageNumber", strconv.Itoa(page))
	query.Set("pageSize", strconv.Itoa(pageSize))
	query.Set("reportName", "RPT_SHAREBONUS_DET")
	//query.Set("columns", "ALL")
	query.Set("columns", "SECURITY_CODE,SECURITY_NAME_ABBR,EX_DIVIDEND_DATE,PRETAX_BONUS_RMB,BONUS_IT_RATIO")
	query.Set("quoteColumns", "")
	//query.Set("source", "WEB")
	//query.Set("client", "WEB")
	filter := "&filter=(SECURITY_CODE%3D%22" + stockid + "%22)"
	u.RawQuery = query.Encode()

	//url := "https://data.eastmoney.com/dataapi/search/company?st=TOTAL_MARKET_VALUE&sr=-1&ps=20&p=2&mainPoint=BK"
	req, err := HttpClient.NewRequest("GET", u.String()+filter)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch data from URL %s: %v", u.String(), err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code %d while fetching data, url is %v", resp.StatusCode, u.String()+filter)
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
	if err := json.Unmarshal([]byte(string(body)), &data); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, err
	}
	fmt.Println(data.Result)

	for _, i := range data.Result.Date {
		if i.Dividenddate == "" {
			continue
		}
		result := db.Exec(` INSERT INTO ft_dividend
					(stockid,stockname,dividenddate,bonus,ratio,upttime)
	               VALUES (?, ?, ?, ?, ?, NOW())
	                ON DUPLICATE KEY UPDATE  stockname=?,dividenddate=?,bonus=?,ratio=?,upttime = NOW()`,
			i.StockID, i.StockName, i.Dividenddate, i.Bonus, i.Ratio,
			i.StockName, i.Dividenddate, i.Bonus, i.Ratio)
		//AI: 修改mysql语句，如果stockid已经存在，不添加数据，只修改字段upttime为现在时间
		if result.Error != nil {
			log.Printf("Failed to store newqtrpt: %v", result.Error)
			return 0, err
		}
	}
	log.Printf("Inserted %d records into database", len(data.Result.Date))
	return len(data.Result.Date), nil
}
