FROM alpine:edge AS builder

RUN apk add --no-cache go
WORKDIR /src
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN go build -o ./sim_board ./main

FROM alpine:3.23.3 AS runner

WORKDIR /app
COPY --from=builder /src/sim_board ./
RUN chmod +x ./sim_board
EXPOSE 6700
CMD ["./sim_board"]