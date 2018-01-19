# build stage
FROM golang:1.9-alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
# RUN go get -u github.com/golang/dep/cmd/dep
ENV D=/go/src/github.com/x/y
ADD main.go $D/
RUN cd $D && go get
# RUN cd $D && dep init && dep ensure
RUN cd $D && go build -o fnserver && cp fnserver /tmp/

# final stage
FROM fnproject/dind
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build-env /tmp/fnserver /app/fnserver
CMD ["./fnserver"]
