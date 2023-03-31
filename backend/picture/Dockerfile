FROM golang:1.18.9

COPY ./ /src/
WORKDIR /src/
RUN ls 
RUN CGO_ENABLED=0 GOOS=linux go build -o recommend

FROM alpine:3.17
COPY --from=0 /src/recommend /
RUN mkdir /var/videos
CMD ./recommend