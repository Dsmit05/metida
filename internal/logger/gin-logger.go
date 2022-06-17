package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type request struct {
	Path      string
	RawQuery  string
	UserAgent string
	IP        string
	Email     string
	Role      string
}

func (r *request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("path", r.Path)
	enc.AddString("rawQuery", r.RawQuery)
	enc.AddString("userAgent", r.UserAgent)
	enc.AddString("ip", r.IP)
	enc.AddString("email", r.Email)
	enc.AddString("role", r.Role)
	return nil
}

type responseError struct {
	Status      int
	Code        int
	Description string
	Error       string
}

func (r *responseError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status", r.Status)
	enc.AddInt("code", r.Code)
	enc.AddString("description", r.Description)
	enc.AddString("error", r.Error)
	return nil
}

func RequestsError(c *gin.Context, status, code int, description string, err error) {
	newContext := c.Copy()
	email := newContext.GetString("email")
	role := newContext.GetString("role")
	path := newContext.Request.URL.Path
	rawQuery := newContext.Request.URL.RawQuery
	userAgent := newContext.Request.UserAgent()
	ip := newContext.ClientIP()

	req := &request{
		Path:      path,
		RawQuery:  rawQuery,
		UserAgent: userAgent,
		IP:        ip,
		Email:     email,
		Role:      role,
	}

	var textError string
	if err != nil {
		textError = err.Error()
	}

	resp := &responseError{
		Status:      status,
		Code:        code,
		Description: description,
		Error:       textError,
	}

	ZapLog.WithOptions(zap.WithCaller(false)).
		Error("api error", zap.Object("request", req), zap.Object("response", resp))
}

type responseSuccess struct {
	Status      int
	Code        int
	Description string
	Data        interface{}
}

func (r *responseSuccess) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status", r.Status)
	enc.AddInt("code", r.Code)
	enc.AddString("description", r.Description)
	enc.AddReflected("data", r.Data)
	return nil
}

func RequestsInfo(c *gin.Context, status, code int, description string, data interface{}) {
	newContext := c.Copy()
	email := newContext.GetString("email")
	role := newContext.GetString("role")
	path := newContext.Request.URL.Path
	rawQuery := newContext.Request.URL.RawQuery
	userAgent := newContext.Request.UserAgent()
	ip := newContext.ClientIP()

	req := &request{
		Path:      path,
		RawQuery:  rawQuery,
		UserAgent: userAgent,
		IP:        ip,
		Email:     email,
		Role:      role,
	}

	resp := &responseSuccess{
		Status:      status,
		Code:        code,
		Description: description,
		Data:        data,
	}

	ZapLog.WithOptions(zap.WithCaller(false)).
		Info("api info", zap.Object("request", req), zap.Object("response", resp))
}
