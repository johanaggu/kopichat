FROM golang:1.24 AS build

WORKDIR /src

COPY ../../../go.mod  ../../../go.sum ./
RUN go mod download

COPY . .
RUN go build -tags lambda.norpc -o main ./cmd/lambda/migrations/main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /src/main main

ENTRYPOINT [ "./main" ]