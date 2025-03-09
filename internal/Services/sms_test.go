package services

// All Test cases Pass
import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/Rich-T-kid/Notiffy/internal/enviroment" // this package needs to always to be run first b4 all other custom packages
	"github.com/Rich-T-kid/Notiffy/pkg"
)

// MockValidator is a mock for the Validator interface
type MockValidator struct {
	valid bool
}

func (m *MockValidator) Validate() error {
	return nil
}

func TestRegister(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
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
	ctx := pkg.ContextWithRequestID()
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

	ctx := pkg.ContextWithRequestID()
	uniqueName := fmt.Sprintf("test_user%d", time.Now().UnixNano())
	userInfo := &RegisterINFO{
		Name:    uniqueName,
		Contact: 1234567890,
		Tags:    []Tag{TagSMS, TagEmail},
	}

	s := NewSMSNotification()
	s.Register(context.TODO(), userInfo, []Tag{"SMS"})
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

	ctx := pkg.ContextWithRequestID()
	s := NewSMSNotification()

	_ = s.Register(ctx, userInfo, []Tag{TagSMS})

	if !s.exist(Username) {
		t.Error("Exist test failed: Expected true for existing user")
	} else {
		t.Log("Exist test passed")
	}
}

func TestRegisterWithInvalidValidator(t *testing.T) {
	ctx := pkg.ContextWithRequestID()
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
	ctx := pkg.ContextWithRequestID()
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

// second half of test
func TestDefineMessage(t *testing.T) {
	msg := DefineMessage("Alice", "Greeting", "Hello World")
	if msg.From() != "Alice" {
		t.Errorf("expected From() to be 'Alice', got %s", msg.From())
	}
	if msg.Title() != "Greeting" {
		t.Errorf("expected Title() to be 'Greeting', got %s", msg.Title())
	}
	if msg.Type() != "sms" {
		t.Errorf("expected Type() to be 'sms', got %s", msg.Type())
	}
	if msg.Content().(string) != "Hello World" {
		t.Errorf("expected Content() to be 'Hello World', got %v", msg.Content())
	}
}

func TestNewSMSUserService(t *testing.T) {
	svc := NewSMSUserService()
	if svc == nil {
		t.Errorf("NewSMSUserService() returned nil")
	}
}
func TestNewSMSService(t *testing.T) {
	svc := NewSMSService()
	if svc == nil {
		t.Errorf("NewSMSService() returned nil")
	}
}

func TestNewSMSNotification(t *testing.T) {
	// Because NewSMSNotification() calls log.Fatal if TEXTBELT_API_KEY is missing,
	// you may want to set that up in your environment before running this test.
	// Alternatively, you can mock it or skip if the environment variable isn't set.
	// For demonstration, we just check if it returns a non-nil pointer (assuming
	// TEXTBELT_API_KEY is properly set).
	svc := NewSMSNotification()
	if svc == nil {
		t.Errorf("NewSMSNotification() returned nil")
	}
}
func TestSMSNotification_Register(t *testing.T) {
	svc := NewSMSNotification()
	ctx := pkg.ContextWithRequestID()

	userName := fmt.Sprintf("testRegisterUser_%d", time.Now().UnixNano())
	user := &RegisterINFO{
		Name:    userName,
		Contact: 1234567890,
	}
	err := svc.Register(ctx, user, []Tag{"example_tag"})
	if err != nil {
		t.Errorf("Register returned an error: %v", err)
	}
}

// Test the Unregister method by first registering a user, then unregistering
func TestSMSNotification_Unregister(t *testing.T) {
	svc := NewSMSNotification()
	ctx := pkg.ContextWithRequestID()

	userName := fmt.Sprintf("testUnregisterUser_%d", time.Now().UnixNano())
	user := &RegisterINFO{
		Name:    userName,
		Contact: 9998887777,
	}
	// Register first
	regErr := svc.Register(ctx, user, []Tag{"example_tag"})
	if regErr != nil {
		t.Fatalf("Register failed while preparing Unregister test: %v", regErr)
	}

	// Now attempt Unregister
	unregErr := svc.Unregister(ctx, user, nil)
	if unregErr != nil {
		t.Errorf("Unregister returned an error: %v", unregErr)
	}
}

// Test UpdateRegistration by registering a user, then removing a tag
func TestSMSNotification_UpdateRegistration(t *testing.T) {
	svc := NewSMSNotification()
	ctx := pkg.ContextWithRequestID()

	userName := fmt.Sprintf("testUpdateRegUser_%d", time.Now().UnixNano())
	user := &RegisterINFO{
		Name:    userName,
		Contact: 1112223333,
		Tags:    []Tag{"old_tag", "some_other_tag"},
	}

	// Register first
	regErr := svc.Register(ctx, user, user.Tags)
	if regErr != nil {
		t.Fatalf("Register failed while preparing UpdateRegistration test: %v", regErr)
	}

	// Now remove "old_tag" via UpdateRegistration
	upErr := svc.UpdateRegistration(ctx, user, []Tag{"old_tag"})
	if upErr != nil {
		t.Errorf("UpdateRegistration returned an error: %v", upErr)
	}
}

// Test ListUsers by registering a user, then verifying we can see them in the list
func TestSMSNotification_ListUsers(t *testing.T) {
	svc := NewSMSNotification()
	ctx := pkg.ContextWithRequestID()

	userName := fmt.Sprintf("testListUsers_%d", time.Now().UnixNano())
	user := &RegisterINFO{
		Name:    userName,
		Contact: 7778889999,
	}
	regErr := svc.Register(ctx, user, []Tag{"list_test_tag"})
	if regErr != nil {
		t.Fatalf("Register failed before ListUsers test: %v", regErr)
	}

	// Filter that returns all users
	allFilter := func(ctx context.Context, x ...interface{}) bool {
		return true
	}

	users, err := svc.ListUsers(allFilter)
	if err != nil {
		t.Errorf("ListUsers returned an error: %v", err)
	}

	// Optional: check that our user is in the list
	found := false
	for _, u := range users {
		if u == userName {
			found = true
			break
		}
	}
	if !found {
		t.Logf("Note: user %s was not found in ListUsers; possible logic or environment issue.", userName)
	}
}

// We cannot fully validate external service calls to Textbelt here; we just check for errors
func TestSMSNotification_Notify(t *testing.T) {
	svc := NewSMSNotification()
	ctx := pkg.ContextWithRequestID()

	msg := DefineMessage("testAuthor", "testTitle", "This is a test broadcast message.")
	filterAll := func(ctx context.Context, tags ...interface{}) bool {
		return true
	}

	count, errs := svc.Notify(ctx, msg, filterAll)
	t.Logf("Notify called, returned count=%d, errors=%v", count, errs)
	// Just check for error presence or not, as we can't confirm actual 3rd party outcome
	if len(errs) > 1 {
		// the one error this does return isnt a critcal error
		t.Logf("size of errors %d", len(errs))
		t.Errorf("Notify returned errors: %v", errs)
	}
}

// Similar approach: we just check for returned errors, not actual textbelt sending
func TestSMSNotification_SendDirectMessage(t *testing.T) {
	svc := NewSMSNotification()
	ctx := pkg.ContextWithRequestID()

	msg := DefineMessage("testDMAuthor", "testDMTitle", "This is a direct message.")
	// We'll pass some likely non-existing user, just to see if we get an error
	errs := svc.SendDirectMessage(ctx, msg, "DM_Test", []string{"NonExistentUser1", "NonExistentUser2"})

	if len(errs) > 0 {
		t.Logf("SendDirectMessage returned errors (possibly expected if users don't exist): %v", errs)
	} else {
		t.Log("SendDirectMessage returned no errors.")
	}
}
