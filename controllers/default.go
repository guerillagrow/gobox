package controllers

import (
	//"net/http"

	"github.com/guerillagrow/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

/*func PageNotFound(rw http.ResponseWriter, r *http.Request) {
	//rw.WriteHeader(404)
}*/

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Prepare() {

	c.TplName = ""
	c.EnableRender = false // for json rendering!!!

}

func (c *ErrorController) Error404() {
	//c.Data["content"] = "page not found"
	//c.TplNames = "404.tpl"
	c.TplName = ""
	c.Data["json"] = ErrorResponse{
		Code:    404,
		Message: "Not found!",
	}
	c.ServeJSON()
}

func (c *ErrorController) Error403() {
	//c.Data["content"] = "page not found"
	//c.TplNames = "404.tpl"
	c.TplName = ""
	c.Data["json"] = ErrorResponse{
		Code:    403,
		Message: "Forbidden!",
	}
	c.ServeJSON()
}

func (c *ErrorController) Error401() {
	//c.Data["content"] = "page not found"
	//c.TplNames = "404.tpl"
	c.TplName = ""
	c.Data["json"] = ErrorResponse{
		Code:    401,
		Message: "Unauthorized!",
	}
	c.ServeJSON()
}

func (c *ErrorController) Error405() {
	//c.Data["content"] = "page not found"
	//c.TplNames = "404.tpl"
	c.TplName = ""
	c.Data["json"] = ErrorResponse{
		Code:    405,
		Message: "Method not allowed!",
	}
	c.ServeJSON()
}

func (c *ErrorController) Error408() {
	//c.Data["content"] = "page not found"
	//c.TplNames = "404.tpl"
	c.TplName = ""
	c.Data["json"] = ErrorResponse{
		Code:    408,
		Message: "Request timeout!",
	}
	c.ServeJSON()
}

func (c *ErrorController) Error500() {
	c.TplName = ""
	c.Data["content"] = "Internal server error!"
	c.Data["json"] = ErrorResponse{
		Code:    500,
		Message: "Internal server error!",
	}
	c.ServeJSON()
}

func (c *ErrorController) ErrorDb() {
	c.TplName = ""
	c.Data["json"] = ErrorResponse{
		Code:    4242,
		Message: "Database error!",
	}
	c.ServeJSON()
}
