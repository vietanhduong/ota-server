package storage_object

import (
	"errors"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"net/http"
	"path/filepath"
	"strings"
)

func ValidateExtension(filename string) error {
	// extract file extension
	ext := strings.ToLower(filepath.Ext(filename))

	// just accept `.ipa` and `.plist`
	if ext == IpaExt || ext == PlistExt {
		return nil
	}
	return cerrors.NewCError(http.StatusBadRequest, errors.New("only 'plist' and 'ipa' extensions are accepted"))
}
