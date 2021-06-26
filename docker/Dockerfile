FROM alpine:3.13 as build

RUN  apk --update add go make musl-dev linux-headers git

ARG  BUILD_DIR=/src/cciu

COPY . $BUILD_DIR
RUN  make -C $BUILD_DIR cciu

###

FROM alpine:3.13
ARG  BUILD_DIR=/src/cciu
RUN  apk --update add ca-certificates
COPY --from=build $BUILD_DIR/bin/cciu /usr/bin/cciu
