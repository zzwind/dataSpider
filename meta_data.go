package main

import (
	"encoding/json"
	"fmt"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//获取所有基金基础信息

type Fund struct {
	Page int `json:"page"`
	List []FundList
}

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

	b := Download(url)

	err := json.Unmarshal(b, &fund)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", "root:123456@/fund")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	stmtIns, err := db.Prepare("INSERT INTO fundmetadata(symbol,sname,publishdate,fundtype) VALUES( ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	for _, j := range fund.List {

		if j.PUBLISHDATE == "" {
			j.PUBLISHDATE = "2000-00-00" //如果没有成立时间，就默认为此时间
		}

		_, err = stmtIns.Exec(j.SYMBOL, j.SNAME, j.PUBLISHDATE, fundtype) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	fmt.Println("done", fundtype)

	//session, err := mgo.Dial("localhost")
	//if err != nil {
	//	panic(err)
	//}
	//defer session.Close()
	//
	//session.SetMode(mgo.Monotonic, true)
	//
	//c := session.DB("test").C("fundMetaData")
	//
	//for _, j := range fund.List {
	//
	//	err = c.Insert(&j)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

}

//func main() {
//	GetMetaData()
//}
