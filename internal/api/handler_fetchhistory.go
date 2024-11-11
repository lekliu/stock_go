package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/history"
)

// 获取股票历史数据
// https://q.stock.sohu.com/cn/300750/lshq.shtml
// https://q.stock.sohu.com/hisHq?code=cn_300750&start=20240704&end=20241101&stat=1&order=D&period=d&callback=history&rt=jsonp&r=0.6059332004266027&0.7019642347784691
func (h *Hdl_fetch) FetchAndStore_history(w http.ResponseWriter, r *http.Request) {
	var FetchURL = "https://q.stock.sohu.com/hisHq"
	//page, pageSize := 1, 50

	id := r.URL.Query().Get("id")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if id == "" || start == "" || end == "" {
		json.NewEncoder(w).Encode(Response{Message: "参数缺失,示例：id=300750&start=20240704&end=20241101 或 id=300750,000064&start=20240704&end=20241101", Success: false})
		return
	}

	go func() {
		recordsFetched, err := stock.FetchAndStore(h.DB, FetchURL, id, start, end)
		if err != nil {
			log.Printf("Failed to fetch and store stocks, url：%s, Err: %v \n", FetchURL, err)
			return
		}
		log.Printf("Fetched page %d successfully", recordsFetched)
	}()

	json.NewEncoder(w).Encode(Response{Message: "Start Data fetched and stored ......", Success: true})
}
