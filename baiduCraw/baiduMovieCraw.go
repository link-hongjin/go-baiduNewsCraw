/*
百度2020最新电影排行榜爬取案例

@author https://github.com/link-hongjin
*/

package main
import (
	"fmt"
	"net/http"
	"log"
	"github.com/PuerkitoBio/goquery"
	"flag"
	"time"
	"strings"
	"github.com/axgle/mahonia"
	"strconv"
	"os"
)

const (
	//目标地址
	url = "http://top.baidu.com/buzz?b=26&fr=topindex"
)

var (
	//定时爬取时间，单位：分钟
	t int
)

//电影数据结构
type movie struct {
	name string
	summary string
	summary_href string
}



func init() {
	flag.IntVar(&t, "t", 0, "定时时间（单位：分钟）")
	flag.Parse()
}

func main() {
	for {
		doc,err := request(url)
		if err != nil {
			log.Fatal(err.Error())
		}

		var data []movie

		doc.Find("tbody >tr").Each(func(i int, s *goquery.Selection){
			o := s.Find(".keyword")

			name := o.Find(".list-title").Text()
			if len(name) > 0 {
				name = convertToString(name, "gbk", "utf-8")
			
				o = s.Find(".tc >a").First()

				href,_ := o.Attr("href")

				//有些链接的 .html 后缀是缺漏的 导致无法请求 需要将其去掉
				href = strings.Replace(href, ".htm", "", -1)
				
				var m movie
				m.name = name
				m.summary_href = href

				//爬取简介
				doc,arr := request(href)
				if arr != nil {
					log.Fatal(arr.Error())
				}

				var summary string
				doc.Find(".lemma-summary >.para").Each(func(i int, s *goquery.Selection) {
					summary += s.Text()
				})
				summary = strings.Replace(summary, "\n", "", -1)
				if len(summary) <= 0 {
					summary = "暂无该电影简介"
				}
				m.summary = summary
				fmt.Println(name, summary)

				data = append(data, m)	
			}
		})
		path := "./baidu_movie_craw_result_" + time.Now().Format("20060102-15:04:05") + ".txt"

		_, err = os.Create(path)
	    if err != nil {
	        log.Fatal(err)
	        return
	    }

	    file, err := os.OpenFile(path, os.O_APPEND | os.O_RDWR, 0644)
	    defer file.Close()
	    if err != nil {
	        log.Fatal(err)
	        return
	    }

	    file.Write([]byte(fmt.Sprintf("%-5s|%-50s|%-150s|\n","排名","电影名","简介")))
	    for key,value := range data {
	    	file.Write([]byte(fmt.Sprintf("%-5s|%-50s|%-150s|\n", strconv.Itoa(key+1), value.name, value.summary)))
	    }


		fmt.Println(fmt.Sprintf("\n\n爬取成功！共爬取%d条数据，已保存到：%s", len(data), path))
		
		if t <= 0 {
	    	return
	    }

	    time.Sleep(time.Minute)
	}
}

func request(url string) (*goquery.Document, error) {
	res,err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return goquery.NewDocumentFromReader(res.Body)
}

func convertToString(src string, srcCode string, tagCode string) string {
    srcCoder := mahonia.NewDecoder(srcCode)
    srcResult := srcCoder.ConvertString(src)
    tagCoder := mahonia.NewDecoder(tagCode)
    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
    result := string(cdata)
    return result
}

