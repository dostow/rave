VERSION := 1.0.8
NAME := $(shell echo $${PWD\#\#*/})
TARGET := ./docker/$(NAME)
all: clean build image scaletag scalepush
$(TARGET): 
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -ldflags="-X main.VERSION=$(VERSION) -X main.BUILD=$(shell git describe --always --long --dirty)" -o $(TARGET) github.com/dostow/rave/cmd/rave
build: $(TARGET)
		@true
image:
	@docker build -t $(NAME):$(VERSION) ./docker
tag: 
	@docker tag $(NAME):$(VERSION) docker.registry/$(NAME):$(VERSION)
push: 
	@docker push docker.registry/$(NAME):$(VERSION)
ktag: 
	@docker tag $(NAME):$(VERSION) gcr.io/dostow-api/$(NAME):$(VERSION)  
kpush:
	@docker push gcr.io/dostow-api/$(NAME):$(VERSION)
scaletag: 
	@docker tag $(NAME):$(VERSION) rg.fr-par.scw.cloud/dostow/$(NAME):$(VERSION)  
scalepush:
	@docker push rg.fr-par.scw.cloud/dostow/$(NAME):$(VERSION) 
msave:
	@docker save $(NAME):$(VERSION) -o $(NAME)_$(VERSION).tar
mcopy: 
	@scp ./$(NAME)_$(VERSION).tar 67.222.154.8:~/
minstall:
	@ssh 67.222.154.8 -C microk8s.ctr -n k8s.io image import $(NAME)_$(VERSION).tar
mpush:
	msave mcopy minstall
clean:
	@rm -f $(TARGET)