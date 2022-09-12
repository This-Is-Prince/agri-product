gen:
	protoc --go_out=pb --proto_path=proto --go_opt=paths=source_relative \
	--go-grpc_out=pb --proto_path=proto --go-grpc_opt=paths=source_relative \
    proto/*.proto

clean:
	rm pb/*.go

run:
	go run *.go