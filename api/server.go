package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"./controllers"
	"./seed"
)

var server = controllers.Server{}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	if os.Getenv("SEED_LOAD") == "1"{
		seed.Load(server.DB)
	}

	if os.Getenv("IMG_DIR") == ""{
		log.Fatalf("Error Image_DIR root not found, %v", err)
	}

	if _, err := os.Stat(os.Getenv("IMG_DIR")+"/"+os.Getenv("IMG_DIR_PRODUCT")); os.IsNotExist(err) {
		err := os.Mkdir(os.Getenv("IMG_DIR")+"/"+os.Getenv("IMG_DIR_PRODUCT"), 0755)
		if err != nil {
			log.Fatalf("Error cant create Image_DIR for product IMAGE, %v", err)
		}
	}else{
		fmt.Println("Found IMG_DIR for product IMAGE")
	}

	server.Run(os.Getenv("WEB_PORT"))
}
