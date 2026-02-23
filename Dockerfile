FROM golang:alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

RUN apk add build-base

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . .

FROM base AS build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags musl -ldflags="-w -s" \
    -o bookmark_service cmd/api/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags musl -ldflags="-w -s" \
    -o bookmark_migration cmd/migrate/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags musl -ldflags="-w -s" \
    -o bookmark_worker cmd/worker/main.go


FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
	grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

FROM scratch AS test
ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

FROM alpine AS final-service

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=build /opt/app/bookmark_service /app/bookmark_service
COPY --from=build /opt/app/bookmark_migration /app/bookmark_migration
COPY --from=build /opt/app/docs /app/docs
COPY --from=build /opt/app/migrations /app/migrations

CMD ["/app/bookmark_service"]

FROM alpine AS final-worker

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=build /opt/app/bookmark_worker /app/bookmark_worker

CMD ["/app/bookmark_worker"]