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
所属板块 - sector
主营业务 - mainbusiness
总股本 - totalshares
流通股本 - outstandingshares
*/

type ResponseData struct {
	Result struct {
		CompanyInfo []struct {
			StockID           string `json:"securityCode"`
			StockName         string `json:"securityShortName"`
			Sector            string `json:"bk"`
			Mainbusiness      string `json:"mainBusiness"`
			Totalshares       string `json:"totalCapital"`
			Outstandingshares string `json:"circulationCapital"`
		} `json:"companyInfo"`
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

	query := u.Query()
	query.Set("st", "TOTAL_MARKET_VALUE")
	query.Set("sr", strconv.Itoa(-1))
	query.Set("p", strconv.Itoa(page))
	query.Set("ps", strconv.Itoa(pageSize))
	query.Set("mainPoint", "BK")

	u.RawQuery = query.Encode()
	resp, err := http.Get(u.String())
	//url := "https://data.eastmoney.com/dataapi/search/company?st=TOTAL_MARKET_VALUE&sr=-1&ps=20&p=2&mainPoint=BK"

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
	if err := json.Unmarshal([]byte(string(body)), &data); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, err
	}
	//fmt.Println(data.Result)

	for _, i := range data.Result.CompanyInfo {
		result := db.Exec(` INSERT INTO ft_theme 
 					(stockid,stockname,sector,mainbusiness,totalshares,outstandingshares,upttime)
	                VALUES (?, ?, ?, ?, ?, ?, NOW())
	                 ON DUPLICATE KEY UPDATE  stockname=?,sector=?,mainbusiness=?,totalshares=?,outstandingshares=?,upttime = NOW()`,
			i.StockID, i.StockName, i.Sector, i.Mainbusiness, i.Totalshares, i.Outstandingshares,
			i.StockName, i.Sector, i.Mainbusiness, i.Totalshares, i.Outstandingshares)
		//AI: 修改mysql语句，如果stockid已经存在，不添加数据，只修改字段upttime为现在时间
		if result.Error != nil {
			log.Printf("Failed to store newqtrpt: %v", result.Error)
			return 0, err
		}
	}
	log.Printf("Inserted %d records into database", len(data.Result.CompanyInfo))
	return len(data.Result.CompanyInfo), nil
	//return 1, nil
}
