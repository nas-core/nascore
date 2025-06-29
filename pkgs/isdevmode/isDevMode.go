package isdevmode

import "os"

func IsDevMode() bool {
	return os.Getenv("nascore_DEV_MODE") == "1" || os.Getenv("nascore_DEV_MODE") == "true"
}
