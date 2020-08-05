/*
百度热搜榜单爬虫练习
涉及知识点：
HTTP、正则匹配、io操作

@author https://github.com/link-hongjin
*/
package main
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"regexp"
	"errors"
	"strconv"
	"os"
	"time"
	"flag"
)


const (
	//目标地址
	url = "https://tophub.today/n/Jb0vmloB1G"
)

var (
	//定时爬取时间，单位：分钟
	t int
)

func init() {
	flag.IntVar(&t, "t", 0, "定时时间（单位：分钟）")
	flag.Parse()
}

func main() {
	for {
		content,err := request2String(url)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		body,err := filter(content, `<div class="jc-c">(?s:(.*?))</div>`)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		body,err = filter(string(body[0][0]), `<tr>(?s:(.*?))</tr>`)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		var news [][]string

		for _,tr := range body {
			t := tr[0]

			td_all,e := filter(t, `<td class="al">(?s:(.*?))</td>`)
			if e != nil {
				log.Fatal(e.Error())
				continue
			}	

			//提取A标签
			a,e := filter(string(td_all[0][0]), `<a .*?href=['"](.*?)['"].*?>(.*?)</a>`)
			if e != nil {
				log.Fatal(e.Error())
				continue
			}	

			var d []string

			//名称
			d = append(d, string(a[0][2]))

			//链接
			d = append(d, string(a[0][1]))

			news = append(news, d)	
		}

		path := "./craw_result_" + time.Now().Format("20060102-15:04:05") + ".txt"

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

	    file.Write([]byte(fmt.Sprintf("%-5s|%-50s|%-100s|\n","排名","热搜头条","链接")))
	    for key,new := range news {
	    	file.Write([]byte(fmt.Sprintf("%-5s|%-50s|%-100s|\n", strconv.Itoa(key+1), new[0], new[1])))
	    }

	    fmt.Println(fmt.Sprintf("爬取成功！共计%d条数据，已保存到：%s", len(news), path))

	    if t <= 0 {
	    	return
	    }

	    time.Sleep(time.Minute)
	}
	

}


func request2String(url string) (string,error) {
	res,err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}

func filter(body string, preg string) ([][]string, error) {
	if len(body) <= 0 {
		return nil,errors.New("body不能为空")
	}

	reg := regexp.MustCompile(preg)
	if reg == nil {
		return nil,errors.New("正则表达式不正确")
	}

	result := reg.FindAllStringSubmatch(body, -1)
	if len(result) <= 0 {
		return nil,errors.New("匹配失败")
	}

	return result, nil

}

