package stock

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"stock/utils/HttpClient"
	"strings"
)

/*
所属板块 - sector
主营业务 - mainbusiness
总股本 - totalshares
流通股本 - outstandingshares
*/

type StockData struct {
	Status int        `json:"status"`
	Hq     [][]string `json:"hq"`
	Code   string     `json:"code"`
	//Stat   []string   `json:"stat"`
}

// FetchAndStore 函数返回插入的数据数量
func FetchAndStore(db *gorm.DB, baseURL string, stockIds, start, end string) (int, error) {
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

	code := "cn_" + strings.ReplaceAll(stockIds, ",", ",cn_")

	query := u.Query()
	query.Set("code", code)
	query.Set("start", start)
	query.Set("end", end)
	query.Set("stat", "1")
	query.Set("order", "A")
	query.Set("callback", "history")
	query.Set("rt", "jsonp")
	query.Set("r", "0.6059332004266027")
	u.RawQuery = query.Encode()

	req, err := HttpClient.NewRequest("GET", u.String())
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
		log.Printf("Unexpected status code %d while fetching data, url is %v", resp.StatusCode, u.String())
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 检查内容编码
	var body []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		// 如果是 Gzip 压缩，解压缩
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Printf("error creating gzip reader: %v", err)
			return 0, err
		}
		defer reader.Close()

		body, err = ioutil.ReadAll(reader)
		if err != nil {
			log.Printf("error reading gzip response: %v", err)
			return 0, err
		}
	} else {
		// 否则直接读取响应体
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error reading response body: %v", err)
			return 0, err
		}
	}

	//fmt.Println(body)
	//fmt.Println(string(body))
	// 去除前面的“history(”和末尾的“)”
	trimmedBody := strings.TrimPrefix(string(body), "history(")
	//小括号后面有换行符，所以去除两个字符
	trimmedBody = trimmedBody[:len(trimmedBody)-2]

	//fmt.Println(trimmedBody)
	// 解析 JSON 数据
	var stockData []StockData
	if err = json.Unmarshal([]byte(trimmedBody), &stockData); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, err
	}

	// 处理并输出 Stat 数据
	for _, data := range stockData {
		stockid := strings.TrimPrefix(data.Code, "cn_")
		fmt.Printf("Code: %s\n", stockid)
		for _, item := range data.Hq {
			hqdate := item[0]
			open := item[1]
			closeprice := item[2]
			changerange := item[4]
			hign := item[5]
			low := item[6]
			turnoverrate := item[9]
			if strings.HasSuffix(changerange, "%") {
				changerange = strings.TrimSuffix(changerange, "%")
			}
			if strings.HasSuffix(turnoverrate, "%") {
				turnoverrate = strings.TrimSuffix(turnoverrate, "%")
			}

			result := db.Exec(` INSERT INTO ft_history
					(stockid,hqdate,open,close,changerange,hign,low,turnoverrate,upttime)
	              VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	               ON DUPLICATE KEY UPDATE  open=?,close=?,changerange=?,hign=?,low=?,turnoverrate=?,upttime = NOW()`,
				stockid, hqdate, open, closeprice, changerange, hign, low, turnoverrate,
				open, closeprice, changerange, hign, low, turnoverrate)
			//AI: 修改mysql语句，如果stockid已经存在，不添加数据，只修改字段upttime为现在时间
			if result.Error != nil {
				log.Printf("Failed to store newqtrpt: %v", result.Error)
				return 0, err
			}
		}
	}

	log.Printf("Inserted %d records into database", len(stockData))
	return len(stockData), nil
}
