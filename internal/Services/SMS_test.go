package services

// All Test cases Pass
import (
	"context"
	"testing"

	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
)

// MockValidator is a mock for the Validator interface
type MockValidator struct {
	valid bool
}

func (m *MockValidator) Validate(interface{}) (bool, error) {
	return m.valid, nil
}

func TestRegister(t *testing.T) {
	ctx := context.TODO()
	userInfo := &RegisterINFO{
		Name:    "Test User",
		Contact: 1234567890,
		Tags:    []Tag{TagSMS, TagEmail},
	}

	s := NewSMSNotification()

	err := s.Register(ctx, userInfo, []Tag{TagSMS})
	if err != nil {
		t.Errorf("Register test failed: %v", err)
	} else {

		t.Log("Register test passed")
	}
}

func TestUnregister(t *testing.T) {
	ctx := context.TODO()
	userInfo := &RegisterINFO{
		Name:    "Test User",
		Contact: 1234567890,
		Tags:    []Tag{TagSMS, TagEmail},
	}

	s := NewSMSNotification()

	err := s.Unregister(ctx, userInfo, []Tag{TagSMS})
	if err != nil {
		t.Errorf("Unregister test failed: %v", err)
	} else {
		t.Log("Unregister test passed")
	}
}

func TestUpdateRegistration(t *testing.T) {
	ctx := context.TODO()
	userInfo := &RegisterINFO{
		Name:    "Test User",
		Contact: 1234567890,
		Tags:    []Tag{TagSMS, TagEmail},
	}

	s := NewSMSNotification()

	err := s.UpdateRegistration(ctx, userInfo, []Tag{"SMS"})
	if err != nil {
		t.Errorf("UpdateRegistration test failed: %v", err)
	} else {
		t.Log("UpdateRegistration test passed")
	}
}

func TestExist(t *testing.T) {
	Username := "Eafv24v2avws"
	userInfo := &RegisterINFO{
		Name:    Username,
		Contact: 1234567890,
		Tags:    []Tag{TagSMS, TagEmail},
	}

	s := NewSMSNotification()

	_ = s.Register(context.TODO(), userInfo, []Tag{TagSMS})

	if !s.exist(Username) {
		t.Error("Exist test failed: Expected true for existing user")
	} else {
		t.Log("Exist test passed")
	}
}

func TestRegisterWithInvalidValidator(t *testing.T) {
	ctx := context.TODO()
	invalidValidator := &MockValidator{valid: false}

	s := NewSMSNotification()

	err := s.Register(ctx, invalidValidator, []Tag{TagSMS})
	if err == nil {
		t.Error("Expected error when registering with invalid validator, got none")
	} else {
		t.Log("Register with invalid validator test passed")
	}
}

func TestUnregisterNonExistentUser(t *testing.T) {
	ctx := context.TODO()
	userInfo := &RegisterINFO{
		Name:    "NonExistent User",
		Contact: 9876543210,
		Tags:    []Tag{TagSMS},
	}

	s := NewSMSNotification()

	err := s.Unregister(ctx, userInfo, []Tag{TagSMS})
	if err != nil {
		t.Log("Unregister non-existent user test passed")
	} else {
		t.Errorf("Unregister non-existent user test failed: %v", err)
	}
}
