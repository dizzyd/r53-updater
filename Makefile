
APP=r53-updater

GOOP = $(shell which goop)
ifndef $(GOOP)
GOOP = .vendor/bin/goop
endif

TOP := $(shell pwd)

all: .vendor
	GOPATH=$(TOP) $(GOOP) go build $(APP)

ifeq ($(GOOP), .vendor/bin/goop)
.vendor: .vendor/bin/goop
else
.vendor:
endif
	GOPATH=$(TOP)/.vendor $(GOOP) install

clean:
	@rm -rf .vendor $(APP)

.vendor/bin/goop:
	@mkdir .vendor
	(cd .vendor && GOPATH=$(TOP)/.vendor go get github.com/nitrous-io/goop)


dbuild:
	docker run --rm -v $(PWD):/usr/src/$(APP) -w /usr/src/$(APP) golang:1.3 make

