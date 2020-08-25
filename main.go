//main.go
package main

import (
	"fmt"

	"fileshare/Config"
	"fileshare/Models"
	"fileshare/Routes"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var err error

func main() {
	Config.DB, err = gorm.Open("mysql", Config.DbURL(Config.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer Config.DB.Close()
	Config.DB.AutoMigrate(&Models.User{}, &Models.Attachment{})
	Config.DB.Model(&Models.Attachment{}).AddForeignKey("user_id", "user(id)", "RESTRICT", "RESTRICT")

	// //
	// os.Setenv("AWS_REGION", "eu-south-1")
	// os.Setenv("S3_BUCKET", "fileshare-app-test")
	// os.Setenv("S3_BUCKET", "fileshare-app-test")
	// os.Setenv("USER_FILE_PATH", "/uploads")
	gin.SetMode(gin.ReleaseMode)
	r := Routes.SetupRouter()
	//running
	r.Run()
}
