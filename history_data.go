package main

import (
	"database/sql"
	"fmt"
	"regexp"

	"time"

	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	LIST_URL = "http://quotes.money.163.com/fund/jzzs_%s.html?start=%s&end=%s"
	BASE_URL = "http://quotes.money.163.com"
)

func StartTdData() {

	db, err := sql.Open("mysql", "root:123456@/fund")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("SELECT id,symbol,sname,publishdate,lastupdate FROM fundmetadata WHERE publishdate!='2000-00-00' ORDER BY id ASC ")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Fetch rows
	for rows.Next() {

		row := new(FundList)
		// get RawBytes from data
		err = rows.Scan(&row.Id, &row.SYMBOL, &row.SNAME, &row.PUBLISHDATE, &row.LASTUPDATE)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Println(row.Id, row.SYMBOL, row.SNAME, row.PUBLISHDATE, row.LASTUPDATE)

		fetchData(fmt.Sprintf(LIST_URL, row.SYMBOL, row.LASTUPDATE, time.Now().Format("2006-01-02")))

		os.Exit(0)

	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	//`http://quotes.money.163.com/fund/jzzs_163113.html?start=2016-09-01&end=2017-01-26`

}

func nextPage(b []byte) {
	reg_next_page := regexp.MustCompile(`<a href="([^>]*?)" class="pages_flip">下一页</a>`) //正则表达式的惰性匹配是从头匹配到尾部，而不能从尾部匹配上来。

	matchNP := reg_next_page.FindAllSubmatch(b, -1)

	for _, j := range matchNP {
		//fmt.Printf("%s\n", BASE_URL+string(j[1]))
		fetchData(BASE_URL + string(j[1]))
	}

}

func fetchData(url string) {
	fmt.Println(url)
	b := Download(url)
	getTdData(b)
	nextPage(b)

}

func getTdData(b []byte) {
	reg1 := regexp.MustCompile(`<tr>\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*</tr>`) //<td>.*?([\d\-\.]*).*?</td>        <td.*?>([\d\-\.]*).*?</td>

	matchTd := reg1.FindAllSubmatch(b, -1)

	for _, j := range matchTd {
		fmt.Printf("%s\n", string(j[1])+string(j[2])+string(j[3])+string(j[4]))
		//fmt.Println(j)
	}

}
