package viewmanager

import (
	"fmt"
	"time"
)

func SetTest() {
	datasetPrefixForTest = fmt.Sprintf("test_%d_", time.Now().Nanosecond())
}
