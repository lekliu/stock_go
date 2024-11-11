package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/newqtrpt"
)

//获取最新的季报数据 new Quarterly Report
//https://data.eastmoney.com/bbsj/202409.html
//https://datacenter-web.eastmoney.com/api/data/get?callback=jQuery&ps=50&p=3&type=RPT_LICO_FN_CPD&sty=ALL&filter=(REPORTDATE%3D%272024-09-30%27)
/*
stdClass Object ( [SECURITY_CODE] => 000504 [SECURITY_NAME_ABBR] => 南华生物 [TRADE_MARKET_CODE] => 069001002001
  [TRADE_MARKET] => 深交所主板 [SECURITY_TYPE_CODE] => 058001001
  [SECURITY_TYPE] => A股 [UPDATE_DATE] => 2021-04-29 00:00:00 [REPORTDATE] => 2021-03-31 00:00:00 [BASIC_EPS] => 0.0031
  [DEDUCT_BASIC_EPS] => [TOTAL_OPERATE_INCOME] => 33131542.07 [PARENT_NETPROFIT] => 957093.09 [WEIGHTAVG_ROE] => 3.54
  [YSTZ] => 79.0189728223 [SJLTZ] => 159.82 [BPS] => 0.088343803482 [MGJYXJJE] => 0.015996334847 [XSMLL] => 83.9038249148
  [YSHZ] => -53.2753 [SJLHZ] => -91.7849 [ASSIGNDSCRPT] => [PAYYEAR] => [PUBLISHNAME] => 医药制造 [ZXGXL] =>
  [NOTICE_DATE] => 2021-04-29 00:00:00 [ORG_CODE] => 10004341 [TRADE_MARKET_ZJG] => 0201 [ISNEW] => 1
  [QDATE] => 2021Q1 [DATATYPE] => 2021年 一季报 [DATAYEAR] => 2021 [DATEMMDD] => 一季报 [EITIME] => 2021-04-28 16:46:49
  [SECUCODE] => 000504.SZ )
*/

func (h *Hdl_fetch) FetchAndStore_newQtRpt(w http.ResponseWriter, r *http.Request) {
	var FetchURL = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	page, pageSize := 1, 100

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
