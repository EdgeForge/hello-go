FROM golang:latest as builder

WORKDIR /app

COPY . ./

# not necessary because no external dependencies, but keeping for reference
# RUN go mod download

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o hello-go


FROM scratch

COPY --from=builder /app/hello-go /

CMD ["/hello-go"]
