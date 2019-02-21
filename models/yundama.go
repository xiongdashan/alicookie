package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/otwdev/galaxylib"
)

type Yundama struct {
	data     map[string]string
	apiuri   string
	Key      int
	codetype string
}

type RetData struct {
	Cid  int    `json:"cid"`
	Ret  int    `json:"ret"`
	Text string `json:"text"`
}

func NewYundama(codetype string) *Yundama {
	rev := &Yundama{}

	sect, _ := galaxylib.GalaxyCfgFile.GetSection("dama")

	rev.data = map[string]string{
		"username": sect["user"],
		"password": sect["pwd"],
		"appid":    "6899",
		"appkey":   "23a76fae3869526bbd5bdbae6d839d47",
		"timeout":  sect["timeout"],
	}

	rev.apiuri = sect["api"]
	rev.codetype = codetype
	return rev
}

func (n *Yundama) Get(imgBuf []byte) string {

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "code.png")

	if err != nil {
		fmt.Println(err)
		return ""
	}

	part.Write(imgBuf)

	for k, v := range n.data {
		writer.WriteField(k, v)
	}

	writer.WriteField("codetype", n.codetype)
	writer.WriteField("method", "upload")

	rq, _ := http.NewRequest("POST", n.apiuri, body)
	rq.Header.Add("Content-Type", writer.FormDataContentType())

	rs, err := http.DefaultClient.Do(rq)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer rs.Body.Close()

	buf, _ := ioutil.ReadAll(rs.Body)

	var revData *RetData

	err = json.Unmarshal(buf, &revData)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println(string(buf))

	if revData.Ret == 0 && revData.Text != "" {
		return revData.Text
	}

	code := n.result(revData.Cid)

	n.Key = revData.Cid

	return code
}

func (n *Yundama) result(cid int) string {

	uri := fmt.Sprintf("%s?cid=%d&method=result", n.apiuri, cid)

	rq, _ := http.NewRequest("POST", uri, nil)

	rs, err := http.DefaultClient.Do(rq)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer rs.Body.Close()

	buf, _ := ioutil.ReadAll(rs.Body)

	var revData *RetData

	json.Unmarshal(buf, &revData)

	if revData.Ret == -3002 {
		time.Sleep(1 * time.Second)
		fmt.Println("正在重试.....")
		return n.result(cid)
	}

	fmt.Println(string(buf))

	return revData.Text
}

func (n *Yundama) ReportErr() {

	//n.report("1")
}

func (n *Yundama) ReportSuc() {
	//n.report("0")
}

func (n *Yundama) report(flag string) {
	rs, err := http.PostForm("http://api.yundama.com/api.php?method=report", nil) //http.DefaultClient.Do(rq)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer rs.Body.Close()

	var retData *RetData

	bufRs, _ := ioutil.ReadAll(rs.Body)

	err = json.Unmarshal(bufRs, &retData)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("上报结果返回---%d\n", retData.Ret)

}
