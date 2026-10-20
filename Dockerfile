ARG IMAGE_REGISTRY=docker.io/library
FROM ${IMAGE_REGISTRY}/golang:1.26-alpine AS build
ARG GOPROXY=https://proxy.golang.org,direct
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/echo-learning .

FROM ${IMAGE_REGISTRY}/alpine:3.23
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /out/echo-learning /app/echo-learning
USER app
ENV PORT=8080 TZ=Asia/Shanghai
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
    CMD wget -qO- "http://127.0.0.1:${PORT}/api/courses" | grep -q '"code":0' || exit 1
ENTRYPOINT ["/app/echo-learning"]
