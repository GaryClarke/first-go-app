// File: internal/request/validate_test.go
package request

import "testing"

func TestValidateFullBookRequest_ValidInput(t *testing.T) {
	// Create FullBookRequest br
	br := FullBookRequest{
		Title:  "Valid Book",
		Author: "Valid Author",
		Year:   1999,
	}

	// errors := ValidateFullBookRequest(br)
	errors := ValidateFullBookRequest(&br)

	// Check errors is empty (len == 0)
	if len(errors) > 0 {
		t.Errorf("expected no validation errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateFullBookRequest_InvalidInput(t *testing.T) {
	// Table-driven tests: we define a list (slice) of test cases to loop over.
	tests := []struct {
		name     string          // A short label for the test case
		br       FullBookRequest // The input data to validate
		wantKeys []string        // The expected error keys we should get back
	}{
		{
			name:     "missing all fields",
			br:       FullBookRequest{},
			wantKeys: []string{"title", "author", "year"},
		},
		{
			name: "missing title",
			br: FullBookRequest{
				Author: "Valid Author", // Valid author
				Year:   1999,           // Valid year
			},
			wantKeys: []string{"title"}, // Only title should fail validation
		},
	}
}
