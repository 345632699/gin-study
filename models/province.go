package models

import "github.com/jinzhu/gorm"

type Province struct {
	Model

	Province string `json:"province"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Img      string `json:"img"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
}

// GetArticles gets a list of articles based on paging constraints
func GetProvince(pageNum int, pageSize int, query []string) ([]Province, error) {
	var provinces []Province
	//quert := []string{"广东省","北京市"}
	//err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&provinces).Error
	err := db.Where("province in (?)", query).Find(&provinces).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return provinces, nil
}

func GetProvinceByName(name string) (Province, error) {
	var province Province
	err := db.Where("province = ? ", name).First(&province).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return province, err
	}
	return province, nil
}
