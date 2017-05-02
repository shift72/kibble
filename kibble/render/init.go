package render

import (
	logging "github.com/op/go-logging"
)

var log *logging.Logger

func init() {
	log = logging.MustGetLogger("render")
}
