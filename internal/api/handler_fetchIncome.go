package api

import (
	"encoding/json"
	"log"
	"net/http"
	"stock/internal/fetcher/income"
)

//获取最新的季报数据 new Quarterly Report
//https://data.eastmoney.com/bbsj/202403/xjll.html
//https://datacenter-web.eastmoney.com/api/data/v1/get?callback=jQuery&sortColumns=NOTICE_DATE%2CSECURITY_CODE&sortTypes=-1%2C-1&pageSize=50&pageNumber=1&reportName=RPT_DMSK_FN_INCOME&columns=ALL&filter=(SECURITY_TYPE_CODE+in+(%22058001001%22%2C%22058001008%22))(TRADE_MARKET_CODE!%3D%22069001017%22)(REPORT_DATE%3D%272024-03-31%27)
/*
	"SECURITY_CODE": "688356",代码
	"SECURITY_NAME_ABBR": "键凯科技",名称
	"NOTICE_DATE": "2024-11-04 00:00:00",公告日期
	"REPORT_DATE": "2024-09-30 00:00:00",报告日期
	"PARENT_NETPROFIT": 32653746.12,净利润
	"TOTAL_OPERATE_INCOME": 185726750.89,营业总收入
	"TOTAL_OPERATE_COST": 152853628.1,营业总支出
	"TOE_RATIO": 17.6446707312,
	"OPERATE_COST": 59607318.87,营业支出
	"OPERATE_EXPENSE": 59607318.87,
	"OPERATE_EXPENSE_RATIO": 22.577066233,
	"SALE_EXPENSE": 4904925.51,销售费用
	"MANAGE_EXPENSE": 39855867.15,管理费用
	"FINANCE_EXPENSE": -1686963.81,财务费用
	"OPERATE_PROFIT": 35545460.69,营业利润
	"TOTAL_PROFIT": 34933524.4,利润总额
	"INCOME_TAX": 2279778.28,
	"OPERATE_INCOME": null,
	"INTEREST_NI": null,
	"INTEREST_NI_RATIO": null,
	"FEE_COMMISSION_NI": null,
	"FCN_RATIO": null,
	"OPERATE_TAX_ADD": 4140962.27,
	"MANAGE_EXPENSE_BANK": null,
	"FCN_CALCULATE": null,
	"INTEREST_NI_CALCULATE": null,
	"EARNED_PREMIUM": null,
	"EARNED_PREMIUM_RATIO": null,
	"INVEST_INCOME": null,
	"SURRENDER_VALUE": null,
	"COMPENSATE_EXPENSE": null,
	"TOI_RATIO": -23.9675653023,营业总收入同比
	"OPERATE_PROFIT_RATIO": -68.156132813854,
	"PARENT_NETPROFIT_RATIO": -67.5,净利润同比
	"DEDUCT_PARENT_NETPROFIT": 25719008.98,
	"DPN_RATIO": -73.295826260978
*/

func (h *Hdl_fetch) FetchAndStore_income(w http.ResponseWriter, r *http.Request) {
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

			log.Printf("Fetched page %d successfully, rpttype: income, rptdate: %s", page, rptdate)
			if recordsFetched < pageSize {
				break
			}
			page++
		}
	}()

	json.NewEncoder(w).Encode(Response{Message: "Start Data fetched and stored ......", Success: true})
}
