package protobuff

import (
	context "context"
	"errors"
	"log"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"

	services "github.com/Rich-T-kid/Notiffy/internal/Services"
	"github.com/Rich-T-kid/Notiffy/pkg"
)

var (
	ErrNotImplemented = errors.New("This GRPC handler is not implemented")
)

// will implement the NotificationServiceClient interface
type GServer struct {
	live          bool // internal check to see if sever is running
	emailRegister services.UserService
	smsRegister   services.UserService
	emailNotifyer services.NotificationService
	smsNotifyer   services.NotificationService
	UnimplementedNotificationServiceServer
}

func NewGServer() *GServer {
	return &GServer{
		live:          true,
		emailRegister: services.NewMailRegister(),
		smsRegister:   services.NewSMSUserService(),
		emailNotifyer: services.NewMailer(),
		smsNotifyer:   services.NewSMSNotification(),
	}
}

// Simple health check method
func (gsrv *GServer) HealthCheck(ctx context.Context, in *empty.Empty) (*HealthCheckResponse, error) {
	requestID, _ := ctx.Value(pkg.RequestIDKey{}).(string)
	startTime, _ := ctx.Value(pkg.StartTime{}).(time.Time)

	log.Printf("Handling HealthCheck RequestID: %s, Started At: %s", requestID, startTime)

	return &HealthCheckResponse{
		Status: "Service is running",
	}, nil
}

// Service agnostic Methods
func (gsrv *GServer) ListUsers(ctx context.Context, in *Tags) (*UserListResponse, error) {
	return nil, ErrNotImplemented
}

// SMS methods
// registration
func (gsrv *GServer) RegisterSMS(ctx context.Context, in *SMSRegisterInfo) (*BasicResponse, error) {
	return nil, ErrNotImplemented
}
func (gsrv *GServer) UnregisterSMS(ctx context.Context, in *SMSRegisterInfo) (*BasicResponse, error) {
	return nil, ErrNotImplemented
}
func (gsrv *GServer) UpdateSMSRegistration(ctx context.Context, in *SMSRegisterInfo) (*BasicResponse, error) {
	return nil, ErrNotImplemented
}

// notification related
func (gsrv *GServer) SMSNotify(ctx context.Context, in *SMSNotifyRequest) (*NotifyResponse, error) {
	return nil, ErrNotImplemented
}

func (gsrv *GServer) SMSSendDirectMessage(ctx context.Context, in *SMSNotifyRequest) (*ErrorArray, error) {
	return nil, ErrNotImplemented
}

// Email Methods
func (gsrv *GServer) RegisterEmail(ctx context.Context, in *EmailRegisterInfo) (*BasicResponse, error) {
	return nil, ErrNotImplemented
}

func (gsrv *GServer) UnregisterEmail(ctx context.Context, in *EmailRegisterInfo) (*BasicResponse, error) {
	return nil, ErrNotImplemented
}

func (gsrv *GServer) UpdateEmailRegistration(ctx context.Context, in *EmailRegisterInfo) (*BasicResponse, error) {
	return nil, ErrNotImplemented
}

// notification related
func (gsrv *GServer) EmailNotify(ctx context.Context, in *EmailNotifyRequest) (*NotifyResponse, error) {
	return nil, ErrNotImplemented
}

func (gsrv *GServer) EmailSendDirectMessage(ctx context.Context, in *EmailSendDirectRequest) (*ErrorArray, error) {
	return nil, ErrNotImplemented
}
