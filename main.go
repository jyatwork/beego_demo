package main

import (
	"apibeego/models"
	_ "apibeego/routers"

	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//USER
type User struct {
	Id         int64 `orm:"column(uid);pk;auto;default(1)"`
	Psw        string
	Uname      string    `orm:"size(100);unique"`
	CreateTime time.Time `orm:"column(create_time);null"` //`orm:"index"`
	LastTime   time.Time `orm:"column(last_time);null"`   //`orm:"index"`
}

func init() {
	models.NewUserManager(&models.DBConfig{Host: beego.AppConfig.String("db_host"), Port: beego.AppConfig.String("db_port"), Database: beego.AppConfig.String("db_name"), Username: beego.AppConfig.String("db_user"), Password: beego.AppConfig.String("db_pass"), MaxIdleConns: beego.AppConfig.DefaultInt("db_max_idle_conn", 1), MaxOpenConns: beego.AppConfig.DefaultInt("db_max_open_conn", 10)})
	//

	orm.Debug = true //打开查询日志
}
func main() {
	//api自动化文档
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()
}
