# Build the manager binary
FROM golang:1.13 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARM=$(if [ "$TARGETVARIANT" = "v7" ]; then echo "7"; fi) GOARCH=$TARGETARCH GO111MODULE=on go build -a -o manager main.go

FROM alpine:3.12.1

ARG GIT_COMMIT="unspecified"
LABEL GIT_COMMIT=$GIT_COMMIT

ARG GIT_TAG=""
LABEL GIT_TAG=$GIT_TAG

ARG COMMIT_TIMESTAMP="unspecified"
LABEL COMMIT_TIMESTAMP=$COMMIT_TIMESTAMP

ARG AUTHOR_EMAIL="unspecified"
LABEL AUTHOR_EMAIL=$AUTHOR_EMAIL

ARG SIGNATURE_KEY="undefined"
LABEL SIGNATURE_KEY=$SIGNATURE_KEY

WORKDIR /
RUN wget https://releases.hashicorp.com/terraform/0.13.5/terraform_0.13.5_linux_${$TARGETARCH}.zip -O terraform.zip
RUN unzip terraform.zip
RUN mv terraform /usr/local/bin/
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]