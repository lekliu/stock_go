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
	"strings"
)

type ResponseData struct {
	Data struct {
		Diff []struct {
			StockID   string `json:"f12"`
			StockName string `json:"f14"`
		} `json:"diff"`
		Total int `json:"total"`
	} `json:"data"`
}

// FetchAndStore 函数返回插入的数据数量
func FetchAndStore(db *gorm.DB, baseURL string, page, pageSize int) (int, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("Failed to parse URL: %v", err)
		return 0, err
	}

	query := u.Query()
	query.Set("cb", "mydata")
	query.Set("pn", strconv.Itoa(page))
	query.Set("pz", strconv.Itoa(pageSize))
	query.Set("po", "1")
	query.Set("np", "1")
	query.Set("fltt", "2")
	query.Set("invt", "2")
	query.Set("fid", "f3")
	query.Set("fs", "m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23")
	query.Set("fields", "f12,f14")

	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
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

	// 去除前面的“mydata(”和末尾的“);”
	trimmedBody := strings.TrimPrefix(string(body), "mydata(")
	trimmedBody = strings.TrimSuffix(trimmedBody, ");")

	var data ResponseData
	if err := json.Unmarshal([]byte(trimmedBody), &data); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, err
	}

	for _, stock := range data.Data.Diff {
		if strings.HasPrefix(stock.StockName, "XD") {
			fmt.Println(stock.StockID, stock.StockName)
			continue
		}
		result := db.Exec(` INSERT INTO ft_stock (stockid, stockname, upttime) 
                                VALUES (?, ?, NOW()) 
                                ON DUPLICATE KEY UPDATE  upttime = NOW()`,
			stock.StockID, stock.StockName)
		//AI: 修改mysql语句，如果stockid已经存在，不添加数据，只修改字段upttime为现在时间
		if result.Error != nil {
			log.Printf("Failed to store stock: %v", result.Error)
			return 0, err
		}
	}
	log.Printf("Inserted %d records into database", len(data.Data.Diff))
	return len(data.Data.Diff), nil
}
