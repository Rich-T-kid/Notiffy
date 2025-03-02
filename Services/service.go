package services

import "context"

// This is for the broad Push notifications. So notifying Users all at once
// Keep minimal
//

// more specific name that string
type Tag string

const (
	TagSMS   = "SMS"
	TagEmail = "Email"
)

type Filter func(ctx context.Context) bool // way to filter out who to notify
type tags []Tag                            // assign user to sub categories, this can be used later for filtering

type UserService interface {
	// Register adds a user to the service and associates them with the provided subcategory.
	// If the user is already registered, they will be added to the new subcategory without duplication.
	Register(ctx context.Context, userid string, subcategory []Tag) error

	// Unregister removes a user from the service.
	// If stripCategory is false, the user is fully unregistered from the service.
	// If stripCategory is true, the user is only removed from the specified subcategories.
	// Multiple subcategories can be provided as variadic arguments.
	Unregister(ctx context.Context, userid string, stripCategory bool, subcategories []Tag) error

	// ListUsers retrieves a list of user IDs that match the provided filter criteria.
	// The filter function can be used to filter users based on tags or other conditions.
	ListUsers(filter Filter) ([]string, error)
}

type NotificationService interface {
	Start(ctx context.Context) error
	// start service, I.E set up any related API and check service is up to run
	UserService
	Notify(ctx context.Context, body Messenger, filter Filter) (int, error)                       // send message body out to the users who fit the criterial
	SendDirectMessage(ctx context.Context, body Messenger, from string, recipient []string) error // is directly to another user/service

	Validate() (bool, error) // Implentation specifc validation that checks weather the current notifcation is good or not
}

type Messenger interface {
	Message() MessageBody  // Returns the message body
	Metadata() MessageMeta // Optional metadata for the message
	Type() string          // return types like "email", "sms", "push"
}

type MessageBody interface {
	Content() interface{} // Body of the message. I.E text body, email body
}

type MessageMeta interface {
	Tags() tags       // Optional tags for filtering
	Priority() int    // Priority of the message (e.g., 1 = High, 5 = Low)
	Timestamp() int64 // Unix timestamp for message scheduling
	Title() string
}

/*
TO be done later dont go crazy implentating stuff youll never need. start with the minimal first
type Notifyer interface {
	Notify()
}
type UserNotification interface {
	DirectMessage
	Scheduler
}

type DirectMessage interface {
	Notify(ctx context.Context, message Messenger, recipients ...string) error
	Validate(message Messenger) bool // Check if the message is suitable for this channel
}

type Scheduler interface {
	Schedule(ctx context.Context, method Notifyer, message Messenger, timestamp int64) error
	Cancel(ctx context.Context, messageID string) error
	ListScheduled(ctx context.Context) ([]Messenger, error)
}
*/
