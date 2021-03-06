package middleware

import (
	"strings"
	"time"

	"github.com/kataras/iris/v12"
)

type RecoverErrorHandler func(ctx iris.Context, err interface{})

var (
	MaxLength      int                 = 1024 * 8
	RecoverHandler RecoverErrorHandler = func(ctx iris.Context, err interface{}) {
		ctx.Application().Logger().Error(err)
		ctx.StatusCode(iris.StatusInternalServerError)
	}

	MediaTypes StringArray = StringArray{
		"application/json",
		"text/xml",
		"application/xml",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
	}
)

func Debug(ctx iris.Context) {
	start := time.Now()
	logger := ctx.Application().Logger()

	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}

			if RecoverHandler != nil {
				RecoverHandler(ctx, err)
			}
			ctx.StopExecution()

			if MediaTypes.Contains(ctx.GetContentType(), false) >= 0 {
				end := time.Now()
				logger.Infof("请求协议[%v]路由[%v]，IP[%v]，Method[%v]，ContentType[%v]，返回：StatusCode[%v]，Body[%v]，执行耗时：[%v]ms\n", ctx.Request().Proto, ctx.Request().URL, ctx.RemoteAddr(), ctx.Method(), ctx.GetContentType(), ctx.GetStatusCode(), string(ctx.Recorder().Body()), end.Sub(start).Milliseconds())
			}
		}
	}()

	rBody, err := ctx.GetBody()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	ctx.Record()
	ctx.Next()

	if MediaTypes.Contains(ctx.GetContentType(), false) >= 0 {
		end := time.Now()
		l := len(rBody)
		if l > 0 && l < MaxLength {
			logger.Infof("请求协议[%v]路由[%v]，IP[%v]，Method[%v]，ContentType[%v]，Body[%v]，返回：StatusCode[%v]，Body[%v]，执行耗时：[%v]ms\n", ctx.Request().Proto, ctx.Request().URL, ctx.RemoteAddr(), ctx.Method(), ctx.GetContentType(), string(rBody), ctx.GetStatusCode(), string(ctx.Recorder().Body()), end.Sub(start).Milliseconds())
		} else {
			logger.Infof("请求协议[%v]路由[%v]，IP[%v]，Method[%v]，ContentType[%v]，返回：StatusCode[%v]，Body[%v]，执行耗时：[%v]ms\n", ctx.Request().Proto, ctx.Request().URL, ctx.RemoteAddr(), ctx.Method(), ctx.GetContentType(), ctx.GetStatusCode(), string(ctx.Recorder().Body()), end.Sub(start).Milliseconds())
		}
	}
}

type StringArray []string

func (array StringArray) Contains(val string, ignoreCase bool) (index int) {
	index = -1
	if ignoreCase {
		for i := 0; i < len(array); i++ {
			if strings.Contains(strings.ToLower(val), strings.ToLower(array[i])) {
				index = i
				return
			}
		}
	} else {
		for i := 0; i < len(array); i++ {
			if strings.Contains(val, array[i]) {
				index = i
				return
			}
		}
	}
	return
}
