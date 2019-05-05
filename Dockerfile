FROM docker.io/library/node:current-slim AS client
WORKDIR /src
COPY ./client .
RUN yarn && yarn build

FROM docker.io/library/golang:alpine AS server
WORKDIR /src
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

RUN apk --no-cache add git
COPY . .
RUN go build -v -o /tmp/podstalk ./cmd/podstalk


FROM scratch
LABEL maintainer="Manuel Hutter (https://github.com/mhutter)"
ENV PORT=8080
EXPOSE 8080

COPY --from=client /src/build /public
COPY --from=server /tmp/podstalk /podstalk
CMD ["/podstalk"]
