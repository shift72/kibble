package utils

import (
	"testing"

	"github.com/op/go-logging"

	"github.com/stretchr/testify/assert"
)

func TestWatchedLogging(t *testing.T) {

	unique := ConfigureWatchedLogging(logging.INFO)

	log.Critical("critical1") // logged
	log.Critical("critical1") // not uniuq
	log.Error("error1")       // logged
	log.Warning("warn1")      // logged
	log.Notice("notice1")     // skipped
	log.Info("info1")         // skipped

	assert.Equal(t, 3, len(unique.Logs()), "store")
}
