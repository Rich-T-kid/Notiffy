package pkg

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMissingRequestID = errors.New("")
	ErrMissingStartTime = errors.New("")
)

type RequestIDKey struct{}
type StartTime struct{}

// for http
func ContextWithRequestID() context.Context {
	requestUUID := uuid.NewString()
	ctx := context.Background()
	ctx = context.WithValue(ctx, RequestIDKey{}, requestUUID)
	ctx = context.WithValue(ctx, StartTime{}, time.Now())
	return ctx
}
func ConstantTimeFormat() string {
	return "02 Jan 06 15:04 MST"
}

/*
provides more context for debugging
Example use case:
	ctx := pkg.ContextWithRequestID()
	requestid := ctx.Value(pkg.RequestIDKey{}).(string)
	startTime := ctx.Value(pkg.StartTime{}).(time.Time)
	formatedTime := startTime.Format("Mon, 02 Jan 2006 15:04:05 MST")
	fmt.Printf("context id %s  context start time %s\n", requestid, formatedTime)
*/
