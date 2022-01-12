package aop

import (
	"context"
	"github.com/forwardalex/Ytool/log"
)

type ServerAspect struct{}

func (a *ServerAspect) After(point *JoinPoint) {
	log.Info(context.Background(), "after proxy", nil)

}
func (a *ServerAspect) Before(point *JoinPoint) (bool, error) {
	log.Info(context.Background(), "proxy before", nil)
	return true, nil
}

func (a *ServerAspect) Finally(point *JoinPoint) {
	log.Info(context.Background(), "proxy finally", nil)
}
