package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/context"
	//"github.com/muesli/cache2go"
	"github.com/satori/go.uuid"
)

// !TODO // !DEBUG

var CSRF *CSRFManager

var ERR_INVALID_CSRF_TOKEN error = errors.New("Invalid or expired CSRF token!")

type CSRFManager struct {
	//Cache *cache2go.CacheTable
}

func (c *CSRFManager) getTokenData(scope string, ctx *context.Context) (string, time.Time) {

	token, _ := ctx.Input.CruSession.Get(fmt.Sprintf("csrf/%s/token", scope)).(string)
	expire, _ := ctx.Input.CruSession.Get(fmt.Sprintf("csrf/%s/expire", scope)).(time.Time)

	return token, expire
}

func (c *CSRFManager) SetToken(scope string, lifetime time.Duration, ctx *context.Context) string {

	/*defer func() {
		ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter.ResponseWriter)
	}()*/

	if int(lifetime) == 0 {
		lifetime = (24 * 7) * time.Hour
	}

	etoken, _ := c.getTokenData(scope, ctx)

	if etoken != "" {
		fmt.Println("CSRFManager.SetToken() -> Update token:", etoken, "; scope:", scope, "SID:", ctx.Input.CruSession.SessionID())
		ctx.Input.CruSession.Set(fmt.Sprintf("csrf/%s/expire", scope), time.Now().Add(lifetime))
		return etoken
	}
	tid, _ := uuid.NewV4()
	token := tid.String()
	expire := time.Now().Add(lifetime)
	err := ctx.Input.CruSession.Set(fmt.Sprintf("csrf/%s/token", scope), token)
	fmt.Println(err)
	err = ctx.Input.CruSession.Set(fmt.Sprintf("csrf/%s/expire", scope), expire)
	fmt.Println(err)
	fmt.Println("CSRFManager.SetToken() -> New token:", token, "; scope:", scope, "SID:", ctx.Input.CruSession.SessionID())
	fmt.Println("------------------------------------------------------------------------------------------------------")

	return token
}

func (c *CSRFManager) ValidateToken(scope string, inputToken string, ctx *context.Context) error {

	token, expire := c.getTokenData(scope, ctx)
	if time.Now().After(expire) {
		return ERR_INVALID_CSRF_TOKEN
	}
	if token != inputToken {
		return ERR_INVALID_CSRF_TOKEN
	}
	return nil
}
