package stat

import (
	log "github.com/sirupsen/logrus"
)

func PingBackends() bool {
	var success bool = true
	log.Debug("Pinging backends")
	return success
}
