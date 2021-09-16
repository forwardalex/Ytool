package aop

import (
	"context"
	"strings"
)

func getContext(point *JoinPoint) context.Context {
	if point.Params != nil && len(point.Params) > 0 {
		for _, param := range point.Params {
			switch param.Interface().(type) {
			case context.Context:
				ctx := param.Interface().(context.Context)

				return ctx
			default:
				continue
			}
		}
	}
	return nil
}

func getRequest(point *JoinPoint) interface{} {
	if point.Params != nil && len(point.Params) > 0 {

		for _, param := range point.Params {
			if strings.Contains(strings.ToUpper(param.Type().String()), strings.ToUpper(point.Method.Name)) {
				return param.Interface()
			}
		}
	}

	return nil
}
