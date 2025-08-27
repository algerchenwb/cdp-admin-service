FROM hub-docker.yuntiancloud.com/esport/base/vr-golang-runtime:1.15-alpine
COPY cdp-admin-service /go
RUN mkdir /go/logs /go/etc /go/public_conf \
    && cd /go \
EXPOSE 38116
CMD ["/go/cdp-admin-service"]