package W

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/S"
	"github.com/valyala/fasthttp"
)

type Posts struct {
	*fasthttp.Args
	M.SS
}

func (p *Posts) GetBool(key string) bool {
	val := p.SS[key]
	return !(val == `` || val == `0` || val == `f` || val == `false`)
}

func (p *Posts) GetJsonMap(key string) M.SX {
	return S.JsonToMap(p.GetStr(key))
}

func (p *Posts) IsSet(key string) bool {
	_, ok := p.SS[key]
	return ok
}

func (p *Posts) GetJsonStrArr(key string) []string {
	return S.JsonToStrArr(p.GetStr(key))
}

func (p *Posts) GetJsonObjArr(key string) []map[string]interface{} {
	return S.JsonToObjArr(p.GetStr(key))
}

func (p *Posts) GetJsonIntArr(key string) []int64 {
	return S.JsonToIntArr(p.GetStr(key))
}

func (p *Posts) FromContext(ctx *Context) {
	p.Args = ctx.RequestCtx.PostArgs()
	p.SS = M.SS{}
	for k, v := range p.Args.All() {
		p.SS[string(k)] = string(v)
	}
	mf, err := ctx.RequestCtx.MultipartForm()
	if err == nil {
		//L.Print(`Multipart Post Data: ` + I.ToStr(len(mf.Value)) + ` keys`)
		for k, v := range mf.Value {
			p.SS[k] = v[0]
			//L.Print(`* ` + k + `: ` + S.IfElse(len(v[0]) < 128, v[0], `length=`+I.ToStr(len(v[0]))) + S.If(len(v) > 1, ` [warning: array form ignored: `+I.ToStr(len(v))+`]`))
		}
	} else {
		L.Print(`Error Parsing Post Data: ` + err.Error())
	}
}

// censor the password string, also when length is too long
func (p *Posts) String() string {
	return p.SS.PrettyFunc(` | `, func(key, val string) string {
		if len(val) > 64 {
			return val[:64] + `...`
		}
		return S.IfElse(key == `pass` || key == `password`, S.Repeat(`*`, len(val)), val)
	})
}
func (p *Posts) NewlineString() string {
	return p.SS.PrettyFunc("\n\t", func(key, val string) string {
		if len(val) > 4096 {
			return val[:4096] + `...`
		}
		return S.IfElse(key == `pass` || key == `password`, S.Repeat(`*`, len(val)), val)
	})
}
