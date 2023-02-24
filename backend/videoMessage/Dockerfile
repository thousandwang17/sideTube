FROM golang:1.18.9

COPY ./ /src/
WORKDIR /src/
RUN ls 
RUN CGO_ENABLED=0 GOOS=linux go build -o videoMessage

FROM alpine:3.17
COPY --from=0 /src/videoMessage /
RUN mkdir /var/videos
CMD ./videoMessage