package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//common validation logic in one file
// Seems unconvetional but id like to keep all of the validation logic in one file

// Validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

var validTLDs = map[string]bool{
	".com": true, ".net": true, ".org": true, ".edu": true,
	".gov": true, ".io": true, ".co": true, ".us": true,
}

// ValidateEmail performs a comprehensive email validation
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	if len(email) > 254 {
		return errors.New("email exceeds the maximum length of 254 characters")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("email does not match the required format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("email must contain a single @ character")
	}

	domain := parts[1]

	if len(domain) < 3 || len(domain) > 255 {
		return errors.New("domain part of the email is invalid")
	}

	// Validate the domain contains at least one dot
	if !strings.Contains(domain, ".") {
		return errors.New("domain must contain a dot (.)")
	}

	// Check if the domain has a valid TLD
	tld := strings.ToLower(domain[strings.LastIndex(domain, "."):])
	if !validTLDs[tld] {
		return fmt.Errorf("invalid top-level domain: %s", tld)
	}

	// Verify the domain has MX records
	//mxRecords, err := net.LookupMX(domain)
	//if err != nil || len(mxRecords) == 0 {
	//	return fmt.Errorf("domain does not have valid MX records: %s", domain)
	//}

	return nil
}
