package province_service

import (
	"gin-study/models"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"

	"github.com/EDDYCJY/go-gin-example/pkg/file"
	"github.com/EDDYCJY/go-gin-example/pkg/qrcode"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"image/png"
	"strconv"
)

type MapPoster struct {
	PosterName string
	Qr         *qrcode.QrCode
}

func NewMapPoster(posterName string, qr *qrcode.QrCode) *MapPoster {
	return &MapPoster{
		PosterName: posterName,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *MapPoster) CheckMergedImage(path string) bool {
	if file.CheckNotExist(path+a.PosterName) == true {
		return false
	}

	return true
}

func (a *MapPoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type MapPosterBg struct {
	Name string
	*MapPoster
	*Rect
	*Pt
	Provinces []models.Province
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

func NewMapPosterBg(name string, ap *MapPoster, rect *Rect, pt *Pt, provinces []models.Province) *MapPosterBg {
	return &MapPosterBg{
		Name:      name,
		MapPoster: ap,
		Rect:      rect,
		Pt:        pt,
		Provinces: provinces,
	}
}

type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	CountProvince int
	CountRecord   int
	ProvinceNames string

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

func (a *MapPosterBg) DrawPoster(d *DrawText, fontName string, regularFontName string) error {
	BoldFontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + fontName
	regularFontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + regularFontName
	boldFontSourceBytes, err := ioutil.ReadFile(BoldFontSource)
	if err != nil {
		return err
	}

	regularFontSourceBytes, err := ioutil.ReadFile(regularFontSource)
	if err != nil {
		return err
	}

	boldFont, err := freetype.ParseFont(boldFontSourceBytes)
	if err != nil {
		return err
	}

	regularFont, err := freetype.ParseFont(regularFontSourceBytes)
	if err != nil {
		return err
	}

	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(boldFont)
	fc.SetFontSize(d.Size0)
	fc.SetClip(d.JPG.Bounds())
	fc.SetDst(d.JPG)
	fc.SetSrc(image.Black)

	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}
	fc1 := freetype.NewContext()
	fc1.SetFontSize(d.Size1)
	fc1.SetFont(regularFont)
	fc1.SetDPI(72)
	fc1.SetClip(d.JPG.Bounds())
	fc1.SetDst(d.JPG)
	fc1.SetSrc(image.Black)
	_, err = fc1.DrawString("我是第"+strconv.Itoa(d.CountRecord+1)+"位点亮的建设者，我", freetype.Pt(d.X1, d.Y1))
	_, err = fc1.DrawString("在"+strconv.Itoa(d.CountProvince)+"（自治区、直辖市）留", freetype.Pt(d.X1, d.Y1+60))
	_, err = fc1.DrawString("下了足迹，见证超级工程的成长", freetype.Pt(d.X1, d.Y1+120))
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	// 插入记录
	models.AddLightRecord(d.ProvinceNames)

	return nil
}

func (a *MapPosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath()
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	if !a.CheckMergedImage(path) {
		mergedF, err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()

		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()

		flagF, err := file.MustOpen("flag3x.png", path)
		if err != nil {
			return "", "", err
		}
		defer flagF.Close()

		bgImage, err := png.Decode(bgF)
		if err != nil {
			return "", "", err
		}
		//qrImage, err := png.Decode(qrF)
		//if err != nil {
		//	return "", "", err
		//}

		flagImage, err := png.Decode(flagF)
		if err != nil {
			return "", "", err
		}

		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		//draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		provinces := a.Provinces
		provinceNames := ""
		for _, province := range provinces {
			x := province.X
			y := province.Y
			provinceNames += province.Province + ","
			draw.Draw(jpg, jpg.Bounds(), flagImage, flagImage.Bounds().Min.Sub(image.Pt(x*3, y*3)), draw.Over)
		}
		countProvince := len(provinces)
		countRecord, err := models.GetLightRecordTotal()
		if err != nil {
			return "", "", err
		}
		err = a.DrawPoster(&DrawText{
			JPG:           jpg,
			Merged:        mergedF,
			CountProvince: countProvince,
			ProvinceNames: provinceNames[0 : len(provinceNames)-1],
			CountRecord:   countRecord,
			Title:         "扫码参与",
			X0:            163 * 3,
			Y0:            547 * 3,
			Size0:         14 * 3,

			SubTitle: "我是第5位点亮的建设者，我在3个省留下了足迹，见证超级工程的成长",
			X1:       160 * 3,
			Y1:       566 * 3,
			Size1:    12 * 3,
		}, "msyhbd.ttf", "msyh.ttf")

		if err != nil {
			return "", "", err
		}
	}

	return fileName, path, nil
}
