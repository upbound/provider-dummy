FROM gcr.io/distroless/static@sha256:a01d47d4036cae5a67a9619e3d06fa14a6811a2247b4da72b4233ece4efebd57

ARG TARGETOS
ARG TARGETARCH

WORKDIR /provider
COPY --from=busybox:1.36.0-musl /bin/sh /bin/sh

ADD bin/$TARGETOS\_$TARGETARCH/provider /usr/local/bin/crossplane-dummy-provider
ADD bin/$TARGETOS\_$TARGETARCH/server /usr/local/bin/dummy-server

USER 65532
ENTRYPOINT ["sh", "-c", "dummy-server & crossplane-dummy-provider"]
