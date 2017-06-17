REPO=malice
NAME=fprot
VERSION=$(shell cat VERSION)

all: gotest build size test

build:
	docker build -t $(REPO)/$(NAME):$(VERSION) .

base:
	docker build -f Dockerfile.base -t $(REPO)/$(NAME):base .

dev: test
	docker build -f Dockerfile.dev -t $(REPO)/$(NAME):$(VERSION) .

size:
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(VERSION)| cut -d' ' -f1)%20GB-blue/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(REPO)/$(NAME)

tar:
	docker save $(REPO)/$(NAME):$(VERSION) -o $(NAME).tar

gotest:
	go get
	go test -v

avtest:
	@echo "===> F-PROT Version"
	@docker run --init --rm --entrypoint=bash $(REPO)/$(NAME):$(VERSION) -c "/usr/local/bin/fpscan --version" > av_version.out
	@echo "===> F-PROT EICAR Test"
	@docker run --init --rm --entrypoint=bash $(REPO)/$(NAME):$(VERSION) -c "/usr/local/bin/fpscan -r EICAR" > av_scan.out || true

test:
	docker run --init -d --name elasticsearch blacktop/elasticsearch
	docker run --init --rm $(REPO)/$(NAME):$(VERSION)
	docker run --init --rm --link elasticsearch $(REPO)/$(NAME):$(VERSION) -V EICAR > results.json
	cat results.json | jq .
	docker rm -f elasticsearch

clean:
	docker-clean stop
	docker rmi $(REPO)/$(NAME):$(VERSION)
	docker rmi $(REPO)/$(NAME):base

.PHONY: build dev size tags test gotest clean
