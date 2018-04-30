FROM golang:alpine as build

WORKDIR /go/src/github.com/mhutter/podstalk
COPY . .
RUN go build -v -o /tmp/podstalk ./cmd/podstalk/podstalk.go

FROM alpine

ENV PORT=8080
EXPOSE 8080

COPY --from=build /tmp/podstalk /bin/podstalk
CMD ["/bin/podstalk"]
