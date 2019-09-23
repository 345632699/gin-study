package api

import (
	"fmt"
	models2 "gin-study/models"
	"gin-study/service/province_service"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/qrcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type provinceData struct {
	Province []string `json:"province"`
}

type provinceJson struct {
	Name string `json:"name"`
}

const (
	QRCODE_URL = "https://github.com/EDDYCJY/blog#gin%E7%B3%BB%E5%88%97%E7%9B%AE%E5%BD%95"
)

func CreatePoster(c *gin.Context) {
	var provinceList provinceData
	appG := app.Gin{C: c}
	err := c.BindJSON(&provinceList)
	if err != nil {
		appG.Response(500, e.ERROR, err)
	}

	for _, province := range provinceList.Province {
		fmt.Println(province)
	}

	provinceService := province_service.Province{}
	provinces, err := provinceService.GetAll(provinceList.Province)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	fmt.Println(provinces)

	qr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	unix := time.Now().Unix()
	formatInt := strconv.FormatInt(unix, 10)
	posterName := province_service.GetPosterFlag() + "-" + formatInt + qr.GetQrCodeExt()
	mapPoster := province_service.NewMapPoster(posterName, qr)
	mapPosterBgService := province_service.NewMapPosterBg(
		"bg3x.png",
		mapPoster,
		&province_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 375 * 3,
			Y1: 667 * 3,
		},
		&province_service.Pt{
			X: 100 * 3,
			Y: 165 * 3,
		},
		provinces,
	)

	_, filePath, err := mapPosterBgService.Generate()
	if err != nil {
		fmt.Println(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})

}

func GetProvince(c *gin.Context) {
	var province provinceJson
	appG := app.Gin{C: c}
	err := c.BindJSON(&province)
	if err != nil {
		appG.Response(500, e.ERROR, err)
	}
	provinceService := province_service.Province{}
	provinceResult, err := provinceService.GetProvinceByName(province.Name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]models2.Province{
		"province": provinceResult,
	})
}
