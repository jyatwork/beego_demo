package models

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//MYSQL
type DBConfig struct {
	Host         string
	Port         string
	Database     string
	Username     string
	Password     string
	MaxIdleConns int //最大空闲连接
	MaxOpenConns int //最大连接数
}

//USER
type User struct {
	Id         int64 `orm:"column(uid);pk;auto;default(1)"`
	Psw        string
	Uname      string    `orm:"size(100);unique"`
	CreateTime time.Time `orm:"column(create_time);null"` //`orm:"index"`
	LastTime   time.Time `orm:"column(last_time);null"`   //`orm:"index"`
}

type UserManager struct {
	DBConf *DBConfig
}

func NewUserManager(dbConf *DBConfig) *UserManager {
	mgr := &UserManager{
		DBConf: dbConf,
	}
	mgr.initDB() //初始化orm
	return mgr
}

func (mgr *UserManager) initDB() {
	//registe driver
	orm.RegisterDriver("mysql", orm.DRMySQL)
	ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", mgr.DBConf.Username, mgr.DBConf.Password, mgr.DBConf.Host, mgr.DBConf.Port, mgr.DBConf.Database)
	log.Infof("datasource=[%s]", ds)
	//注册数据库
	err := orm.RegisterDataBase("default", "mysql", ds, mgr.DBConf.MaxIdleConns, mgr.DBConf.MaxOpenConns)
	if err != nil {
		log.Error(err)
	}
	//	orm.RegisterDataBase("A", "mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	//	orm.RegisterDataBase("B", "mysql", "root:root@tcp(127.0.0.1:3306)/test1?charset=utf8")

	//	orm.RegisterModel(new(User))

	//	orm.RunSyncdb("A", false, false)
	//	orm.RunSyncdb("B", false, false)
	// 注册数据库
	orm.RegisterModel(new(User))
	// 自动建表
	orm.RunSyncdb("default", false, true)
}

func (user *User) GetTableName() string {

	return "user"
}

func AddUser(user *User) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(user)
	log.Infof("Insert user:%+v", user)
	if err != nil {
		log.Error("AddUser error:", err.Error())
		return 0, err
	}
	return id, nil
}

func UpdateUserById(user *User) error {
	o := orm.NewOrm()
	v := User{Id: user.Id}
	// ascertain id exists in the database
	if err := o.Read(&v); err == nil {
		if num, err := o.Update(user, "Psw", "LastTime"); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}

	return nil
}

func GetUserByFilter(field string, user interface{}) ([]orm.Params, error) {

	var (
		num  int64
		maps []orm.Params
		err  error
	)
	o := orm.NewOrm()
	switch field {
	case "uname":
		num, err = o.Raw("select * FROM user WHERE uname=? ", user).Values(&maps)
	case "uid":

		num, err = o.Raw("select * FROM user WHERE uid=? ", user).Values(&maps)
	default:
		return nil, errors.New("field type error")

	}
	if err == nil && num > 0 {
		return maps, nil
	}
	return nil, err
}

func VerifyUser(user *User) ([]orm.Params, error) {
	//orm.Debug = true
	var (
		num  int64
		maps []orm.Params
		err  error
	)
	o := orm.NewOrm()

	if user.Id != 0 {
		num, err = o.Raw("select * FROM user WHERE uid=? and psw=? ", user.Id, user.Psw).Values(&maps)

	} else if user.Uname != "" {
		num, err = o.Raw("select * FROM user WHERE uname=? and psw=? ", user.Uname, user.Psw).Values(&maps)
	} else {
		return nil, errors.New("uid and uname not found")
	}
	//beego.Debug("row nums: ", num)
	if err == nil && num > 0 {
		return maps, nil
	}

	return nil, err
}

func DeleteUser(uid int64) (interface{}, error) {
	o := orm.NewOrm()
	user := User{Id: uid}
	_, err := o.Delete(&user)
	if err != nil {
		log.Error("QueryUser error:", err.Error())
		return nil, err
	}
	return user, nil
}
