module github.com/forwardalex/Ytool

go 1.16

require (
	bou.ke/monkey v1.0.2
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/buger/jsonparser v1.1.1
	github.com/coreos/bbolt v1.3.6 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/emicklei/proto v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/gin-gonic/gin v1.7.4
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/jhump/protoreflect v1.10.1
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tal-tech/go-zero v1.2.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace (
	github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6
	google.golang.org/grpc v1.43.0 => google.golang.org/grpc v1.26.0
)
