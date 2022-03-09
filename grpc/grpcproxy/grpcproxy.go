package grpchttpproxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/forwardalex/Ytool/grpc/grpcproxy/grpcall"
	"github.com/forwardalex/Ytool/log"
	"github.com/forwardalex/Ytool/utils/jsonconv"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"net/http"
)

// ProxyInterceptor TODO
type ProxyInterceptor func(ctx context.Context, traceId string, r *http.Request, req *GrpcReq, body string,
	handler ProxyHandler) *HttpResp

// ProxyHandler TODO
type ProxyHandler func(ctx context.Context, traceId string, req *GrpcReq, body string, headers []string) *HttpResp

var (
	// grpcEnter grpcall客户端
	grpcEnter *grpcall.EngineHandler
)

// Start 服务启动
func Start(l net.Listener, prefix string, target string) error {
	return StartWithInterceptor(l, prefix, target, nil)
}

// StartWithInterceptor 有中间件功能的代理服务启动
func StartWithInterceptor(l net.Listener, prefix string, target string, interceptor ProxyInterceptor) error {

	var err error
	// 初始化grpc客户端（利用反射）
	var handler = DefaultEventHandler{}
	grpcEnter, err = grpcall.New(
		grpcall.SetHookHandler(&handler),
	)
	if err != nil {
		log.Info(context.TODO(), "Proxy started", "")
	}
	grpcall.SetKeepAliveTime(1000 * 24 * time.Hour)(grpcEnter)
	grpcEnter.SetMode(grpcall.ProtoReflectMode)
	grpcEnter.Init(target)

	// 注册HTTP Server
	http.HandleFunc("/"+prefix+"/", func(w http.ResponseWriter, r *http.Request) {

		resp := &HttpResp{
			RetCode: 0,
			BizCode: 0,
			Message: "success",
			Payload: NullStruct{},
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                                                                                                                                                                                                              //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,platform") //header的类型
		w.Header().Set("content-type", "application/json")                                                                                                                                                                                                                                                                              //返回数据格式是json
		if r.Header.Get("tencent-leakscan") == "TST(Tencent Security Team)" {                                                                                                                                                                                                                                                           // 过滤扫描
			respByte, _ := json.Marshal(resp)
			w.Write(respByte)
			return
		}

		// 解析请求路径
		greq, isKeepAlive, err := parseUrlToGrpc(prefix, r)
		if err != nil {
			log.Info(r.Context(), err.Error(), "")
			resp.RetCode = -1
			resp.Message = err.Error()
			respByte, _ := json.Marshal(resp)
			w.Write(respByte)
			return
		}

		if isKeepAlive { // 直接返回200
			resp.RetCode = 200
			resp.BizCode = 0
			respByte, _ := json.Marshal(resp)
			w.Write(respByte)
			return
		}

		// 获取请求参数
		body, _ := ioutil.ReadAll(r.Body)
		body, dataType, _, err := jsonparser.Get(body, "body")
		if dataType != jsonparser.Object || err != nil {
			resp.RetCode = -1
			resp.Message = "request body error"
			respByte, _ := json.Marshal(resp)
			w.Write(respByte)
			return
		}

		// 请求处理
		ctx := context.Background()
		traceId := uuid.New().String()
		ctx, _ = log.WriteHeader(&ctx, log.TraceStringKey, traceId)
		//log.Info(ctx," test trace ")
		if interceptor != nil {
			resp = interceptor(ctx, traceId, r, greq, string(body), BaseHandler)
		} else {
			traceidStr := fmt.Sprintf("traceid:%s", traceId)
			resp = BaseHandler(ctx, traceId, greq, string(body), []string{traceidStr})
		}

		// key自动转换下划线
		respByte, _ := json.Marshal(jsonconv.JsonSnakeCase{Value: RetResult{Body: resp}})
		w.Write(respByte)
	})

	httpS := &http.Server{
		Handler: nil,
	}
	httpS.Serve(l)
	//err = http.ListenAndServe(":"+fmt.Sprint(port), nil)
	if err != nil {
		fmt.Println("http proxy start failed")
		return err
	}
	return nil
}

// RetResult TODO
type RetResult struct {
	Body interface{} `json:"body"`
}

// BaseHandler 基础请求处理
func BaseHandler(ctx context.Context, traceId string, greq *GrpcReq, body string, headers []string) *HttpResp {
	resp := &HttpResp{
		RetCode: 0,
		BizCode: 0,
		Message: "success",
		Payload: NullStruct{},
		TraceID: traceId,
	}

	payload, dataType, _, err := jsonparser.Get([]byte(body), "payload")
	if dataType != jsonparser.Object || err != nil {
		resp.RetCode = -1
		resp.Message = "request payload error"
		return resp
	}

	res, err := grpcEnter.CallWithCtxAndHeaders(ctx, greq.Server+"."+greq.Service, greq.Method, string(payload), headers)
	if err != nil {
		bizErr := getErrInfo(err.Error())
		if bizErr == nil {
			resp.RetCode = -1
			resp.Message = err.Error()
		} else {
			resp.BizCode = bizErr.ErrCode
			resp.Message = bizErr.ErrMsg
		}
		return resp
	}

	json.Unmarshal([]byte(res.Data), &resp.Payload)
	return resp
}

// 解析业务错误信息
func getErrInfo(msg string) *BizResp {
	var resp BizResp
	err := json.Unmarshal([]byte(msg), &resp)
	if err != nil {
		return nil
	}
	return &resp
}

// GrpcReq TODO
type GrpcReq struct {
	Server  string
	Service string
	Method  string
}

// NullStruct TODO
type NullStruct struct {
}

// HttpResp TODO
type HttpResp struct {
	RetCode int         `json:"retcode"`  // 系统错误码
	BizCode int         `json:"bizcode"`  // 业务错误码
	TraceID string      `json:"trace_id"` // 请求ID
	Message string      `json:"message"`  // 错误消息
	Payload interface{} `json:"payload"`  // 接口内容
}

// BizResp 业务返回
type BizResp struct {
	ErrCode int    `json:"errcode"` // 业务错误码
	ErrMsg  string `json:"errmsg"`  // 错误消息
}

func parseUrlToGrpc(prefix string, r *http.Request) (*GrpcReq, bool, error) {
	originPath := r.URL.Path
	path := strings.Replace(originPath, "/"+prefix+"/", "", -1)
	path = strings.TrimSuffix(path, "/") // 去掉结尾/
	pathInfo := strings.Split(path, "/")
	lenPath := len(pathInfo)
	if lenPath == 1 { // 此处认为服务根路径，心跳请求
		return nil, true, nil
	}
	if lenPath != 3 {
		errmsg := fmt.Sprintf("grpc path parse failed！originPath[%s] path[%s] len[%d]", originPath, path, lenPath)
		return nil, false, errors.New(errmsg)
	}
	return &GrpcReq{
		Server:  pathInfo[0],
		Service: pathInfo[1],
		Method:  pathInfo[2],
	}, false, nil
}

// DefaultEventHandler TODO
type DefaultEventHandler struct {
	sendChan chan []byte
}

// OnReceiveData TODO
func (h *DefaultEventHandler) OnReceiveData(md metadata.MD, resp string, respErr error) {
}

// OnReceiveTrailers TODO
func (h *DefaultEventHandler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
}
