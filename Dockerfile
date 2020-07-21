FROM golang:1.14 as build

WORKDIR /src

ADD go.mod /src/
ADD go.sum /src/
RUN go mod download

ADD . /src/
RUN CGO_ENABLED=0 GOOS=linux make all

FROM alpine:3.11.6
COPY --from=build /src/build/codeowners-verifier /usr/local/bin

CMD [ "codeowners-verifier", "help" ]
