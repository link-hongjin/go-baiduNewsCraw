/*
go craw package
*/
package gocraw
import (
	"os"
	"net/http"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/garyburd/redigo/redis"
	"github.com/tealeg/xlsx"
)


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

func fileOut(data []map[string]string, path string) error {

	xfile := xlsx.NewFile()

	sheet,err := xfile.AddSheet("sheet1")
	if err != nil {
		fmt.Println("创建xlsx文件失败！")
		return err
	}

	for i,m := range data {
		if 0 == i {
			row := sheet.AddRow()
			row.SetHeigthCM(1)
			for k,_ := range m {
				cell := row.AddCell()
				cell.Value = k
			}
		}
		row := sheet.AddRow()
		row.SetHeigthCM(1)
		for _,v := range m {
			cell := row.AddCell()
			cell.Value = v
		}
	}

	err = xfile.Save(path)
	if err != nil {
		fmt.Println("保存xlsx文件失败！")
		return err
	}

	fmt.Printf("保存成功：%s\n", path)
}


func redisConnect(conn string) (,error){
	c,err := redis.Dial("tcp", conn)

	if err != nil {
		return err
	}

	return c, nil
}





