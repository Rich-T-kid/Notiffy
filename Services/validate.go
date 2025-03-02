package services

var validTags = func(input string) bool {
	return Tag(input) == TagSMS || Tag(input) == TagEmail
}

// Seems unconvetional but id like to keep all of the validation logic in one file
