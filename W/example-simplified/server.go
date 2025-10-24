package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kokizzu/gotro/D/Rd"
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/X"
)

var VERSION string
var DEBUG_MODE = (VERSION == ``)
var LISTEN_ADDR = `:3001`
var ROOT_DIR string

var ASSETS = [][2]string{
	//// http://api.jquery.com/ 1.11.1
	{`js`, `jquery`},
}

var ROUTERS = map[string]W.Action{
	``:                            LoginExample,
	`login_example`:               LoginExample,
	`post_values_example`:         PostValuesExample,
	`named_params_example/:value`: NamedParamsExample,
	`query_string_example`:        QueryStringExample,
}

var WEBMASTER_EMAILS = M.SS{
	`test@test.com`: `Test`,
}

func AjaxResponse() W.Ajax {
	return W.Ajax{SX: M.SX{`is_success`: true}}
}

func LoginExample(ctx *W.Context) {
	user_id := ctx.Session.GetStr(`user_id`)
	if ctx.IsAjax() {
		ajax := AjaxResponse()
		posts := ctx.Posts()
		a := posts.GetStr(`a`)
		switch a {
		case `login`:
			username := posts.GetStr(`username`) // $_POST
			password := posts.GetStr(`password`)
			if username != password {
				ajax.Set(`is_success`, false)
				ajax.Error(`301 Wrong username or password; username is case sensitive`)
				ctx.AppendAjax(ajax)
				return
			}
			id := X.ToS(rand.Intn(1000))
			ctx.Session.Login(M.SX{
				`username`: username,
				`user_id`:  id,
				`level`:    M.SX{},
			})
			L.Describe(ctx.Session)
			ajax.Set(`logged`, id)
		case `logout`:
			ctx.Session.Logout()
		default:
			ajax.Error(`Unknown action`)
		}
		ctx.AppendAjax(ajax)
		return
	}
	ctx.Render(`login_example`, M.SX{
		`title`:   `Login Example`,
		`user_id`: user_id,
	})
}

// this function handles posted form values
func PostValuesExample(ctx *W.Context) {
	if ctx.IsAjax() {
		ajax := AjaxResponse()
		value := ctx.Posts().GetStr(`value`)
		ajax.Set(`value`, value)
		ctx.AppendAjax(ajax)
		return
	}
	ctx.Render(`post_values_example`, M.SX{
		`title`: `Post Values Example`,
	})
}

// this function handles /named_params_example/foo
func NamedParamsExample(ctx *W.Context) {
	ctx.Render(`named_params_example`, M.SX{
		`title`: `Named Params Example`,
		`value`: ctx.ParamStr(`value`),
	})
}

// this function handles /query_string_example?something=a&something_else=123
func QueryStringExample(ctx *W.Context) {
	params := ctx.QueryParams()
	//params.GetStr(`something`) // $_GET['something']
	//params.GetInt(`something_else`)
	data := M.SX{}
	for key, value := range params.All() {
		data.Set(X.ToS(key), X.ToS(value))
	}
	ctx.Render(`query_string_example`, M.SX{
		`title`: `Query String Example`,
		`data`:  data,
	})
}

const PROJECT_NAME = `Gotro Example`

// filter the page that may or may may not be accessed
func AuthFilter(ctx *W.Context) {
	L.Trace()
	handled := false
	if ctx.Session.GetStr(`user_id`) > `` {
		// logged in
	} else {
		// you can block the page for non-logged-in users here (handled=true)
	}
	ctx.Session.Touch()
	if !handled {
		cpu := L.PercentCPU()
		if cpu > 95.0 {
			W.Sessions.Inc(`throttle_counter`)
			if !ctx.IsAjax() {
				ctx.Error(503, `Server Overloaded`)
			} else {
				ctx.AppendString(`{"errors":["error 503: server is overloaded, please wait for a moment.."]}`)
			}
			fmt.Println(`Throttled: ` + F.ToS(cpu) + ` %`)
			return
		}
		ctx.Next()(ctx)
	}
}

// initialize loading time
func init() {
	var err error
	ROOT_DIR, err = os.Getwd() // filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		_, path, _, _ := runtime.Caller(0)
		slash_pos := strings.LastIndex(path, `/`) + 1
		ROOT_DIR = path[:slash_pos]
	} else {
		ROOT_DIR += `/`
	}
}

func main() {
	redis_conn := Rd.NewRedisSession(``, ``, 9, `session::`)
	global_conn := Rd.NewRedisSession(``, ``, 10, `session::`)
	W.InitSession(`Aaa`, 2*24*time.Hour, 1*24*time.Hour, *redis_conn, *global_conn)
	W.Mailers = map[string]*W.SmtpConfig{
		``: {
			Name:     `Mailer Daemon`,
			Username: `test.test`,
			Password: `123456`,
			Hostname: `smtp.gmail.com`,
			Port:     587,
		},
	}
	W.Assets = ASSETS
	W.Webmasters = WEBMASTER_EMAILS
	W.Routes = ROUTERS
	W.Filters = []W.Action{AuthFilter}
	// web engine
	server := W.NewEngine(DEBUG_MODE, false, PROJECT_NAME+VERSION, ROOT_DIR)
	server.StartServer(LISTEN_ADDR)
}
