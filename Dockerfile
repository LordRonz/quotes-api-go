FROM golang:alpine AS build

ARG PORT=8080
ARG DB_HOST
ARG DB_USER
ARG DB_PASS
ARG DB_NAME
ARG DB_PORT
ARG REDIS_URL
ARG REDIS_PASS

ENV PORT=${PORT}
ENV DB_HOST=${DB_HOST}
ENV DB_USER=${DB_USER}
ENV DB_PASS=${DB_PASS}
ENV DB_NAME=${DB_NAME}
ENV DB_PORT=${DB_PORT}
ENV REDIS_URL=${REDIS_URL}
ENV REDIS_PASS=${REDIS_PASS}

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /backend


FROM alpine:latest AS run

# Copy the application executable from the build image
COPY --from=build /backend /backend

WORKDIR /app
EXPOSE 8080
CMD [ "/backend" ]