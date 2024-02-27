FROM golang:latest

LABEL maintainer="StoychoMihaylov <st.mihaylov.mihaylov@gmail.com>"

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
ENV PORT 8000

RUN go build
CMD ["./bitshare-chain"]
