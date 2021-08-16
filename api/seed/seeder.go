package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"../models"
)

var users = []models.User{
	models.User{
		Nickname: "Hyna Tester",
		Email:    "hyna1234@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Martin Brunch",
		Email:    "martin@gmail.com",
		Password: "password",
	},
}


func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Product{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Product{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	/*
		err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error: %v", err)
		}
	*/

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		/*
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}*/
	}
}
