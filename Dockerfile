# dev image
FROM golang:1.14-alpine as dev

RUN apk add --update --no-cache bash inotify-tools curl git make

ENV CODE=/usr/src/app

RUN mkdir -p ${CODE}/.gopath

# allows static linking for alpine
ENV CGO_ENABLED=0
ENV GOPATH="${CODE}/.gopath"
ENV PATH="${PATH}:${GOPATH}/bin"

WORKDIR ${CODE}

CMD ["air", "-c", "air.toml"]

# build image
FROM dev as build

COPY . ${CODE}/
RUN make compile

# production image
FROM alpine as production

RUN apk --no-cache add ca-certificates

COPY --from=build /usr/src/app/main /root/main
