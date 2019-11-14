ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/sansay_exporter  /bin/sansay_exporter

EXPOSE      9116
ENTRYPOINT  [ "/bin/sansay_exporter" ]
