package cerrors

import (
	"github.com/vietanhduong/ota-server/pkg/logger"
	"io"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		logger.Logger.Fatal(err)
	}
}
