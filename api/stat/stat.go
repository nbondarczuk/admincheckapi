package stat

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"admincheckapi/api/resource"
)

const (
	ServiceAlive   int32 = 1
	ServiceHealthy int32 = 2
	ServiceError   int32 = 3
	ServiceDead    int32 = 4
)

var (
	alive, healthy int32
	requestId      string
)

//
// bToMb connverts bytes to Mbytes
//
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

//
// Info provides mem stat info
//
func Info() resource.Stat {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s := resource.Stat{
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC,
	}

	return s
}

func SetAlive(s int32) {
	var val = atomic.LoadInt32(&alive)
	log.Debugf("Alive status: %d -> %d", val, s)
	atomic.StoreInt32(&alive, s)
}

func IsAlive() bool {
	var val = atomic.LoadInt32(&alive)
	log.Debugf("Alive status: %d", val)
	return val == ServiceAlive
}

func SetHealthy(s int32) {
	var val = atomic.LoadInt32(&healthy)
	log.Debugf("Health status: %d -> %d", val, s)
	atomic.StoreInt32(&healthy, s)
}

func IsHealthy() bool {
	if PingBackends() == false {
		log.Errorf("Error pinging backends, service unhealthy")
		SetHealthy(ServiceError)
	} else {
		log.Debugf("Success pinging backends")
	}

	var val = atomic.LoadInt32(&healthy)
	log.Debugf("Health status: %d", val)
	return val == ServiceHealthy
}

func RequestId() string {
	requestId = fmt.Sprintf("%d", time.Now().UnixNano())
	log.Debugf("Current requestId: %s", requestId)
	return requestId
}
