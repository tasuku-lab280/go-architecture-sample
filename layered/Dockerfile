FROM golang:1.25-alpine

RUN apk --no-cache add tzdata
ENV TZ=Asia/Tokyo

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

EXPOSE 8080
CMD ["go", "run", "./cmd/api"]
