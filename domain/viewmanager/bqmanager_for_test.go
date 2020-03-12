package viewmanager

import (
	"fmt"
	"time"
)

func SetTest() string {
	datasetPrefixForTest = fmt.Sprintf("test_%d_", time.Now().Nanosecond())
	return datasetPrefixForTest
}
