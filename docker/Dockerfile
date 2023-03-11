FROM alpine:3.17 as build

RUN  apk --update add go make musl-dev linux-headers git

ARG  BUILD_DIR=/src/cciu

COPY . $BUILD_DIR
RUN  env CGO_ENABLED=0 make -C $BUILD_DIR cciu

###

FROM alpine:3.17
ARG  BUILD_DIR=/src/cciu
RUN  apk --update add ca-certificates
COPY --from=build $BUILD_DIR/bin/cciu /usr/bin/cciu
