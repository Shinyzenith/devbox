/*
* SPDX-License-Identifier: GPL-3.0-only
*
* imageutil_test.go
*
* Created by:	Aakash Sen Sharma
* Copyright:	(C) 2022, Aakash Sen Sharma & Contributors
 */

package imageutil

import (
	"testing"
)

func TestShortnameResolver(t *testing.T) {
	if _, err := resolveShortnameUrl("alpine", "latest"); err != nil {
		t.Errorf("alpine:latest shortname failed to resolve")
	}

	if url, err := resolveShortnameUrl("wrong_name", "wrong_tag"); err == nil {
		t.Errorf("Inaccurate image_name and image_tag successfully resolved to %s", url)
	}
}
