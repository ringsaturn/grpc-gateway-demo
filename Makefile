pb:
	protoc -I . \
		--go_out ./proto/helloworld --go_opt paths=source_relative \
		--go-grpc_out ./proto/helloworld --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto/helloworld \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out ./openapiv2 \
		--openapiv2_opt logtostderr=true \
		helloworld.proto
