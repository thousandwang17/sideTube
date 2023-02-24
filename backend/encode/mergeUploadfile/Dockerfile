FROM golang:1.18.9

COPY ./ /src/
WORKDIR /src/
RUN ls 
RUN CGO_ENABLED=0 GOOS=linux go build -o fileMerge

FROM alpine:3.17
RUN mkdir /var/videos
COPY --from=0 /src/fileMerge /
CMD ./fileMerge