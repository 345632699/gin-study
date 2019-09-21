package models

type VisitRecord struct {
	Model
	Ip string `json:"ip"`
}

func GetVisitRecordTotal() (int, error) {
	var count int
	if err := db.Model(&VisitRecord{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func AddVisitRecord(data map[string]interface{}) error {
	visitRecord := VisitRecord{
		Ip:         data["ip"].(string),
	}
	if err := db.Create(&visitRecord).Error; err != nil {
		return err
	}

	return nil
}