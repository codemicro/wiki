FROM golang:1 as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

ADD https://github.com/sass/dart-sass/releases/download/1.52.1/dart-sass-1.52.1-linux-x64.tar.gz /opt/
RUN ["tar", "-C", "/opt/", "-xzvf", "/opt/dart-sass-1.52.1-linux-x64.tar.gz"]
RUN ["/opt/dart-sass/sass", "wiki/static:wiki/static"]
RUN rm wiki/static/*.scss

RUN CGO_ENABLED=1 GOOS=linux go build -a -buildvcs=false -installsuffix cgo -ldflags "-extldflags '-static'" -o main github.com/codemicro/wiki/wiki

FROM alpine
COPY --from=builder /build/main /
WORKDIR /run

CMD ["../main"]
