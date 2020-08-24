//Config/Database.go
package Config

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "host.docker.internal"
	}

	dbConfig := DBConfig{
		Host:     host,
		Port:     3306,
		User:     "root",
		Password: "my-secret-pw",
		DBName:   "fileshare",
	}
	return &dbConfig
}
func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=CET",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
