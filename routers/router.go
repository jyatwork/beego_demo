// @APIVersion 1.0.0
// @Title mobile API
// @Description mobile has every tool to get any job done, so codename for the new mobile APIs.
// @Contact astaxie@gmail.com
package routers

import (
	"apibeego/controllers"

	"github.com/astaxie/beego"
)

func init() {

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/addUser",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/health",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/creatUser",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
