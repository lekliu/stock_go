package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/stock"
)

//获取最新的股票列表
//
//http://quote.eastmoney.com/center/gridlist.html#hs_a_board
//http://54.push2.eastmoney.com/api/qt/clist/get?
// cb=mydata&pn=2&pz=50&po=1&np=1&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23
//&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152
/*
stdClass Object ( [f1] => 2 [f2] 最新价 => 2.66 [f3] => 1.14 [f4] => 0.03 [f5] => 726838 [f6] => 193026761 [f7] => 3.04
     [f8] => 3.22 [f9] 市盈率 => 1891.58 [f10] => 0.6 [f11] => 0.38 [f12] 代码=> 600759 [f13] => 1 [f14] 名称 => 洲际油气
	 [f15] => 2.7 [f16] => 2.62 [f17] => 2.64 [f18] => 2.63 [f20] 总市值 => 6020929998 [f21] => 6006769488
	 [f22] => 0 [f23] => 1.14 [f24] => 37.82 [f25] => 48.6 [f62] => -3714610 [f115] => 28.24 [f128] => -
	 [f140] => - [f141] => - [f136] => - [f152] => 2 )
*/

func (h *Hdl_fetch) FetchAndStore_stocks(w http.ResponseWriter, r *http.Request) {
	var FetchURL = "https://54.push2delay.eastmoney.com/api/qt/clist/get"
	page, pageSize := 1, 200

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
	return
}
