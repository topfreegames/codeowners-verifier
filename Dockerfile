FROM golang:1.19 as build

WORKDIR /src

ADD go.mod /src/
ADD go.sum /src/
RUN go mod download

ADD . /src/
ARG GOOS=linux
ARG GOARCH=amd64
RUN CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH make all

FROM alpine:3.17
COPY --from=build /src/build/codeowners-verifier /usr/local/bin

CMD ["codeowners-verifier", "help" ]
