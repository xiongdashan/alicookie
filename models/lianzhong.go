package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/otwdev/galaxylib"
)

type Lianzhong struct {
	url    string
	RSData *LZRsData
	RQData *LZRqData
	//buf    []byte
}

type LZRqData struct {
	CaptchaData      string `json:"captchaData"`
	CaptchaMaxLength int    `json:"captchaMaxLength"`
	CaptchaMinLength int    `json:"captchaMinLength"`
	CaptchaType      int    `json:"captchaType"`
	Password         string `json:"password"`
	SoftwareID       int    `json:"softwareId"`
	SoftwareSecret   string `json:"softwareSecret"`
	Username         string `json:"username"`
	WorkerTipsID     int    `json:"workerTipsId"`
}

type LZRsData struct {
	Code int `json:"code"`
	Data struct {
		CaptchaID   string `json:"captchaId"`
		Recognition string `json:"recognition"`
	} `json:"data"`
	Message string `json:"message"`
	Ts      int    `json:"ts"`
}

func NewLianzhong(maxLen, mixLen, codetype int) *Lianzhong {
	lz := &Lianzhong{}
	lz.RQData = &LZRqData{
		CaptchaMaxLength: maxLen,
		CaptchaMinLength: mixLen,
	}
	cnf, _ := galaxylib.GalaxyCfgFile.GetSection("lianzhong")
	lz.RQData.Password = cnf["pwd"]
	lz.RQData.Username = cnf["user"]
	lz.RQData.SoftwareID = 12836
	lz.RQData.SoftwareSecret = "4AV9PNCHJfSkPplqKekZoeLXLiPkRyhkPr1LDkFR"
	lz.RQData.CaptchaType = codetype
	lz.url = cnf["api"]
	//lz.buf = buf
	return lz
}

func (l *Lianzhong) Get(buf []byte) string {
	l.RQData.CaptchaData = base64.StdEncoding.EncodeToString(buf)
	bufRS := l.rs("upload", l.RQData)
	var reData *LZRsData
	err := json.Unmarshal(bufRS, &reData)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	l.RSData = reData
	if reData.Code == 0 {
		return reData.Data.Recognition
	}
	fmt.Println(reData.Message)
	return ""
}

func (l *Lianzhong) ReportErr() {
	param := map[string]interface{}{
		"softwareId":     l.RQData.SoftwareID,
		"softwareSecret": l.RQData.SoftwareSecret,
		"username":       l.RQData.Username,
		"password":       l.RQData.Password,
		"captchaId":      l.RSData.Data.CaptchaID,
	}
	buf := l.rs("report-error", param)
	if buf == nil {
		return
	}
	fmt.Println(string(buf))
}

func (l *Lianzhong) rs(path string, param interface{}) []byte {
	uri := fmt.Sprintf("%s/%s", l.url, path)
	buf, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	rq, _ := http.NewRequest("POST", uri, bytes.NewBuffer(buf))
	rq.Header.Add("Content-type", "application/json")
	rq.Header.Add("Connection", "keep-alive")
	rs, err := http.DefaultClient.Do(rq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rs.Body.Close()

	bufRS, _ := ioutil.ReadAll(rs.Body)
	return bufRS
}
