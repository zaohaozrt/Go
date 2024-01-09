package common

import (
	"fmt"
	"ginEssential/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database.err:" + err.Error())
	}
	db.AutoMigrate(&model.User{}) //自动建表
	return db

}
func GetDB() *gorm.DB {
	return DB
}
