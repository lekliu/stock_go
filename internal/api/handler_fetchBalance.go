package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/balance"
)

//获取最新的季报数据 new Quarterly Report
//https://data.eastmoney.com/bbsj/202403/xjll.html
//https://datacenter-web.eastmoney.com/api/data/v1/get?callback=jQuery&sortColumns=NOTICE_DATE%2CSECURITY_CODE&sortTypes=-1%2C-1&pageSize=50&pageNumber=1&reportName=RPT_DMSK_FN_BALANCE&columns=ALL&filter=(SECURITY_TYPE_CODE+in+(%22058001001%22%2C%22058001008%22))(TRADE_MARKET_CODE!%3D%22069001017%22)(REPORT_DATE%3D%272024-03-31%27)
/*
	"SECURITY_CODE": "688356",代码
	"SECURITY_NAME_ABBR": "键凯科技",名称
	"SECURITY_TYPE_CODE": "058001001",
	"TRADE_MARKET_CODE": "069001001006",
	"DATE_TYPE_CODE": "004",
	"REPORT_TYPE_CODE": "001",
	"DATA_STATE": "2",
	"NOTICE_DATE": "2024-11-04 00:00:00",公告日期
	"REPORT_DATE": "2024-09-30 00:00:00",报告日期
	"TOTAL_ASSETS": 1340971881.8,总资产
	"FIXED_ASSET": 395699546.68,
	"MONETARYFUNDS": 129161234.54,货币资金
	"MONETARYFUNDS_RATIO": -25.8144989273,
	"ACCOUNTS_RECE": 75543560.03,应收账款
	"ACCOUNTS_RECE_RATIO": -38.7609071561,
	"INVENTORY": 104429689.13,存货
	"INVENTORY_RATIO": 42.072069267,
	"TOTAL_LIABILITIES": 73605499.87,总负债
	"ACCOUNTS_PAYABLE": 4884678.74,应付账款
	"ACCOUNTS_PAYABLE_RATIO": 21.0203939366,
	"ADVANCE_RECEIVABLES": null,
	"ADVANCE_RECEIVABLES_RATIO": null,
	"TOTAL_EQUITY": 1267366381.93,股东权益合计
	"TOTAL_EQUITY_RATIO": -1.7669919663,
	"TOTAL_ASSETS_RATIO": -0.559416759,总资产同比
	"TOTAL_LIAB_RATIO": 26.1399776707,总负债同比
	"CURRENT_RATIO": 1237.2183852846,
	"DEBT_ASSET_RATIO": 5.4889666867,
*/

func (h *Hdl_fetch) FetchAndStore_balance(w http.ResponseWriter, r *http.Request) {
	var FetchURL = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	page, pageSize := 1, 100

	rptdate := r.URL.Query().Get("rptdate")
	if rptdate == "" {
		json.NewEncoder(w).Encode(Response{Message: "参数缺失,示例：rptdate=2024-03-31 或 2024-06-30 或 2024-09-30 或 2024-12-31", Success: false})
		return
	}

	go func() {
		for {
			recordsFetched, err := stock.FetchAndStore(h.DB, FetchURL, page, pageSize, rptdate)
			if err != nil {
				log.Printf("Failed to fetch and store stocks, url：%s, Err: %v \n", FetchURL, err)
				return
			}

			log.Printf("Fetched page %d successfully, rpttype: balance, rptdate: %s", page, rptdate)
			if recordsFetched < pageSize {
				break
			}
			page++
		}
	}()

	json.NewEncoder(w).Encode(Response{Message: "Start Data fetched and stored ......", Success: true})
}
