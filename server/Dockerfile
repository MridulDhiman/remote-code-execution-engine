FROM golang:1.22-alpine as builder 

WORKDIR /build

COPY go.mod go.sum ./
# Downloads all the modules and saves them in local cache, whenever there is change in dependencies
RUN go mod download && go mod verify

## copy all the code into container
COPY . .

## build the go code 
RUN go build -o bin .

###

FROM scratch

WORKDIR /

## copies /build/bin from /bin/ directory 
COPY --from=builder /build/bin /bin

## setup amqp url
# ENV AMQP_URL ${AMQP_URL}

EXPOSE 8080

# ENTRYPOINT [ "/bin", "server" ]









