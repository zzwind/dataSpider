package main

import (
	"encoding/json"
	"fmt"

	"database/sql"

	"log"

	_ "github.com/lib/pq"
)

//获取所有基金基础信息

//定义获取的json结构
type Fund struct {
	Page int `json:"page"`
	List []FundList
}

//定义基金json的结构
type FundList struct {
	//OFPROFILE8  string
	SNAME       string //基金名称
	PUBLISHDATE string //基金成立时间
	SYMBOL      string //基金代码
	Id          int    //ID
	LASTUPDATE  string //最后更新时间
}

func GetMetaData(fundtype string, url string) {

	var fund Fund
	//下载url
	b := Download(url)
	//解析下载的json数据
	err := json.Unmarshal(b, &fund)
	if err != nil {
		panic(err)
	}
	//链接数据库
	db, err := sql.Open("postgres", "user=postgres dbname=fund password=123456 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	//初始化insert sql语句
	stmtIns, err := db.Prepare("INSERT INTO fundmetadata(symbol,sname,publishdate,fundtype) VALUES( $1, $2, $3, $4)") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	for _, j := range fund.List {

		if j.PUBLISHDATE == "" {
			j.PUBLISHDATE = "2000-01-01" //如果没有成立时间，就默认为此时间
		}

		_, err = stmtIns.Exec(j.SYMBOL, j.SNAME, j.PUBLISHDATE, fundtype) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	fmt.Println("done", fundtype)

}
