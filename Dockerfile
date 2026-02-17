FROM golang:1.25-alpine
WORKDIR /app
COPY ./ ./
RUN go build -o /docker
RUN go mod download

ENV PORT=80
ENV LOG_LEVEL=info
ENV FIREBASE_FILE_PATH=./files/firebase-credentials.json

EXPOSE 80

CMD ["/docker"]