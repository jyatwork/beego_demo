package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["apibeego/controllers:UserController"] = append(beego.GlobalControllerRouter["apibeego/controllers:UserController"],
		beego.ControllerComments{
			Method: "AddUser",
			Router: `/addUser`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apibeego/controllers:UserController"] = append(beego.GlobalControllerRouter["apibeego/controllers:UserController"],
		beego.ControllerComments{
			Method: "CreatUser",
			Router: `/creatUser`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["apibeego/controllers:UserController"] = append(beego.GlobalControllerRouter["apibeego/controllers:UserController"],
		beego.ControllerComments{
			Method: "Health",
			Router: `/health`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
