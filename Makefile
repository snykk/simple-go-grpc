.PHONY: proto run-server run-client

proto:
	protoc \
	--proto_path=proto \
	--go_out=proto/pb \
	--go_opt=paths=source_relative \
    --go-grpc_out=proto/pb \
	--go-grpc_opt=paths=source_relative \
    proto/*.proto

run-client:
	cd client && go run .

run-server:
	cd server && go run .