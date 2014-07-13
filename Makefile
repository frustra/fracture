all: goprotobuf

goprotobuf:
	cd protobuf; protoc --proto_path=../../../..:../../../../code.google.com/p/gogoprotobuf/protobuf/:. --gogo_out=. *.proto

deps:
	# OS X: brew install protobuf
	# ArchLinux: pacman -S protobuf
	go get -u code.google.com/p/gogoprotobuf/gogoproto
	go get -u code.google.com/p/gogoprotobuf/protoc-gen-gogo
	go get -u ./...
