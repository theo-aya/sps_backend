package models

import(
	"github.com/Qwerci/sps_backend/database"
)


func SyncDatabase() {
	err := database.DB.AutoMigrate(
		 &User{},&Contact{},&Message{},
	)

	if err != nil {
		panic("failed to migrate tables")
	}
}