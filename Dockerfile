FROM golang:alpine as build

WORKDIR /go/src
COPY . .
RUN go build -o /tmp/podstalk .

FROM alpine

ENV PORT=8080
EXPOSE 8080

COPY --from=build /tmp/podstalk /bin/podstalk
CMD ["/bin/podstalk"]
