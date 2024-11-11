# 生成代码

提示词：
1、生成一个go语言开发项目，代码目录架构符合主流的设计方法
2、实现从web网站上下载数据存放到SQL中，将此方法以openapi的方式对外暴露
3、url是https://54.push2delay.eastmoney.com/api/qt/clist/get?cb=mydata&pn=1&pz=400&po=1&np=1&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23&fields=f12,f14
4、url中，pn代表页码，默认从1开始；pz是一次做返回的记录条数，默认为400；
5、f12是返回的股票代码stockid， f14是返回的股票名称stockname
6、下载的数据是json格式，返回字段中total为总记录数
7、数据库类型为mysql，存放在本地，账号是root,密码是root；库名为dbstock, 
8、数据存放到名为tbstock的表中,字段中id为自增字段，不需要处理；股票代码stockid，股票名称stockname


## 测试
go run stock/cmd/server/main.go
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/fetch_stocks" -Method GET
$response | Format-List