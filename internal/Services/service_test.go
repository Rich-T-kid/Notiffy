package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Rich-T-kid/Notiffy/pkg"
)

func TestServiceSMSNotification_Register(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	smsService := NewSMSService()
	var name = fmt.Sprintf("TestUser%d", time.Now().UnixNano())
	userInfo := &RegisterINFO{
		Name:    name,
		Contact: 9085252880,
		Tags:    []Tag{"SMS"},
	}

	// Test Registration
	err := smsService.Register(ctx, userInfo, []Tag{"SMS"})
	assert.NoError(t, err, "Expected no error while registering a user")

	// Test Duplicate Registration
	err = smsService.Register(ctx, userInfo, []Tag{"SMS"})
	assert.Error(t, err, "Expected error when registering a duplicate user")
}

func TestServiceSMSNotification_Notify(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	smsService := NewSMSService()
	body := DefineMessage("TestSender", "TestTitle", "TestMessage")

	// Always notify
	filter := func(ctx context.Context, input ...interface{}) bool {
		return true
	}
	// Test Notification
	notified, errors := smsService.Notify(ctx, body, filter)
	assert.LessOrEqual(t, len(errors), 1, "Expected no errors during notification")

	// Test Empty Notification
	emptyFilter := func(ctx context.Context, input ...interface{}) bool {
		return false
	}
	notified, _ = smsService.Notify(ctx, body, emptyFilter)
	assert.Equal(t, 0, notified, "Expected no users to be notified")
}

func TestServiceSMSNotification_Unregister(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	smsService := NewSMSService()
	userInfo := &RegisterINFO{
		Name:    "TestUser",
		Contact: 9085252880,
		Tags:    []Tag{"SMS"},
	}

	// Test Unregistration
	err := smsService.Unregister(ctx, userInfo, []Tag{"SMS"})
	assert.NoError(t, err, "Expected no error while unregistering a user")

	// Test Unregister Non-existent User
	nonExistentUser := &RegisterINFO{Name: "NonExistent", Contact: 1234567890}
	err = smsService.Unregister(ctx, nonExistentUser, []Tag{"SMS"})
	assert.Error(t, err, "Expected error when unregistering a non-existent user")
}

func TestServiceSMSNotification_UpdateRegistration(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	smsService := NewSMSService()
	userInfo := &RegisterINFO{
		Name:    "TestUser",
		Contact: 9085252880,
		Tags:    []Tag{"SMS"},
	}

	// Test Registration Update
	smsService.Register(ctx, userInfo, userInfo.Tags)
	err := smsService.UpdateRegistration(ctx, userInfo, []Tag{"Sports", "Dance"})
	assert.NoError(t, err, "Expected no error while updating registration")

	// Test Invalid Update
	invalidUserInfo := &RegisterINFO{Name: "InvalidUser", Contact: 0}
	err = smsService.UpdateRegistration(ctx, invalidUserInfo, []Tag{"Sports"})
	assert.Error(t, err, "Expected error with invalid user information")
}

func TestServiceSMSNotification_ListUsers(t *testing.T) {
	smsService := NewSMSService()
	filter := func(ctx context.Context, input ...interface{}) bool {
		return true
	}
	users, err := smsService.ListUsers(filter)
	assert.NoError(t, err, "Expected no error while listing users")
	assert.Greater(t, len(users), 0, "Expected to retrieve at least one user")
}

func TestServiceEmailNotification_Register(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	emailService := NewMailNotiffyer()
	var name = fmt.Sprintf("TestUser%d", time.Now().UnixNano())
	userInfo := &EmailReigisterInfo{
		Name:  name,
		Email: "richiebbaah@gmail.com",
		Tags:  []Tag{"Email"},
	}

	// Test Registration
	err := emailService.Register(ctx, userInfo, []Tag{"Email"})
	assert.NoError(t, err, "Expected no error while registering a user")

	// Test Duplicate Registration
	err = emailService.Register(ctx, userInfo, []Tag{"Email"})
	assert.Error(t, err, "Expected error when registering a duplicate user")
}

func TestServiceEmailNotification_Notify(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	emailService := NewMailNotiffyer()
	body := DefineMail("TestSubject", "TestBody", "richiebbaah@gmail.com", []Tag{"Email"})

	// Always notify
	filter := func(ctx context.Context, input ...interface{}) bool {
		return true
	}
	users, _ := emailService.ListUsers(filter)

	// Test Notification
	notified, errors := emailService.Notify(ctx, body, filter)
	assert.LessOrEqual(t, len(errors), 1, "Expected no errors during notification")
	if len(users) > 1 {
		assert.Greater(t, notified, 0, "Expected at least one user to be notified")
	}
	// Test Empty Notification
	emptyFilter := func(ctx context.Context, input ...interface{}) bool {
		return false
	}
	notified, errors = emailService.Notify(ctx, body, emptyFilter)
	assert.Equal(t, 0, notified, "Expected no users to be notified")
}

func TestServiceEmailNotification_Unregister(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	emailService := NewMailNotiffyer()
	userInfo := &EmailReigisterInfo{
		Name:  "TestUser",
		Email: "richiebbaah@gmail.com",
		Tags:  []Tag{"Email"},
	}

	// Test Unregistration
	err := emailService.Unregister(ctx, userInfo, []Tag{"Email"})
	assert.NoError(t, err, "Expected no error while unregistering a user")

	// Test Unregister Non-existent User
	nonExistentUser := &EmailReigisterInfo{Name: "NonExistent", Email: "nonexistent@example.com"}
	err = emailService.Unregister(ctx, nonExistentUser, []Tag{"Email"})
	assert.Error(t, err, "Expected error when unregistering a non-existent user")
}

func TestServiceEmailNotification_UpdateRegistration(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
	emailService := NewMailNotiffyer()
	userInfo := &EmailReigisterInfo{
		Name:  "TestUser",
		Email: "richiebbaah@gmail.com",
		Tags:  []Tag{"Email"},
	}

	// Test Registration Update
	emailService.Register(ctx, userInfo, userInfo.Tags)
	err := emailService.UpdateRegistration(ctx, userInfo, []Tag{"Sports", "Dance"})
	assert.NoError(t, err, "Expected no error while updating registration")

	// Test Invalid Update
	invalidUserInfo := &EmailReigisterInfo{Name: "InvalidUser", Email: "invalid@email.com"}
	err = emailService.UpdateRegistration(ctx, invalidUserInfo, []Tag{"Sports"})
	assert.Error(t, err, "Expected error with invalid user information")
}

func TestServiceEmailNotification_ListUsers(t *testing.T) {
	emailService := NewMailNotiffyer()
	filter := func(ctx context.Context, input ...interface{}) bool {
		return true
	}
	users, err := emailService.ListUsers(filter)
	assert.NoError(t, err, "Expected no error while listing users")
	assert.Greater(t, len(users), 0, "Expected to retrieve at least one user")
}
