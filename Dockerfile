FROM golang:1.18-alpine AS gobuild

ARG REVISION

WORKDIR /api

COPY go.mod .
COPY go.sum .
COPY vendor/ vendor
COPY cmd/ cmd
COPY pkg/ pkg
COPY migrations/ migrations
RUN mkdir -p images
COPY thirdparty/wait-for .
RUN chmod +x wait-for

RUN go build -o app -mod=vendor -ldflags \
    "-X 'main.buildTimestamp=$(date '+%b %d %Y %T')' -X main.revision=$REVISION" \
    cmd/app/*.go


FROM alpine:3.15

RUN apk --no-cache add ca-certificates

COPY --from=gobuild /api/app .
COPY --from=gobuild /api/migrations .
COPY --from=gobuild /api/wait-for .

EXPOSE 8000

CMD ["./app"]
