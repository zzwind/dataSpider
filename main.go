package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//this program for data spider
//以沉浸的心态做好这件事，忘掉努力奋斗坚持。
func main() {
	t := time.Now()

	//要执行的任务1：元数据  2：历史数据
	mission := 2

	if mission == 1 {

		//fmt.Println("STW")
		//Download("http://quotes.money.163.com/fund/jzzs_163113.html?start=2015-12-01&end=2017-01-26")
		GetMetaData("开放型", "http://quotes.money.163.com/fn/service/netvalue.php?host=/fn/service/netvalue.php&query=STYPE:FDO&fields=no,PUBLISHDATE,SYMBOL,SNAME,NAV,PCHG,M12RETRUN,SLNAVG,LJFH,ZJZC,SGZT,OFPROFILE15&sort=PCHG&order=desc&count=6000&type=query")
		GetMetaData("货币型", "http://quotes.money.163.com/fn/service/netvalue.php?host=/fn/service/netvalue.php&page=0&query=TYPE1:MONEY&fields=no,PUBLISHDATE,SYMBOL,SNAME,CUR4,CURNAV_001,CURNAV_010,CURNAV_011,OFPROFILE8,OFPROFILE15&sort=CUR4&order=desc&count=600&type=query")
		GetMetaData("保本型", "http://quotes.money.163.com/fn/service/netvalue.php?host=/fn/service/netvalue.php&page=0&query=STYPE:FDO;TYPE3:BBX&fields=no,PUBLISHDATE,SYMBOL,SNAME,NAV,PCHG,M12RETRUN,SLNAVG,LJFH,ZJZC,SGZT,OFPROFILE15&sort=PCHG&order=desc&count=200&type=query")
		GetMetaData("股票型", "http://quotes.money.163.com/fn/service/netvalue.php?host=/fn/service/netvalue.php&page=0&query=STYPE:FDO;TYPE3:GPX&fields=no,PUBLISHDATE,SYMBOL,SNAME,NAV,PCHG,M12RETRUN,SLNAVG,LJFH,ZJZC,SGZT,OFPROFILE15&sort=PCHG&order=desc&count=300&type=query")
		GetMetaData("指数型", "http://quotes.money.163.com/fn/service/netvalue.php?host=/fn/service/netvalue.php&page=0&query=STYPE:FDO;TYPE3:ZSX&fields=no,PUBLISHDATE,SYMBOL,SNAME,NAV,PCHG,M12RETRUN,SLNAVG,LJFH,ZJZC,SGZT,OFPROFILE15&sort=PCHG&order=desc&count=900&type=query")
		GetMetaData("混合型", "http://quotes.money.163.com/fn/service/netvalue.php?host=/fn/service/netvalue.php&page=0&query=STYPE:FDO;TYPE3:HHX&fields=no,PUBLISHDATE,SYMBOL,SNAME,NAV,PCHG,M12RETRUN,SLNAVG,LJFH,ZJZC,SGZT,OFPROFILE15&sort=PCHG&order=desc&count=3000&type=query")
		GetMetaData("债券型", "http://quotes.money.163.com/fn/service/netvalue.php?host=/fn/service/netvalue.php&page=0&query=STYPE:FDO;TYPE3:ZQX&fields=no,PUBLISHDATE,SYMBOL,SNAME,NAV,PCHG,M12RETRUN,SLNAVG,LJFH,ZJZC,SGZT,OFPROFILE15&sort=PCHG&order=desc&count=2000&type=query")

	} else if mission == 2 {
		StartTdData()
	}

	fmt.Println(time.Since(t).Seconds())
}

//<tbody>\n<tr>\n(<td>(.*)</td>\n){3}</tr>\n</tbody>
//<tbody>[\s\S\n]*</tbody>

func Download(url string) []byte {
	r, err := http.Get(url)
	defer r.Body.Close()
	if err != nil {
		panic("Download error " + url)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic("Download error " + url)
	}

	return b

	//ioutil.WriteFile("tbody.txt", b, os.ModePerm)

	//b, err := ioutil.ReadFile("tbody.txt")
	//if err != nil {
	//	panic(err)
	//}
	//
	//t := time.Now()

	//hmtx:=(string)b

	//reg0 := regexp.MustCompile(`<tbody>[\s\S\n]*</tbody>`)

	//matchTbody := reg0.Find(b)

	//regTd := regexp.MustCompile(`<td>(?:<span class=".*?">)?(.*?)(?:</span>)?</td>`) //<td>.*?([\d\-\.]*).*?</td>        <td.*?>([\d\-\.]*).*?</td>
	//matchTd := regTd.FindAllSubmatch(b, -1)
	//
	//for _, j := range matchTd {
	//	fmt.Printf("%s\n", j[1])
	//}

	//fmt.Println(time.Since(t).Seconds())

}
