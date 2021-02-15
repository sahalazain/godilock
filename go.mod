module github.com/sahalazain/godilock

go 1.15

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc v1.34.0 => google.golang.org/grpc v1.29.0

replace go.etcd.io/etcd/api/v3 v3.5.0-pre => go.etcd.io/etcd/api/v3 v3.0.0-20210210213918-d33a1c91f0cf

replace go.etcd.io/etcd/pkg/v3 v3.5.0-pre => go.etcd.io/etcd/pkg/v3 v3.0.0-20210210213918-d33a1c91f0cf

require (
	github.com/coocood/freecache v1.1.1
	github.com/go-redis/redis/v8 v8.5.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	github.com/stretchr/testify v1.6.1
	go.etcd.io/etcd/client/v3 v3.0.0-20210210213918-d33a1c91f0cf
)
