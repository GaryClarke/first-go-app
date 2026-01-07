// File: internal/request/validate.go
package request

func ValidateFullBookRequest(br *FullBookRequest) map[string]string {
	// Make errors map to hold errors
	errors := make(map[string]string)

	// Validate title != ""
	if br.Title == "" {
		errors["title"] = "title is required"
	}

	// Validate author != ""
	if br.Author == "" {
		errors["author"] = "author is required"
	}

	// Validate year > 0
	if br.Year < 1 {
		errors["year"] = "year must be a positive integer"
	}

	// return errors map
	return errors
}
