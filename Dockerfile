FROM golang:alpine AS build
LABEL maintainer="Manuel Hutter (https://github.com/mhutter)"
WORKDIR /src
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

RUN apk --no-cache add git
COPY . .
RUN go build -v -o /tmp/podstalk ./cmd/podstalk


FROM scratch
ENV PORT=8080
EXPOSE 8080

COPY --from=build /tmp/podstalk /podstalk
CMD ['/podstalk']
