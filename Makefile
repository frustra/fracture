protobuf-deps:
	# OS X: brew install protobuf
	# ArchLinux: pacman -S protobuf
	go get -u code.google.com/p/gogoprotobuf/{proto,protoc-gen-gogo,gogoproto}

protobufs:
	protoc --gogo_out=. protobuf/*.proto

