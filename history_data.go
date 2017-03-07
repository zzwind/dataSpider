package main

import (
	"database/sql"
	"fmt"
	"regexp"

	"time"

	"bytes"

	"encoding/binary"

	"math"

	_ "github.com/go-sql-driver/mysql"
)

const (
	LIST_URL = "http://quotes.money.163.com/fund/jzzs_%s.html?start=%s&end=%s&sort=TDATE&order=asc"
	BASE_URL = "http://quotes.money.163.com"
)

func StartTdData() {

	db, err := sql.Open("postgres", "user=postgres dbname=fund password=123456 sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT id,symbol,sname,publishdate,lastupdate FROM fundmetadata WHERE publishdate!='2000-01-01' AND fundtype='指数型' ORDER BY id ASC ")
	if err != nil {
		panic(err.Error())
	}
	// Fetch rows
	for rows.Next() {

		row := new(FundList)
		err = rows.Scan(&row.Id, &row.SYMBOL, &row.SNAME, &row.PUBLISHDATE, &row.LASTUPDATE)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(row.Id, row.SYMBOL, row.SNAME, row.PUBLISHDATE, row.LASTUPDATE)
		row.LASTUPDATE = row.LASTUPDATE[:10] //数据库查找出来的日期带时间，这里把时间截取掉
		fetchData(fmt.Sprintf(LIST_URL, row.SYMBOL, row.LASTUPDATE, time.Now().Format("2006-01-02")), db, row.Id)

		//os.Exit(0)

	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	//`http://quotes.money.163.com/fund/jzzs_163113.html?start=2016-09-01&end=2017-01-26`

}

func nextPage(b []byte, db *sql.DB, id int) {
	reg_next_page := regexp.MustCompile(`<a href="([^>]*?)" class="pages_flip">下一页</a>`) //正则表达式的惰性匹配是从头匹配到尾部，而不能从尾部匹配上来。
	matchNP := reg_next_page.FindAllSubmatch(b, -1)
	for _, j := range matchNP {
		fetchData(BASE_URL+string(j[1]), db, id)
	}
}

func fetchData(url string, db *sql.DB, id int) {
	fmt.Println(url)
	b := Download(url)
	getTdData(b, db, id)
	nextPage(b, db, id)
}

func getTdData(b []byte, db *sql.DB, id int) {
	reg1 := regexp.MustCompile(`<tr>\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>\s*</tr>`) //<td>.*?([\d\-\.]*).*?</td>        <td.*?>([\d\-\.]*).*?</td>
	stmtIns, err := db.Prepare("INSERT INTO history(fundid,value,valuetotal,date,rate) VALUES( $1, $2, $3, $4, $5 )")
	stmtUpd, err := db.Prepare("update fundmetadata set lastupdate=$1 WHERE id=$2")
	if err != nil {
		panic(err.Error())
	}

	matchTd := reg1.FindAllSubmatch(b, -1)
	for _, j := range matchTd {
		//2015-09-28 0.9980 0.9980 -0.30%
		rate := string(j[4])
		rate = rate[:len(rate)-1]
		//有些基金的累计值为--，so 在这里做了检测
		if bytes.Equal(j[3], []byte("--")) {
			fmt.Println(j[3])
			fmt.Printf("%s\n", j[3])
			//binary.LittleEndian.PutUint16(j[3], 0)
			j[3] = Float32ToByte(0.00)
		}

		//valuetotal, _ := strconv.ParseFloat(vtp, 32)

		_, err = stmtIns.Exec(id, j[2], j[3], j[1], rate) // Insert tuples (i, i^2)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		//fmt.Printf("update fundmetadata set lastupdate=%s WHERE id=%s\n", j[1], id)

		_, err := stmtUpd.Exec(j[1], id)

		if err != nil {

			panic(err.Error()) // proper error handling instead of panic in your app
		}

		//fmt.Printf("%s\n", string(j[1])+string(j[2])+string(j[3])+string(j[4]))
		//fmt.Println(j)
	}

}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}
