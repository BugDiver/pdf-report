cd gauge-proto
PATH=$PATH:$GOPATH/bin protoc --go_out=../gauge_messages spec.proto messages.proto
cd ..
go fmt ./...