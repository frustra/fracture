protobuf-deps:
	# OS X: brew install protobuf
	# ArchLinux: pacman -S protobuf
	go get -u code.google.com/p/gogoprotobuf/{proto,protoc-gen-gogo,gogoproto}

protobufs:
	cd protobuf; protoc --proto_path=../../../..:../../../../code.google.com/p/gogoprotobuf/protobuf/:. --gogo_out=. *.proto

