FROM golang:1.20.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /stage
COPY . /stage
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /stage/bin/${TARGETOS}_${TARGETARCH}/server .

###
FROM gcr.io/distroless/static@sha256:a01d47d4036cae5a67a9619e3d06fa14a6811a2247b4da72b4233ece4efebd57

ARG TARGETOS
ARG TARGETARCH

COPY --from=builder /stage/bin/${TARGETOS}_${TARGETARCH}/server /usr/local/bin/server-dummy

USER 65532
ENTRYPOINT ["server-dummy"]
