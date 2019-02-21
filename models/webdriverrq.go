package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/otwdev/galaxylib"
	"github.com/tebeka/selenium"
)

type WebDriverRq struct {
	wd        selenium.WebDriver
	yd        *Yundama
	codeRetry int
	thirdCode IThirdCode
}

func NewWebDriverRq(thirdCode IThirdCode) *WebDriverRq {
	return &WebDriverRq{
		thirdCode: thirdCode,
	}
}

func (w *WebDriverRq) Rq(user, pwd string) map[string]string {

	hostURL := galaxylib.GalaxyCfgFile.MustValue("selenium", "host")

	wd, err := selenium.NewRemote(selenium.Capabilities{"browserName": "chrome"}, hostURL)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer wd.Quit()

	w.wd = wd

	//w.wd.SwitchFrame

	err = wd.Get("https://auth.alipay.com/login/index.htm?bizFrom=mrchportal&goto=https%3A%2F%2Fenterpriseportal.alipay.com%2Findex.htm")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	time.Sleep(1 * time.Second)

	elm, err := wd.FindElement(selenium.ByID, "J-loginMethod-tabs")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	aryLi, _ := elm.FindElements(selenium.ByTagName, "li")
	aryLi[1].Click()

	time.Sleep(1 * time.Second)

	inputUser, err := wd.FindElement(selenium.ByID, "J-input-user")

	if !w.rendError(67, err) {
		return nil
	}

	inputUser.Click()
	inputUser.SendKeys(user)

	time.Sleep(1 * time.Second)

	w.inputPwd(pwd)

	if err := w.validateCode(); err != nil {
		return nil
	}

	loginBtn, err := wd.FindElement(selenium.ByID, "J-login-btn")

	if !w.rendError(85, err) {
		return nil
	}

	err = loginBtn.Click()

	time.Sleep(1 * time.Second)

	title, _ := wd.Title()

	for title == "登录 - 支付宝" {
		//w.screenshot("login")

		fmt.Println("再次尝试登录....")

		w.inputPwd(pwd)
		//time.Sleep(1 * time.Second)
		btn, _ := wd.FindElement(selenium.ByID, "J-login-btn")
		btn.Click()
		time.Sleep(1 * time.Second)
		title, _ = wd.Title()
		fmt.Println(title)
	}

	w.screenshot("submit")

	time.Sleep(1 * time.Second)

	cookies, err := wd.GetCookies()

	if !w.rendError(96, err) {
		return nil
	}

	rev := make(map[string]string)

	for _, v := range cookies {
		fmt.Printf("%s----%s\n", v.Name, v.Value)
		rev[v.Name] = v.Value
	}

	shotBuf, _ := wd.Screenshot()

	ioutil.WriteFile("./img/main.png", shotBuf, os.ModeType)

	return rev
}

func (w *WebDriverRq) validateCode() error {

	codePanel, _ := w.wd.FindElement(selenium.ByID, "J-checkcode")
	showcode, _ := codePanel.IsDisplayed()
	fmt.Println(showcode)

	if showcode {

		fmt.Println("获取验证码.....")
		code := w.getCode()

		inputCode, err := w.wd.FindElement(selenium.ByID, "J-input-checkcode")
		if !w.rendError(95, err) {
			return err
		}

		inputCode.Click()
		inputCode.SendKeys(code)

		time.Sleep(1 * time.Second)

		_, err = w.wd.FindElement(selenium.ByClassName, "sl-checkcode-suc")

		if err != nil {
			fmt.Println(err)
			if w.codeRetry == 5 {
				return err
			}
			w.codeRetry++

			fmt.Printf("第%d次尝试验证码....\n", w.codeRetry)

			imgCode, _ := w.wd.FindElement(selenium.ByID, "J-checkcode-img")

			imgCode.Click()

			return w.validateCode()
		}
	}
	return nil
}

func (w *WebDriverRq) inputPwd(pwd string) {
	inputPwd, err := w.wd.FindElement(selenium.ByID, "password_rsainput")

	if !w.rendError(76, err) {
		return
	}

	inputPwd.Click()

	for _, v := range pwd {

		time.Sleep(100 * time.Millisecond)
		inputPwd.SendKeys(string(v))

	}

	time.Sleep(1 * time.Second)
}

func (w *WebDriverRq) getCode() string {

	buf, err := w.wd.Screenshot()

	if !w.rendError(103, err) {
		return ""
	}

	ioutil.WriteFile("./img/code.png", buf, os.ModeType)

	nc := NewCutImg()
	imgBuf := nc.Cut(buf)

	//return ""

	fmt.Println("请求打码....")
	return w.thirdCode.Get(imgBuf.Bytes())

}

func (w *WebDriverRq) screenshot(name string) {
	buf, _ := w.wd.Screenshot()
	ioutil.WriteFile(fmt.Sprintf("./img/%s.png", name), buf, os.ModeType)
}

func (w *WebDriverRq) printBody(wd selenium.WebDriver) {
	body, _ := wd.PageSource()

	fmt.Println(body)
}

func (w *WebDriverRq) rendError(sign int, err error) bool {
	if err == nil {
		return true
	}
	fmt.Printf("%d--%s\n", sign, err.Error())
	return false
}
