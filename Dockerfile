FROM golang:1.25-alpine AS stage

WORKDIR /app

RUN apk add --no-cache git


COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api/main.go

FROM scratch

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=stage /app/server .

EXPOSE 8080

CMD [ "./server" ]