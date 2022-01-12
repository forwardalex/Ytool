package grpcall

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

var (
	callMaxTime = 60 * time.Second

	descSourceController = NewDescSourceEntry()
)

const (
	// ProtoSetMode TODO
	ProtoSetMode = iota // 0
	// ProtoFilesMode TODO
	ProtoFilesMode // 1
	// ProtoReflectMode TODO
	ProtoReflectMode // 2
)

// DescSourceEntry TODO
type DescSourceEntry struct {
	ctx    context.Context
	cancel context.CancelFunc

	descSource   DescriptorSource
	descMode     int
	descLastTime time.Time

	protoset    multiString
	protoFiles  multiString
	importPaths multiString
}

// NewDescSourceEntry TODO
func NewDescSourceEntry() *DescSourceEntry {
	ctx, cancel := context.WithCancel(context.Background())

	desc := &DescSourceEntry{
		ctx:    ctx,
		cancel: cancel,
	}

	return desc
}

// SetProtoSetFiles TODO
func SetProtoSetFiles(fileName string) (string, error) {
	data, err := fileReadByte(fileName)
	if err != nil {
		return "", err
	}

	descSourceController.SetProtoSetFiles(fileName)
	return makeHexMD5(data), nil
}

// SetProtoFiles TODO
func SetProtoFiles(importPath string, protoFile string) {
	descSourceController.SetProtoFiles(importPath, protoFile)
}

// SetMode TODO
func SetMode(mode int) {
	descSourceController.SetMode(mode)
}

// GetDescSource TODO
func GetDescSource() (DescriptorSource, error) {
	return descSourceController.GetDescSource()
}

// GetRemoteDescSource TODO
func GetRemoteDescSource(target string) error {
	// parse proto by grpc reflet api
	return nil
}

// InitDescSource TODO
func InitDescSource() error {
	return descSourceController.InitDescSource()
}

// AysncNotifyDesc TODO
func AysncNotifyDesc() {
	descSourceController.AysncNotifyDesc()
}

// SetProtoSetFiles TODO
func (d *DescSourceEntry) SetProtoSetFiles(fileName string) {
	d.protoset.Set(fileName)
}

// SetProtoFiles TODO
func (d *DescSourceEntry) SetProtoFiles(importPath string, protoFile string) {
	d.importPaths.Set(importPath)
	d.protoFiles.Set(protoFile)
}

// SetMode TODO
func (d *DescSourceEntry) SetMode(mode int) {
	d.descMode = mode
}

// GetDescSource TODO
func (d *DescSourceEntry) GetDescSource() (DescriptorSource, error) {
	switch d.descMode {
	case ProtoSetMode:
		return d.descSource, nil

	case ProtoFilesMode:
		return d.descSource, nil

	default:
		return d.descSource, errors.New("only eq ProtoSetMode and ProtoFilesMode")
	}
}

// InitDescSource TODO
func (d *DescSourceEntry) InitDescSource() error {
	var err error
	var desc DescriptorSource

	switch d.descMode {
	case ProtoSetMode:
		// parse proto by protoset

		if d.protoset.IsEmpty() {
			return errors.New("protoset null")
		}

		for _, f := range d.protoset {
			ok := pathExists(f)
			if !ok {
				return errors.New("protoset file not exist")
			}
		}

		desc, err = DescriptorSourceFromProtoSets(d.protoset...)
		if err != nil {
			return errors.New("Failed to process proto descriptor sets")
		}

		d.descSource = desc

	case ProtoFilesMode:
		// parse proto by protoFiles
		d.descSource, err = DescriptorSourceFromProtoFiles(
			d.importPaths,
			d.protoFiles...,
		)
		if err != nil {
			return errors.New("Failed to process proto source files")
		}

	default:
		return errors.New("only eq ProtoSetMode and ProtoFilesMode")
	}

	return nil
}

// AysncNotifyDesc TODO
func (d *DescSourceEntry) AysncNotifyDesc() {
	go func() {
		q := make(chan os.Signal, 1)
		signal.Notify(q, syscall.SIGINT)

		for {
			select {
			case <-q:
				d.InitDescSource()

			case <-d.ctx.Done():
				return
			}
		}
	}()
}

// Close TODO
func (d *DescSourceEntry) Close() {
	d.cancel()
}

// EngineHandler TODO
type EngineHandler struct {
	// grpc clients
	target      string
	clients     map[string]*grpc.ClientConn
	clientsLock sync.RWMutex

	eventHandler InvocationEventHandler
	descCtl      *DescSourceEntry
	invokeCtl    *InvokeHandler

	dialTime      time.Duration
	keepAliveTime time.Duration
	typeCacher    *protoTypesCache
	recvMsgSize   int

	ctx    context.Context
	cancel context.CancelFunc
}

// Option TODO
type Option func(*EngineHandler) error

// SetDialTime TODO
func SetDialTime(val time.Duration) Option {
	return func(o *EngineHandler) error {
		o.dialTime = val
		return nil
	}
}

// SetKeepAliveTime TODO
func SetKeepAliveTime(val time.Duration) Option {
	return func(o *EngineHandler) error {
		o.keepAliveTime = val
		return nil
	}
}

// SetRecvMsgSize TODO
func SetRecvMsgSize(val int) Option {
	return func(o *EngineHandler) error {
		o.recvMsgSize = val
		return nil
	}
}

// SetCtx TODO
func SetCtx(val context.Context, cancel context.CancelFunc) Option {
	return func(o *EngineHandler) error {
		o.ctx = val
		o.cancel = cancel
		return nil
	}
}

// SetDescSourceCtl TODO
func SetDescSourceCtl(val *DescSourceEntry) Option {
	return func(o *EngineHandler) error {
		o.descCtl = val
		return nil
	}
}

// SetHookHandler TODO
func SetHookHandler(handler InvocationEventHandler) Option {
	return func(o *EngineHandler) error {
		o.eventHandler = handler
		return nil
	}
}

// New TODO
func New(options ...Option) (*EngineHandler, error) {
	e := new(EngineHandler)

	// default values
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.dialTime = 10 * time.Second
	e.keepAliveTime = 64 * time.Second
	e.eventHandler = defaultInEventHooker
	e.descCtl = NewDescSourceEntry()
	e.clients = make(map[string]*grpc.ClientConn, 10)
	e.typeCacher = newProtoTypeCache()
	e.recvMsgSize = 100 * 1024 * 1024

	for _, opt := range options {
		if opt != nil {
			if err := opt(e); err != nil {
				return nil, err
			}
		}
	}

	return e, nil
}

// DoConnect TODO
func (e *EngineHandler) DoConnect(target string) (*grpc.ClientConn, error) {
	e.clientsLock.RLock() // read lock
	if conn, ok := e.clients[target]; ok {
		e.clientsLock.RUnlock()
		return conn, nil
	}

	e.clientsLock.RUnlock()

	ctx, _ := context.WithTimeout(e.ctx, e.dialTime)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(e.recvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    e.keepAliveTime,
			Timeout: e.keepAliveTime,
		}))

	cc, err := BlockingDial(ctx, target, opts...)
	if err != nil {
		return cc, err
	}

	e.clientsLock.Lock() // write lock
	defer e.clientsLock.Unlock()

	e.clients[target] = cc
	return cc, err
}

// Init TODO
func (e *EngineHandler) Init(target string) error {
	e.target = target
	if e.descCtl.descMode == ProtoReflectMode {
		var (
			cc        *grpc.ClientConn
			connErr   error
			refClient *grpcreflect.Client

			addlHeaders multiString
			reflHeaders multiString

			descSource DescriptorSource
		)
		md := MetadataFromHeaders(append(addlHeaders, reflHeaders...))
		refCtx := metadata.NewOutgoingContext(e.ctx, md)
		cc, connErr = e.DoConnect(target)
		if connErr != nil {
			return connErr
		}
		refClient = grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(cc))
		descSource = DescriptorSourceFromServer(e.ctx, refClient)
		e.descCtl.descSource = descSource
	} else {
		e.descCtl.InitDescSource()
	}

	err := e.InitFormater()
	if err != nil {
		return err
	}

	return nil
}

// SetMode TODO
func (e *EngineHandler) SetMode(mode int) {
	e.descCtl.SetMode(mode)
}

// SetProtoSetFiles TODO
func (e *EngineHandler) SetProtoSetFiles(fileName string) (string, error) {
	data, err := fileReadByte(fileName)
	if err != nil {
		return "", err
	}

	e.descCtl.SetProtoSetFiles(fileName)
	return makeHexMD5(data), nil
}

// InitFormater TODO
func (e *EngineHandler) InitFormater() error {
	formater, err := ParseFormatterByDesc(e.descCtl.descSource, true)
	if err != nil {
		return err
	}

	inEventHandler := SetDefaultEventHandler(e.descCtl.descSource, formater)
	e.invokeCtl = newInvokeHandler(inEventHandler, e.eventHandler)
	return nil
}

// Close TODO
func (e *EngineHandler) Close() error {
	e.cancel()

	e.clientsLock.Lock()
	defer e.clientsLock.Unlock()
	for _, cc := range e.clients {
		cc.Close()
	}

	return nil
}

// CallWithCtx TODO
func (e *EngineHandler) CallWithCtx(ctx context.Context, serviceName, methodName, data string) (*ResultModel, error) {
	return e.invokeCall(ctx, nil, e.target, serviceName, methodName, data, nil)
}

// CallWithCtxAndHeaders TODO
func (e *EngineHandler) CallWithCtxAndHeaders(ctx context.Context, serviceName, methodName, data string,
	rpcHeaders multiString) (*ResultModel, error) {
	return e.invokeCall(ctx, nil, e.target, serviceName, methodName, data, rpcHeaders)
}

// Call TODO
func (e *EngineHandler) Call(serviceName, methodName, data string) (*ResultModel, error) {
	return e.invokeCall(e.ctx, nil, e.target, serviceName, methodName, data, nil)
}

// CallWithAddr TODO
func (e *EngineHandler) CallWithAddr(serviceName, methodName, data string) (*ResultModel, error) {
	return e.invokeCall(e.ctx, nil, e.target, serviceName, methodName, data, nil)
}

// CallWithAddrCtx TODO
func (e *EngineHandler) CallWithAddrCtx(ctx context.Context, serviceName, methodName, data string) (*ResultModel,
	error) {
	return e.invokeCall(ctx, nil, e.target, serviceName, methodName, data, nil)
}

// CallWithClient TODO
func (e *EngineHandler) CallWithClient(client *grpc.ClientConn, serviceName, methodName, data string) (*ResultModel,
	error) {
	return e.CallWithClientCtx(nil, client, serviceName, methodName, data)
}

// CallWithClientCtx TODO
func (e *EngineHandler) CallWithClientCtx(ctx context.Context, client *grpc.ClientConn, serviceName, methodName,
	data string) (*ResultModel, error) {
	if client == nil {
		return nil, errors.New("invalid grpc client")
	}

	return e.invokeCall(ctx, client, "", serviceName, methodName, data, nil)
}

// invokeCall request target grpc server
func (e *EngineHandler) invokeCall(ctx context.Context, gclient *grpc.ClientConn, target, serviceName, methodName,
	data string, rpcHeaders multiString) (*ResultModel, error) {
	if serviceName == "" || methodName == "" {
		return nil, errors.New("serverName or methodName is null")
	}

	if gclient == nil && target == "" {
		return nil, errors.New("target addr is null")
	}

	if ctx == nil {
		ctx = e.ctx
	}

	var (
		err     error
		cc      *grpc.ClientConn
		connErr error
		// refClient *grpcreflect.Client

		addlHeaders multiString

		// reflHeaders multiString

		descSource DescriptorSource
	)

	descSource = e.descCtl.descSource

	// parse proto by grpc reflet api
	// if e.descCtl.descMode == ProtoReflectMode {
	// 	md := MetadataFromHeaders(append(addlHeaders, reflHeaders...))
	// 	refCtx := metadata.NewOutgoingContext(e.ctx, md)
	// 	cc, connErr = e.DoConnect(target)
	// 	if connErr != nil {
	// 		return nil, connErr
	// 	}
	// 	refClient = grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(cc))
	// 	descSource = DescriptorSourceFromServer(e.ctx, refClient)
	// }

	if gclient == nil {
		cc, connErr = e.DoConnect(target)
		if connErr != nil {
			return nil, connErr
		}
	} else {
		cc = gclient
	}

	var inData io.Reader
	inData = strings.NewReader(data)
	rf, err := RequestParserFor(descSource, inData)
	if err != nil {
		return nil, errors.New("request parse and format failed")
	}

	result, err := e.invokeCtl.InvokeRPC(ctx, descSource, cc, serviceName, methodName,
		append(addlHeaders, rpcHeaders...),
		rf.Next,
	)
	return result, err
}

// ListServices TODO
func (e *EngineHandler) ListServices() ([]string, error) {
	return e.descCtl.descSource.ListServices()
}

// ListMethods TODO
func (e *EngineHandler) ListMethods(svc string) ([]string, error) {
	return ListMethods(e.descCtl.descSource, svc)
}

// ServMethodModel TODO
type ServMethodModel struct {
	PackageName     string
	ServiceName     string
	FullServiceName string
	MethodName      string
	FullMethodName  string
}

// ListServiceAndMethods TODO
func (e *EngineHandler) ListServiceAndMethods() (map[string][]ServMethodModel, error) {
	servList, err := e.ListServices()
	if err != nil {
		return nil, err
	}

	m := map[string][]ServMethodModel{}
	for _, svc := range servList {
		fullMethodList, err := e.ListMethods(svc)
		servMethodModelList := []ServMethodModel{}
		for _, method := range fullMethodList {
			cs := strings.Split(method, ".")
			if len(cs) < 3 {
				return nil, errors.New("method split failed")
			}

			dto := ServMethodModel{
				MethodName:      cs[len(cs)-1],
				ServiceName:     cs[len(cs)-2],
				PackageName:     strings.Join(cs[:len(cs)-2], "."),
				FullMethodName:  method,
				FullServiceName: svc,
			}
			servMethodModelList = append(servMethodModelList, dto)
		}

		if err != nil {
			return nil, err
		}

		m[svc] = servMethodModelList
	}

	return m, nil
}

// ExtractProtoType TODO
func (e *EngineHandler) ExtractProtoType(svc, mth string) (proto.Message, proto.Message, error) {
	var (
		descSource DescriptorSource
	)

	// get types from cache
	key := e.typeCacher.makeKey(svc, mth)
	model, ok := e.typeCacher.get(key)
	if ok {
		return model.reqType, model.respType, nil
	}

	descSource = e.descCtl.descSource
	dsc, err := descSource.FindSymbol(svc)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil, errors.New("not find service in pb descriptor")
		}

		return nil, nil, errors.New("query service failed in pb descriptor")
	}

	sd, ok := dsc.(*desc.ServiceDescriptor)
	if !ok {
		return nil, nil, errors.New("not expose service")
	}

	mtd := sd.FindMethodByName(mth)
	if mtd == nil {
		return nil, nil, fmt.Errorf("service %q does not include a method named %q", svc, mth)
	}

	var ext dynamic.ExtensionRegistry
	alreadyFetched := map[string]bool{}
	if err = fetchAllExtensions(descSource, &ext, mtd.GetInputType(), alreadyFetched); err != nil {
		return nil, nil, fmt.Errorf("error resolving server extensions for message %s: %v", mtd.GetInputType().
			GetFullyQualifiedName(), err)
	}

	if err = fetchAllExtensions(descSource, &ext, mtd.GetOutputType(), alreadyFetched); err != nil {
		return nil, nil, fmt.Errorf("error resolving server extensions for message %s: %v", mtd.GetOutputType().
			GetFullyQualifiedName(), err)
	}

	msgFactory := dynamic.NewMessageFactoryWithExtensionRegistry(&ext)
	req := msgFactory.NewMessage(mtd.GetInputType())
	reply := msgFactory.NewMessage(mtd.GetOutputType())

	// set types to cache
	e.typeCacher.set(key, req, reply)
	return req, reply, nil
}

type multiString []string

// String TODO
func (s *multiString) String() string {
	return strings.Join(*s, ",")
}

// IsEmpty TODO
func (s *multiString) IsEmpty() bool {
	if len(*s) > 0 {
		return false
	}

	return true
}

// Set TODO
func (s *multiString) Set(value string) error {
	*s = append(*s, value)
	return nil
}
