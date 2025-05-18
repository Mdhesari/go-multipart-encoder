package multipart

import (
	"strings"
	"testing"
)

// TestRequestStruct is a sample struct for testing the EncodeMultipart function
type TestRequestStruct struct {
	Name        string `form:"name"`
	Email       string
	Age         int `form:"user_age"`
	IsActive    bool
	IgnoreField string `form:"-"`
	ProfilePic  []byte `form:"avatar" filename:"profile.jpg"`
	Tags        []string
	Numbers     []int
	UserInfo    struct {
		Address string
		Phone   string
	}
}

func TestEncodeMultipart(t *testing.T) {
	req := TestRequestStruct{
		Name:        "John Doe",
		Email:       "john@example.com",
		Age:         30,
		IsActive:    true,
		IgnoreField: "This should be ignored",
		ProfilePic:  []byte("fake image data"),
		Tags:        []string{"golang", "multipart"},
		Numbers:     []int{1, 2, 3},
		UserInfo: struct {
			Address string
			Phone   string
		}{
			Address: "123 Main St",
			Phone:   "555-1234",
		},
	}

	buf, contentType, err := EncodeMultipart(req)
	if err != nil {
		t.Fatalf("EncodeMultipart returned an error: %v", err)
	}

	if buf == nil {
		t.Fatal("EncodeMultipart returned a nil buffer")
	}

	if contentType == "" {
		t.Fatal("EncodeMultipart returned an empty content type")
	}

	// The multipart boundary is random, so we can't check the exact output
	// But we can check that it contains expected substrings
	output := buf.String()

	// Check form fields
	if !contains(output, `name="name"`) || !contains(output, "John Doe") {
		t.Error("Missing or incorrect 'name' field")
	}

	if !contains(output, `name="email"`) || !contains(output, "john@example.com") {
		t.Error("Missing or incorrect 'email' field")
	}

	if !contains(output, `name="user_age"`) || !contains(output, "30") {
		t.Error("Missing or incorrect 'user_age' field")
	}

	if !contains(output, `name="isactive"`) || !contains(output, "true") {
		t.Error("Missing or incorrect 'isactive' field")
	}

	// Check that ignored field is actually ignored
	if contains(output, "This should be ignored") {
		t.Error("Field with form:\"-\" tag was not ignored")
	}

	// Check file field
	if !contains(output, `name="avatar"`) || !contains(output, `filename="profile.jpg"`) {
		t.Error("Missing or incorrect file field")
	}

	// Check slice fields
	if !contains(output, `name="tags"`) || !contains(output, "golang") || !contains(output, "multipart") {
		t.Error("Missing or incorrect 'tags' field")
	}

	if !contains(output, `name="numbers"`) || !contains(output, "1") || !contains(output, "2") || !contains(output, "3") {
		t.Error("Missing or incorrect 'numbers' field")
	}

	// Check struct field (JSON-encoded)
	if !contains(output, `name="userinfo"`) || !contains(output, "123 Main St") || !contains(output, "555-1234") {
		t.Error("Missing or incorrect 'userinfo' field")
	}
}

func TestIODetectExtension(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "PNG file",
			data:     []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expected: ".png",
		},
		{
			name:     "JPEG file",
			data:     []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46},
			expected: ".jpg",
		},
		{
			name:     "GIF file",
			data:     []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x0A, 0x00},
			expected: ".gif",
		},
		{
			name:     "PDF file",
			data:     []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35},
			expected: ".pdf",
		},
		{
			name:     "Unknown file type",
			data:     []byte{0x00, 0x01, 0x02, 0x03},
			expected: ".bin",
		},
		{
			name:     "Too short",
			data:     []byte{0x00, 0x01},
			expected: ".bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getExtensionFromContent(tt.data)
			if result != tt.expected {
				t.Errorf("IODetectExtension() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
