package protobuff

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	logger "github.com/Rich-T-kid/Notiffy/internal/log"
	"github.com/Rich-T-kid/Notiffy/pkg"
)

// might not need this file
// for grpc
func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Generate a new Request ID
	requestID := uuid.NewString()

	// Add Request ID and Start Time to the context
	ctx = context.WithValue(ctx, pkg.RequestIDKey{}, requestID)
	ctx = context.WithValue(ctx, pkg.StartTime{}, time.Now())

	// TODO : remove the below line once it becomes overlogging
	// Optionally log the request details
	logger.Info(fmt.Sprintf("Received gRPC request: %s, RequestID: %s", info.FullMethod, requestID))

	// Continue to the actual gRPC method handler
	return handler(ctx, req)
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
