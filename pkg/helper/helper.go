package helper

import (
	"fmt"
	"time"
)

func GenerateUniqueID() string {
	return fmt.Sprintf("client-%d", time.Now().UnixNano())
}
