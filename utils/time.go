package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func LogExecTime(start time.Time, msg string, format ...any) {
	elapsed := time.Since(start)
	log.Infof("%s took %s", fmt.Sprintf(msg, format...), elapsed)
}
