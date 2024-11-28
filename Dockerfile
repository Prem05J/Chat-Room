FROM golang:1.22.2 AS build-stage 

WORKDIR /app 

COPY go.mod go.sum ./

RUN go mod download 

COPY . .   

RUN CGO_ENABLED=0 GOOS=linux go build -o /api main.go 


FROM scratch AS build-release-stage 


WORKDIR / 

COPY --from=build-stage /api /api 


EXPOSE 8080 

ENTRYPOINT ["/api"]



