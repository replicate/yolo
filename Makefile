GO_FILES := $(shell find . -name '*.go')

all: build

build: $(GO_FILES)
	go build

sdxl-hello: build
	time ./yolo push --base r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41 --dest r8.im/anotherjesse/yolo examples/hello-world/predict.py --ast examples/hello-world/predict.py

sdxl-no-watermarks: build
	time ./yolo push --base r8.im/stability-ai/sdxl@sha256:1bfb924045802467cf8869d96b231a12e6aa994abfe37e337c63a4e49a8c6c41 --dest r8.im/anotherjesse/yolo examples/sdxl-no-watermarks/predict.py
