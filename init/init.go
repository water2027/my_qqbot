package init

import (
	"os"
	"github.com/joho/godotenv"
)

func init(){
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		godotenv.Load()
	}
}