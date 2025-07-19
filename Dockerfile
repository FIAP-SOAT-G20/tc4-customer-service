# Multi-stage build for AWS Lambda
FROM golang:1.24 as build
WORKDIR /app

# Copy dependencies list
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optional lambda.norpc tag
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -ldflags="-w -s" -o main .

# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /app/main ./main
ENTRYPOINT [ "./main" ]