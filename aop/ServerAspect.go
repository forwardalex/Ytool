package aop

import (
	"github.com/forwardalex/Ytool/log"
)

type ServerAspect struct{}

func (a *ServerAspect) After(point *JoinPoint) {
	log.Info("after proxy", nil)

}
func (a *ServerAspect) Before(point *JoinPoint) (bool, error) {
	log.Info("proxy before", nil)
	return true, nil
}

func (a *ServerAspect) Finally(point *JoinPoint) {
	log.Info("proxy finally", nil)
}
