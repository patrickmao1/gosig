package utils

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func LogExecTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Infof("%s took %s", name, elapsed)
}
