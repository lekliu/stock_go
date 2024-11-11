package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/cashflow"
)

//获取最新的季报数据 new Quarterly Report
//https://data.eastmoney.com/bbsj/202403/xjll.html
//https://datacenter-web.eastmoney.com/api/data/v1/get?callback=jQueryi&sortColumns=NOTICE_DATE%2CSECURITY_CODE&sortTypes=-1%2C-1&pageSize=50&pageNumber=1&reportName=RPT_DMSK_FN_CASHFLOW&columns=ALL&filter=(SECURITY_TYPE_CODE+in+(%22058001001%22%2C%22058001008%22))(TRADE_MARKET_CODE!%3D%22069001017%22)(REPORT_DATE%3D%272024-03-31%27)
/*
	"SECURITY_CODE": "000851",代码
	"SECURITY_NAME_ABBR": "ST高鸿",名称
	"NOTICE_DATE": "2024-11-05 00:00:00",公告日期
	"REPORT_DATE": "2024-09-30 00:00:00",报告日期
	"NETCASH_OPERATE": -363166542.26,经营性现金流
	"NETCASH_OPERATE_RATIO": -457.6597430416,
	"SALES_SERVICES": 1401562694.02,
	"SALES_SERVICES_RATIO": 1766.2387575966,
	"PAY_STAFF_CASH": 97746385.67,
	"PSC_RATIO": 123.1792594951,
	"NETCASH_INVEST": 273161965.27,投资性现金流
	"NETCASH_INVEST_RATIO": 344.2366525733,
	"RECEIVE_INVEST_INCOME": 5670268,
	"RII_RATIO": 7.1456290541,
	"CONSTRUCT_LONG_ASSET": 61358720.16,
	"CLA_RATIO": 77.3237973053,
	"NETCASH_FINANCE": 10635900.79,融资性现金流
	"NETCASH_FINANCE_RATIO": 13.4032821203,
	"CCE_ADD": -79352957.69,净现金流
	"CCE_ADD_RATIO": 91.8969056396,净现金流同比
	"CUSTOMER_DEPOSIT_ADD": null,
	"CDA_RATIO": null,
	"DEPOSIT_IOFI_OTHER": null,
	"DIO_RATIO": null,
	"LOAN_ADVANCE_ADD": null,
	"LAA_RATIO": null,
	"RECEIVE_INTEREST_COMMISSION": null,
	"RIC_RATIO": null,
	"INVEST_PAY_CASH": null,
	"IPC_RATIO": null,
	"BEGIN_CCE": null,
	"BEGIN_CCE_RATIO": null,
	"END_CCE": null,
	"END_CCE_RATIO": null,
	"RECEIVE_ORIGIC_PREMIUM": null,
	"ROP_RATIO": null,
	"PAY_ORIGIC_COMPENSATE": null,
	"POC_RATIO": null
*/

func (h *Hdl_fetch) FetchAndStore_cashflow(w http.ResponseWriter, r *http.Request) {
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

			log.Printf("Fetched page %d successfully, rpttype: cashflow, rptdate: %s", page, rptdate)
			if recordsFetched < pageSize {
				break
			}
			page++
		}
	}()

	json.NewEncoder(w).Encode(Response{Message: "Start Data fetched and stored ......", Success: true})
}
