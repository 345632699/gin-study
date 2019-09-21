package province_service

import (
	"github.com/EDDYCJY/go-gin-example/models"
)

type Province struct {
	ID         int
	Provinceid int
	Province   string
	Lat        string
	Lng        string
	X          int
	Y          int
	Img        string
	Name       string
	Desc       string

	PageNum  int
	PageSize int
}

func (t *Province) GetAll(query []string) ([]models.Province, error) {
	var (
		provinces []models.Province
	)

	provinces, err := models.GetProvince(t.PageNum, t.PageSize, query)
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func (t *Province) GetProvinceByName(name string) (models.Province, error) {
	var province models.Province
	province, err := models.GetProvinceByName(name)
	if err != nil {
		return province, err
	}

	return province, nil

}

func (t *Province) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}

	return maps
}
