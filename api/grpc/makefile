outDirectory = protobuff
protoFile = notify.proto
client = client.go
server = server.go

generate:
	mkdir -p $(outDirectory)
	# Generate Go client and server code from .proto files
	protoc --go_out=$(outDirectory) --go-grpc_out=$(outDirectory) --proto_path=. $(protoFile)

	# Create Go client and server files if they don't already exist
	touch $(outDirectory)/$(client)
	touch $(outDirectory)/$(server)
	
clean:
	rm -rf $(outDirectory)
