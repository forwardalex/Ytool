package log

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"net"
	"os"
	"time"
)

var (
	sed              = rand.NewSource(time.Now().Unix())
	localServiceName string
	localServiceIP   string
)

type key int

type Header struct {
	TraceID string `json:"trace_id"`
}

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// TraceIDKey TODO
	TraceIDKey     key    = 100
	logHeader      key    = 101
	serverName     key    = 102
	TraceStringKey string = "TraceID"
)

// GetLogHeader 获取log头信息
//func GetLogHeader(ctx context.Context) string {
//	timeStr := time.Now().Format("2006-01-02 15:04:05.000")
//	if headerStr, ok := ctx.Value(logHeader).(string); ok && headerStr != "" {
//		return "|report_time=" + timeStr + headerStr
//	}
//	sn := GetServiceName(ctx)
//	ip := getServiceIP()
//	header := GetHeader(ctx)
//	traceID, _ := GetTraceID(ctx)
//	// glog.Info("trace_id=", traceID, ", err=", err)
//
//	headerStr := "|service_name=" + sn + "|server_ip=" + ip + "|uin=" + strconv.FormatUint(header["uin"].(uint64), 10) +
//		"|trace_id=" + traceID + "|"
//	ctx = context.WithValue(ctx, logHeader, headerStr)
//	return "|report_time=" + timeStr + headerStr
//}
func getServiceIP() string {
	return localServiceIP
}

// GetTraceID return uin+seq  as an unique trace_id
func GetTraceID(ctx context.Context) (string, error) {
	// Read metadata from client.
	if ctx == nil {
		return "", errors.New("ctx is nil")
	}

	if tmpTraceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return tmpTraceID, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("GetHeader: failed to get metadata")
	}

	if headerStr, ok := md["header-bin"]; ok && len(headerStr) > 0 {

		return "", nil
	}
	return "", errors.New("no header-bin is founded")
}

// GetServiceName TODO
func GetServiceName(ctx context.Context) string {
	if localServiceName != "" {
		return localServiceName
	}
	sn := os.Getenv("service_name")
	if sn != "" {
		return sn
	}
	sn, _ = ctx.Value(serverName).(string)
	if sn == "" {
		sn = os.Getenv("POD_NAME")
	}

	if sn == "" {
		sn = os.Getenv("POD_IP")
	}
	return sn
}

func GetOutHeader(ctx context.Context) metadata.MD {
	if ctx == nil || &ctx == nil {
		return nil
	}

	// Read metadata from client.
	md, getok := metadata.FromOutgoingContext(ctx)
	if getok {
		return md
	} else {
		return nil
	}
}

func GetInComeHeader(ctx context.Context) metadata.MD {
	if ctx == nil || &ctx == nil {
		return nil
	}

	// Read metadata from client.
	md, getok := metadata.FromIncomingContext(ctx)
	if getok {
		return md
	} else {
		return nil
	}
}

// WriteHeader 填充head
func WriteHeader(ctx *context.Context, key string, value string) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("ctx is empty")
	}
	//*ctx = metadata.AppendToOutgoingContext(*ctx, "header-bin", value)
	md := metadata.Pairs(key, value)
	newctx := metadata.NewOutgoingContext(*ctx, md)
	return newctx, nil
}

func initServerIP() error {
	if podIP := os.Getenv("POD_IP"); podIP != "" {
		localServiceIP = podIP
		return nil
	}

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localServiceIP = ipnet.IP.String()
				return nil
			}
		}
	}
	return errors.New("no_local_ip_found")
}
