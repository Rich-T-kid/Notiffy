package protobuff

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"

	"github.com/Rich-T-kid/Notiffy/api/grpc/protobuff"
	services "github.com/Rich-T-kid/Notiffy/internal/Services"
	"github.com/Rich-T-kid/Notiffy/pkg"
)

// TODO: GRPC doesnt allow for error handling like htttp. need to return errs
var (
	ErrNotImplemented = errors.New("this grpc handler is not implemented")
)

const (
	SuccsefulOperation = "Succseful Operation"
	SuccsefulStatus    = 200
	FailedOperation    = "Failed Operation"
	FailedStatus       = 400
)

// will implement the NotificationServiceClient interface
type GServer struct {
	live          bool // internal check to see if sever is running
	emailRegister services.UserService
	smsRegister   services.UserService
	emailNotifyer services.NotificationService
	smsNotifyer   services.NotificationService
	protobuff.UnimplementedNotificationServiceServer
}

func NewGServer() *GServer {
	ctx := pkg.ContextWithRequestID()
	emailNotifyer := services.NewMailNotiffyer()
	smsNotifyer := services.NewSMSNotification()
	err := emailNotifyer.Start(ctx)
	if err != nil {
		panic(fmt.Sprintf("email service failed to start for grpc server with error : %v", err))
	}
	err = smsNotifyer.Start(ctx)
	if err != nil {
		panic(fmt.Sprintf("sms service failed to start for grpc server with error %v", err))
	}
	return &GServer{
		live:          true,
		emailRegister: services.NewMailRegister(),
		smsRegister:   services.NewSMSUserService(),
		emailNotifyer: emailNotifyer,
		smsNotifyer:   smsNotifyer,
	}
}

// Works
// Simple health check method
func (gsrv *GServer) HealthCheck(ctx context.Context, in *empty.Empty) (*protobuff.HealthCheckResponse, error) {
	requestID, _ := ctx.Value(pkg.RequestIDKey{}).(string)
	startTime, _ := ctx.Value(pkg.StartTime{}).(time.Time)

	log.Printf("Handling HealthCheck RequestID: %s, Started At: %s", requestID, startTime)
	timestamp := time.Now().Format(pkg.ConstantTimeFormat())
	return &protobuff.HealthCheckResponse{
		Status: fmt.Sprintf("Service is running  Timestamp:%s", timestamp),
	}, nil
}

// WORKS
// Service agnostic Methods
// Returns all usersName associated with the number of tags passed in
func (gsrv *GServer) ListUsers(ctx context.Context, in *protobuff.Tags) (*protobuff.UserListResponse, error) {
	fmt.Printf("Recieved grpc Request to return users with %v tags full object %v \n", in.Topics, in)
	filter := func(ctx context.Context, input ...interface{}) bool {
		// input is an array of tags
		usertags, ok := input[0].([]string)
		if !ok {
			givenType := fmt.Sprintf("%T", input[0])
			fmt.Printf("this filter expects a string but type %s was passed in instead \n", givenType)
			return false
		}
		//fmt.Printf("Checking if any tags in (usertags) %s exist in (passed in tags) %s\n", usertags, temptags)
		return existIn(in.Topics, usertags)
	}
	smsUsers, err := gsrv.smsNotifyer.ListUsers(filter)
	if err != nil {
		return nil, err
	}
	emailUsers, err := gsrv.emailNotifyer.ListUsers(filter)
	if err != nil {
		return nil, err
	}
	// pass to both email and sms services.
	smsUsers = append(smsUsers, emailUsers...)
	// join results

	// return results
	// this is just for code clarity renaming the joined results of email users and sms users into one  results array that is returned
	var results = smsUsers
	return &protobuff.UserListResponse{Users: results}, nil
}

// Works
// SMS methods
// registration
func (gsrv *GServer) RegisterSMS(ctx context.Context, in *protobuff.SMSRegisterInfo) (*protobuff.BasicResponse, error) {
	newTags := stringtoTags(in.Tags.Topics)
	smsr := services.NewSMSRegister(in.Name, in.ContactNumber, newTags)
	if err := gsrv.smsRegister.Register(ctx, smsr, smsr.Tags); err != nil {
		return &protobuff.BasicResponse{
			Message: FailedOperation,
			Status:  FailedStatus,
		}, err
	}
	return &protobuff.BasicResponse{
		Message: SuccsefulOperation,
		Status:  SuccsefulStatus,
	}, nil
}

// Works
func (gsrv *GServer) UnregisterSMS(ctx context.Context, in *protobuff.SMSRegisterInfo) (*protobuff.BasicResponse, error) {
	newTags := stringtoTags(in.Tags.Topics)
	smsr := services.NewSMSRegister(in.Name, in.ContactNumber, newTags)
	if err := gsrv.smsRegister.Unregister(ctx, smsr, smsr.Tags); err != nil {
		return &protobuff.BasicResponse{
			Message: FailedOperation,
			Status:  FailedStatus,
		}, err
	}
	return &protobuff.BasicResponse{
		Message: SuccsefulOperation,
		Status:  SuccsefulStatus,
	}, nil
}

// Works
func (gsrv *GServer) UpdateSMSRegistration(ctx context.Context, in *protobuff.SMSRegisterInfo) (*protobuff.BasicResponse, error) {
	newTags := stringtoTags(in.Tags.Topics)
	fmt.Printf("Raw object passed in %v tags passed in %s tags to string %s\n", in, in.Tags.Topics, newTags)
	smsr := services.NewSMSRegister(in.Name, in.ContactNumber, newTags)
	fmt.Printf("sms register object %+v\n", smsr)
	if err := gsrv.smsRegister.UpdateRegistration(ctx, smsr, newTags); err != nil {
		return &protobuff.BasicResponse{
			Message: FailedOperation,
			Status:  FailedStatus,
		}, err
	}
	return &protobuff.BasicResponse{
		Message: SuccsefulOperation,
		Status:  SuccsefulStatus,
	}, nil
}

// Works
// notification related
func (gsrv *GServer) SMSNotify(ctx context.Context, in *protobuff.SMSNotifyRequest) (*protobuff.NotifyResponse, error) {
	msg := services.DefineMessage("notiffy", in.Message.Title, in.Message.Message)
	filter := func(ctx context.Context, input ...interface{}) bool {
		usertags, ok := input[0].([]string)
		if !ok {
			givenType := fmt.Sprintf("%T", input[0])
			fmt.Printf("this filter expects a string but type %s was passed in instead \n", givenType)
			return false
		}
		return existIn(in.Tags.Topics, usertags)
	}
	n, errs := gsrv.smsNotifyer.Notify(ctx, msg, filter)
	return &protobuff.NotifyResponse{
		Notified: int64(n),
		Errors:   errorsToStrings(errs),
	}, nil
}

// Works
func (gsrv *GServer) SMSSendDirectMessage(ctx context.Context, in *protobuff.SMSSendDirectRequest) (*protobuff.ErrorArray, error) {
	msg := services.DefineMessage(in.From, in.Message.Title, in.Message.Message)

	errs := gsrv.smsNotifyer.SendDirectMessage(ctx, msg, in.From, in.Recipients)
	toString := errorsToStrings(errs)
	return &protobuff.ErrorArray{
		Errors: toString,
	}, nil
}

//Works

// Email Methods
func (gsrv *GServer) RegisterEmail(ctx context.Context, in *protobuff.EmailRegisterInfo) (*protobuff.BasicResponse, error) {
	newTags := stringtoTags(in.Tags.Topics)
	emailrgst := &services.EmailReigisterInfo{
		Name:  in.Name,
		Email: in.Email,
		Tags:  newTags,
	}
	if err := gsrv.emailRegister.Register(ctx, emailrgst, emailrgst.Tags); err != nil {
		return &protobuff.BasicResponse{
			Message: FailedOperation,
			Status:  FailedStatus,
		}, err
	}
	return &protobuff.BasicResponse{
		Message: SuccsefulOperation,
		Status:  SuccsefulStatus,
	}, nil
}

// Works
func (gsrv *GServer) UnregisterEmail(ctx context.Context, in *protobuff.EmailRegisterInfo) (*protobuff.BasicResponse, error) {
	newTags := stringtoTags(in.Tags.Topics)
	emailrgst := &services.EmailReigisterInfo{
		Name:  in.Name,
		Email: in.Email,
		Tags:  newTags,
	}
	if err := gsrv.emailRegister.Unregister(ctx, emailrgst, emailrgst.Tags); err != nil {
		return &protobuff.BasicResponse{
			Message: FailedOperation,
			Status:  FailedStatus,
		}, err
	}
	return &protobuff.BasicResponse{
		Message: SuccsefulOperation,
		Status:  SuccsefulStatus,
	}, nil
}

// Works
func (gsrv *GServer) UpdateEmailRegistration(ctx context.Context, in *protobuff.EmailRegisterInfo) (*protobuff.BasicResponse, error) {
	newTags := stringtoTags(in.Tags.Topics)
	fmt.Printf("Raw object passed in %v tags passed in %s tags to string %s\n", in, in.Tags.Topics, newTags)
	emailrgst := &services.EmailReigisterInfo{
		Name:  in.Name,
		Email: in.Email,
		Tags:  newTags,
	}
	fmt.Printf("emailRegist object %v\n", emailrgst)
	if err := gsrv.emailRegister.UpdateRegistration(ctx, emailrgst, emailrgst.Tags); err != nil {
		return &protobuff.BasicResponse{
			Message: FailedOperation,
			Status:  FailedStatus,
		}, err
	}
	return &protobuff.BasicResponse{
		Message: SuccsefulOperation,
		Status:  SuccsefulStatus,
	}, nil
}

// Works
// notification related
func (gsrv *GServer) EmailNotify(ctx context.Context, in *protobuff.EmailNotifyRequest) (*protobuff.NotifyResponse, error) {
	mail := services.DefineMail(in.Body.Subject, in.Body.Body, in.Body.To, stringtoTags(in.Body.Tags.Topics))
	filter := func(ctx context.Context, input ...interface{}) bool {
		usertags, ok := input[0].([]string)
		if !ok {
			givenType := fmt.Sprintf("%T", input[0])
			fmt.Printf("this filter expects a string array but type %s was passed in instead \n", givenType)
			return false
		}
		// Tags in passed in take precident over whats in the mail struct
		return existIn(in.Tags.Topics, usertags)
	}
	n, errs := gsrv.emailNotifyer.Notify(ctx, mail, filter)
	return &protobuff.NotifyResponse{
		Notified: int64(n),
		Errors:   errorsToStrings(errs),
	}, nil
}

// Works
func (gsrv *GServer) EmailSendDirectMessage(ctx context.Context, in *protobuff.EmailSendDirectRequest) (*protobuff.ErrorArray, error) {
	mail := services.DefineMail(in.Message.Subject, in.Message.Body, in.Message.To, stringtoTags(in.Message.Tags.Topics))
	errs := gsrv.emailNotifyer.SendDirectMessage(ctx, mail, in.From, in.Recipients)
	return &protobuff.ErrorArray{
		Errors: errorsToStrings(errs),
	}, nil
}

func stringtoTags(s []string) []services.Tag {
	var res []services.Tag
	for i := range s {
		item := services.Tag(s[i])
		res = append(res, item)
	}
	return res
}

func errorsToStrings(errors []error) []string {
	strings := make([]string, len(errors))
	for i, err := range errors {
		if err != nil {
			strings[i] = err.Error()
		} else {
			strings[i] = "" // Handle nil errors gracefully
		}
	}
	return strings
}

// Plz optimize. Leetcode prep will be worth somthing
// implement a Set dataStructure in pkg directory
func existIn[T comparable](input []T, elements []T) bool {
	for _, e := range elements { // Iterate over elements
		for _, v := range input { // Iterate over input
			if v == e { // If any element in input matches an element in elements
				return true
			}
		}
	}
	return false
}
