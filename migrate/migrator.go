package migrate

import (
	"Blog/models"
	"Blog/storage"
	"fmt"
)

func RunMigrations() {
	db := storage.DB

	err := db.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		panic("Ошибка миграции: " + err.Error())
	}

	fmt.Println("Миграция завершена")
}
