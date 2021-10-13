package mime

import "testing"

func TestHeaders(t *testing.T) {
	if HeaderContentType != "Content-Type" {
		t.Error("Invalid content type")
	}

	if HeaderUserAgent != "User-Agent" {
		t.Error("Invalid content type")
	}

	if ApplicationTypeJSON != "application/json" {
		t.Error("Invalid content type")
	}

	if ApplicationTypeXML != "application/xml" {
		t.Error("Invalid content type")
	}
}