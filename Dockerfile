FROM golang:alpine as build

ENV PORT=8080
EXPOSE 8080

WORKDIR /go/src/
COPY . .
RUN go build -o /tmp/podstalk .

FROM alpine

COPY --from=build /tmp/podstalk /bin/podstalk
CMD ["/bin/podstalk"]
