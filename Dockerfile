FROM golang:1.25-alpine AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM gcr.io/distroless/base-debian13:nonroot AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/main /main

USER 12345
CMD [ "/main" ]
