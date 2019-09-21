package models

type LightRecord struct {
	Model
	Provinces string `json:"provinces"`
}

func GetLightRecordTotal() (int, error) {
	var count int
	if err := db.Model(&LightRecord{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func AddLightRecord(provinceNames string) error {
	lightRecord := LightRecord{
		Provinces:         provinceNames,
	}
	if err := db.Create(&lightRecord).Error; err != nil {
		return err
	}

	return nil
}