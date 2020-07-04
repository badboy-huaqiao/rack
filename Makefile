# Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved. 
# SPDX-License-Identifier: Apache-2.0

GO=CGO_ENABLED=0 GO111MODULE=on GO
CGO=CGO_ENABLED=1 GO111MODULE=on GO

MICROSERVICES=cmd/rack/rack

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)

GOFLAGS=-ldflags "-X rack.Version=$(VERSION)"
GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)
	$(GO) build ./...

cmd/rack/rack: 
	$(GO) build $(GOFLAGS) -o $@ ./cmd/rack

clean:
	rm -f $(MICROSERVICES)

run:
	cd bin && ./rack-launch.sh