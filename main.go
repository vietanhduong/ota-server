package main

import (
	"fmt"
	"github.com/vietanhduong/ota-server/pkg/app"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
)

func main() {
	application := app.App{}
	application.Initialize()
	application.Run(fmt.Sprintf(":%s", env.GetEnvAsStringOrFallback("PORT", "8080")))
}



