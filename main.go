package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/otwdev/galaxylib"
	"github.com/otwdev/getalipaycookie/models"
)

func main() {

	app := echo.New()

	app.POST("/api/alicookies", func(ctx echo.Context) error {
		param := &struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		if err := ctx.Bind(param); err != nil {
			return ctx.JSON(http.StatusOK, echo.Map{
				"ret": 1,
				"msg": err.Error(),
			})
		}

		var thirdCode models.IThirdCode

		codeName := galaxylib.GalaxyCfgFile.MustValue("platform", "name")

		if codeName == "lianzhong" {
			thirdCode = models.NewLianzhong(4, 4, 1001)
		} else {
			thirdCode = models.NewYundama("1004")
		}

		wd := models.NewWebDriverRq(thirdCode)

		cookies := wd.Rq(param.Username, param.Password)

		return ctx.JSON(http.StatusOK, echo.Map{
			"ret":  0,
			"data": cookies,
		})

	})

	port := galaxylib.GalaxyCfgFile.MustValue("host", "port")
	app.Start(port)
}
