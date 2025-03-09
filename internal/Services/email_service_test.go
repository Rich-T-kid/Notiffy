package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test for successful registration
func TestSuccessfulRegister(t *testing.T) {
	ctx := context.TODO()
	mailer := NewMailer()

	uniqueName := fmt.Sprintf("John_%d", time.Now().UnixNano())
	user := &EmailReigisterInfo{Name: uniqueName, Email: "john@example.com"}
	err := mailer.Register(ctx, user, []Tag{})

	assert.NoError(t, err)
	assert.True(t, mailer.exist(user.Name))
}

// Test registration with existing user
func TestRegisterExistingUser(t *testing.T) {
	ctx := context.TODO()
	mailer := NewMailer()

	uniqueName := fmt.Sprintf("John_%d", time.Now().UnixNano())
	user := &EmailReigisterInfo{Name: uniqueName, Email: "john@example.com"}
	err := mailer.Register(ctx, user, []Tag{})

	assert.NoError(t, err)

	err = mailer.Register(ctx, user, []Tag{})
	assert.Error(t, err)
	assert.Equal(t, ErrUsernameExists, err)
}

// Test registering multiple users
func TestRegisterMultipleUsers(t *testing.T) {
	ctx := context.TODO()
	mailer := NewMailer()

	for i := 0; i < 3; i++ {
		uniqueName := fmt.Sprintf("User_%d_%d", i, time.Now().UnixNano())
		user := &EmailReigisterInfo{Name: uniqueName, Email: fmt.Sprintf("user%d@example.com", i)}
		err := mailer.Register(ctx, user, []Tag{})
		assert.NoError(t, err)
		assert.True(t, mailer.exist(user.Name))
	}
}

// Test register, unregister, and check existence
func TestRegisterUnregisterAndCheck(t *testing.T) {
	ctx := context.TODO()
	mailer := NewMailer()

	uniqueName := fmt.Sprintf("John_%d", time.Now().UnixNano())
	user := &EmailReigisterInfo{Name: uniqueName, Email: "john@example.com"}
	err := mailer.Register(ctx, user, []Tag{})
	assert.NoError(t, err)

	err = mailer.Unregister(ctx, user, []Tag{})
	assert.NoError(t, err)
	assert.False(t, mailer.exist(user.Name))
}

// Test unregistering a non-existent user
func TestUnregisterNonExistentUserEmail(t *testing.T) {
	ctx := context.TODO()
	mailer := NewMailer()

	uniqueName := fmt.Sprintf("Ghost_%d", time.Now().UnixNano())
	user := &EmailReigisterInfo{Name: uniqueName, Email: "ghost@example.com"}

	err := mailer.Unregister(ctx, user, []Tag{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must exist")
}

// Test checking existence of a non-existent user
func TestExistenceOfNonExistentUser(t *testing.T) {
	mailer := NewMailer()

	assert.False(t, mailer.exist("NonExistentUser"))
}

func TestValidateEmail(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.io",
		"firstname.lastname@domain.co",
		"support@venmo.com",
		"noreply@redditmail.com",
		"notifications@codecrafters.discoursemail.com",
	}

	for _, email := range validEmails {
		assert.NoError(t, ValidateEmail(email), "Expected valid email: %s", email)
	}

}
