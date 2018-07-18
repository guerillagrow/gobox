package routers

import (
	//"log"

	"github.com/guerillagrow/gobox/controllers"
	"github.com/guerillagrow/gobox/models"

	"github.com/guerillagrow/beego"
	"github.com/guerillagrow/beego/context"
	"github.com/guerillagrow/beego/plugins/auth"
)

func CustomAuthFilter(ctx *context.Context) {

	// !DEBUG new session id on every request

	userObj, ok := ctx.Input.Session("user").(models.User)
	//log.Println("Router.CustomAuthFilter() SID:", ctx.Input.CruSession.SessionID(), "user:", userObj)
	if userObj.ID < 1 || !ok {
		a := &auth.BasicAuth{Secrets: models.UserAuth, Realm: "GoBox"}
		email := a.CheckAuth(ctx.Request)
		if email == "" {
			a.RequireAuth(ctx.ResponseWriter, ctx.Request)
			return
		}

		user, _ := models.GetUserByEmail(email)

		ctx.Input.CruSession.Set("user", user)

		ctx.Input.CruSession.Set("user/id", user.ID)
		//log.Println("Router.CustomAuthFilter() SERR #1:", serr)
		ctx.Input.CruSession.Set("user/name", user.Name)
		//log.Println("Router.CustomAuthFilter() SERR #2:", serr)
		ctx.Input.CruSession.Set("user/email", user.Email)
		//log.Println("Router.CustomAuthFilter() SERR #3:", serr)
		ctx.Input.CruSession.Set("user/isadmin", user.IsAdmin)
		//log.Println("Router.CustomAuthFilter() SERR #4:", serr)

	}

}

func init() {

	//authPlugin := auth.NewBasicAuthenticator(models.UserAuth, "GoBox Backoffice")
	//beego.InsertFilter("*", beego.BeforeRouter, authPlugin)
	//beego.InsertFilter("*", beego.BeforeExec, authPlugin)
	beego.InsertFilter("*", beego.BeforeRouter, CustomAuthFilter)

	beego.Router("/", &controllers.DashboarController{})
	serviceNS := beego.NewNamespace("svc",
		beego.NSRouter("/sensors/temperature", &controllers.ServiceSensors{}, "get:GetTemp"),
		beego.NSRouter("/sensors/humidity", &controllers.ServiceSensors{}, "get:GetHumidity"),
		beego.NSRouter("/sys/time", &controllers.ServiceSys{}, "get:GetTime"),
		beego.NSRouter("/sys/pistats", &controllers.ServiceSys{}, "get:GetPiStats"),
		beego.NSRouter("/user", &controllers.ServiceUser{}),
		beego.NSRouter("/export", &controllers.ServiceExport{}),
		beego.NSRouter("/relay", &controllers.ServiceRelay{}),
		beego.NSRouter("/gobox", &controllers.ServiceGoBox{}))
	beego.AddNamespace(serviceNS)

	beego.SetStaticPath("/static", "static")
	beego.SetStaticPath("/export", "export")

	//beego.ErrorHandler("404", controllers.PageNotFound)
}
