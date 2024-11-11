package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/dividend"
)

// 获取公司分红送股数据
// https://data.eastmoney.com/yjfp/detail/000063.html
// https://datacenter-web.eastmoney.com/api/data/v1/get?sortColumns=REPORT_DATE&sortTypes=-1&pageSize=50&pageNumber=1&reportName=RPT_SHAREBONUS_DET&columns=ALL&quoteColumns=&source=WEB&client=WEB&filter=(SECURITY_CODE%3D"000063")
func (h *Hdl_fetch) FetchAndStore_dividend(w http.ResponseWriter, r *http.Request) {
	var FetchURL = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	page, pageSize := 1, 50

	id := r.URL.Query().Get("id")
	if id == "" {
		json.NewEncoder(w).Encode(Response{Message: "id缺失,请在id中输入股票代码 ", Success: false})
		return
	}

	go func() {
		recordsFetched, err := stock.FetchAndStore(h.DB, FetchURL, page, pageSize, id)
		if err != nil {
			log.Printf("Failed to fetch and store stocks, url：%s, Err: %v \n", FetchURL, err)
			return
		}
		log.Printf("Fetched page %d successfully", recordsFetched)
	}()

	json.NewEncoder(w).Encode(Response{Message: "Start Data fetched and stored ......", Success: true})
}
