FROM golang:1.15.2 as builder

# Copy all the source files for the vote service.
RUN mkdir /svc
ADD . /svc

# We specify that we now wish to execute any further commands inside the /svc directory.
WORKDIR /svc

# Build the binary
ENV GOPROXY=direct
RUN go build -o vote ./cmd/vote


# For the real image, we'll only copy the binaryso that the image size is small.
#FROM gcr.io/distroless/base-debian10
#COPY --from=builder /svc/vote /svc
#COPY --from=builder /svc/templates /svc

EXPOSE 8080
CMD ["./vote"]