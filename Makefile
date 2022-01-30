all: gobuild

gobuild:
	@mkdir -p ./build
	go build -o build/csg ./cmd