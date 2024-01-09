package main

import (
	"ginEssential/common"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var house = make(map[string]*Hub) //管理所有hub

func main() {
	InitConfig()
	common.DB = common.InitDB()
	defer common.DB.Close()
	r := gin.Default()
	r = CollectRoute(r)
	r.Run(":8888")
}
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
