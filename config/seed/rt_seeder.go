package seed

import (
	"github/com/cl0ky/e-voting-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedRT(db *gorm.DB) error {
	var count int64
	db.Model(&models.RT{}).Count(&count)
	if count > 0 {
		return nil
	}

	rts := []models.RT{
		{Id: uuid.New(), Name: "RT 16", Region: "RW 08, Kelurahan Palmeriam, Kecamatan Matraman, Kota Jakarta Timur"},
		{Id: uuid.New(), Name: "RT 17", Region: "RW 08, Kelurahan Palmeriam, Kecamatan Matraman, Kota Jakarta Timur"},
		{Id: uuid.New(), Name: "RT 18", Region: "RW 08, Kelurahan Palmeriam, Kecamatan Matraman, Kota Jakarta Timur"},
		{Id: uuid.New(), Name: "RT 19", Region: "RW 08, Kelurahan Palmeriam, Kecamatan Matraman, Kota Jakarta Timur"},
		{Id: uuid.New(), Name: "RT 20", Region: "RW 08, Kelurahan Palmeriam, Kecamatan Matraman, Kota Jakarta Timur"},
	}
	for _, rt := range rts {
		if err := db.Create(&rt).Error; err != nil {
			return err
		}
	}
	return nil
}
