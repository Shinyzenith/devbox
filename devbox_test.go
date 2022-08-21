package main

import (
	"testing"
)

func TestShortnameResolver(t *testing.T) {
	if _, err := ResolveShortnameUrl("alpine", "latest"); err != nil {
		t.Errorf("alpine:latest shortname failed to resolve.")
	}

	if url, err := ResolveShortnameUrl("wrong_name", "wrong_tag"); err == nil {
		t.Errorf("Inaccurate image_name and image_tag succesfully resolved to %s", url)
	}
}
