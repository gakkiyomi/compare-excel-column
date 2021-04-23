package handler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gakkiyomi/galang/array"
	str "github.com/gakkiyomi/galang/string"
	"github.com/gakkiyomi/galang/utils"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

// Compare 查询2个Excel文件中列的重复数据
// @Summary 查询2个Excel文件中列的重复数据
// @Description 查询2个Excel文件中列的重复数据
// @Accept  multipart/form-data
// @Produce  application/octet-stream
// @Param file1 formData file true "文件1"
// @Param file2 formData file true "文件2"
// @Param line1 query string  true "文件1数据是从第几行开始"
// @Param column1 query string  true "文件1要对比第几列"
// @Param line2 query string  true "文件2数据是从第几行开始"
// @Param column2 query string  true "文件2要对比第几列"
// @Success 200 {string} Resp "{"code":200,"data":{},"msg":"OK"}"
// @Failure 401 {string} Resp "{"code":401,"data":{},"msg":"UNAUTHORIZED"}"
// @Failure 500 {string} Resp "{"code":500,"data":{},"msg":"SERVERERROR"}"
// @Failure 80002 {string} Resp "{"code":80002,"data":{},"msg":"PARAMCHECKFAILED"}"
// @Failure 80003 {string} Resp "{"code":80003,"data":{},"msg":"DATABASEOPFAILED"}""
// @Router /api/compare [post]
func Compare(c *gin.Context) {

	var (
		line1   int
		line2   int
		file1   multipart.File
		file2   multipart.File
		column1 int
		column2 int
		err1    error
	)

	line1, _ = utils.Transform.StringToInt(c.Query("line1"))
	line2, _ = utils.Transform.StringToInt(c.Query("line2"))
	column1, _ = utils.Transform.StringToInt(c.Query("column1"))
	column2, _ = utils.Transform.StringToInt(c.Query("column2"))

	fmt.Println(line1, line2, column1, column2)

	file1, _, err := c.Request.FormFile("file1")
	if err != nil {
		c.JSON(500, "系统错误")
		return
	}
	file2, _, err = c.Request.FormFile("file2")
	if err != nil {
		c.JSON(500, "系统错误")
		return
	}
	defer file1.Close()
	defer file2.Close()

	temp := line1 - 1
	temp2 := column1 - 1
	temp3 := line2 - 1
	temp4 := column2 - 1

	ch1 := make(chan int)
	ch2 := make(chan int)
	list1 := make([]string, 0)
	list2 := make([]string, 0)
	build := str.String.NewStringBuilder("")

	go func() {
		bytes, _ := ioutil.ReadAll(file1)
		f, err := xlsx.OpenBinary(bytes)
		if err != nil {
			err1 = errors.New("解析文件错误，请检查文件格式")
			return
		}
		for _, sheet := range f.Sheets {
			for i, row := range sheet.Rows {

				if i < temp {
					continue
				}

				if len(row.Cells) < temp2 || row.Cells[temp2].Value == "" {
					continue
				}
				name := row.Cells[temp2].Value
				list1 = append(list1, name)
			}
		}
		ch1 <- 1
	}()

	go func() {
		bytes2, _ := ioutil.ReadAll(file2)
		f2, err := xlsx.OpenBinary(bytes2)
		if err != nil {
			err1 = errors.New("解析文件错误，请检查文件格式")
			return
		}

		for _, sheet := range f2.Sheets {
			for i, row := range sheet.Rows {
				if i < temp3 {
					continue
				}

				if len(row.Cells) < temp4 || row.Cells[temp4].Value == "" {
					continue
				}
				name := row.Cells[temp4].Value
				list2 = append(list2, name)
			}
		}
		ch2 <- 1
	}()

	<-ch1
	<-ch2
	if err1 != nil {
		c.JSON(500, "解析文件错误，请检查文件格式")
		return
	}
	intersect := array.Array.GetIntersectForString(list1, list2)
	for _, name := range intersect {
		build.Append(name).Append("\n")
	}

	content := []byte(build.ToString())
	err = ioutil.WriteFile("./result.txt", content, 0644)
	if err != nil {
		c.JSON(500, "服务繁忙")
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=result.txt")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Header("Content-Type", "application/octet-stream")
	//c.Writer.Write([]byte(build.ToString()))
	c.File("./result.txt")
}
