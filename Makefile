
APP=r53-updater

GOOP = $(shell which goop)
ifndef $(GOOP)
GOOP = bin/goop
endif

TOP := $(shell pwd)

all: .vendor
	GOPATH=$(TOP) $(GOOP) go install $(APP)

ifeq ($(GOOP), bin/goop)
.vendor: bin/goop
else
.vendor:
endif
	GOPATH=$(TOP)/.vendor $(GOOP) install

clean:
	@rm -rf pkg bin .vendor .goop

bin/goop:
	@mkdir -p bin
	@mkdir .goop
	(cd .goop && GOPATH=$(TOP)/.goop go get github.com/nitrous-io/goop)
	mv .goop/bin/goop bin

