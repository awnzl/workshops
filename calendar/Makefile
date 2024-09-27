GOBUILD	:= go build
GOTEST	:= go test

GOBASE	:= $(shell pwd)
CMD		:= $(GOBASE)/cmd
GOBIN	:= $(GOBASE)/bin

APPS	:= $(notdir $(wildcard $(CMD)/*))
TESTS	:= $(GOBASE)/internal/...
# DEVENV  := $(GOBASE)/.env

FLAGS	:= -gcflags="all=-N -l"

define BUILD
$(GOBUILD) $(2) -o $(GOBIN)/$(1) $(CMD)/$(1)/*.go

endef

.PHONY: all
all: build

.PHONY: build
build:
	$(foreach app,$(APPS),$(call BUILD,$(app)))
	# cp -n -f $(DEVENV) $(GOBIN)

.PHONY: delve
delve:
	$(foreach app,$(APPS),$(call BUILD,$(app), $(FLAGS)))

.PHONY: test
test:
	$(GOTEST) $(TESTS)

.PHONY: clean
clean:
	rm -rf $(GOBIN)/$(BINARY_NAME)
