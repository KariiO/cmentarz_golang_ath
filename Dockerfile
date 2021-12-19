FROM golang:alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . ./

RUN go get -u github.com/KariiO/cmentarz_golang_ath

RUN CGO_ENABLED=0 GOOS=linux go build main.go

EXPOSE 3000

CMD ["./main"]