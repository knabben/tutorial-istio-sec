module github.com/knabben/tutorial-istio-sec/1-grpc/grpc/server

go 1.20

require (
	github.com/knabben/tutorial-istio-sec/1-grpc/grpc/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.55.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/knabben/tutorial-istio-sec/1-grpc/grpc/proto => ../proto
