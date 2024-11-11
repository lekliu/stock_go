package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/theme"
)

//获取公司题材 theme
//https://data.eastmoney.com/gstc/
//https://data.eastmoney.com/dataapi/search/company?st=TOTAL_MARKET_VALUE&sr=-1&ps=20&p=2&mainPoint=BK

func (h *Hdl_fetch) FetchAndStore_theme(w http.ResponseWriter, r *http.Request) {
	var FetchURL = "https://data.eastmoney.com/dataapi/search/company"
	page, pageSize := 1, 50

	go func() {
		for {
			recordsFetched, err := stock.FetchAndStore(h.DB, FetchURL, page, pageSize)
			if err != nil {
				log.Printf("Failed to fetch and store stocks, url：%s, Err: %v \n", FetchURL, err)
				return
			}

			log.Printf("Fetched page %d successfully", page)
			if recordsFetched < pageSize {
				break
			}
			page++
		}
	}()

	json.NewEncoder(w).Encode(Response{Message: "Start Data fetched and stored ......", Success: true})
}
