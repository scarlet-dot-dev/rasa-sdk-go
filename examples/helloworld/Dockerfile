# Copyright (c) 2020 Eddy <eddy@scarlet.dev>
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

# Builder
FROM golang:1-alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN apk update && apk upgrade && apk add -U --no-cache ca-certificates build-base
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd/rasa-actions

# Runtime
FROM scratch
COPY --from=builder /build/main /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
ENTRYPOINT ["./main"]
