package storage_object

import (
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"net/http"
	"path/filepath"
	"strings"
)

var AllowedExtensions = []string{".ipa", ".plist", ".apk", ".png", ".jpg", ".jpeg"}

func ValidateExtension(filename string) error {
	// extract file extension
	ext := strings.ToLower(filepath.Ext(filename))

	for _, e := range AllowedExtensions {
		if e == ext {
			return nil
		}
	}

	return cerrors.NewCError(http.StatusBadRequest, "only 'plist' and 'ipa' extensions are accepted")
}
