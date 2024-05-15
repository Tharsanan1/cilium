package xds

import (
	"log"
	"github.com/cilium/cilium/pkg/logging"
	"github.com/cilium/cilium/pkg/logging/logfields"
)

const (
	Subsys = "gateway-xds"
)
type Logger struct {
	Debug bool
}

var cilium_log = logging.DefaultLogger.WithField(logfields.LogSubsys, Subsys)

func (logger Logger) Debugf(format string, args ...interface{}) {
	if logger.Debug {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (logger Logger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO]"+format+"\n", args...)
}

func (logger Logger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format+"\n", args...)
}

func (logger Logger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR]"+format+"\n", args...)
}
