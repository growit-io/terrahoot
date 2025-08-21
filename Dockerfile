FROM golang:1.25-alpine@sha256:72a77777b4d52d80820a9b5828c06cca1cd3c24721b6ca25b0269b5f6fb00daa AS builder

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk add --no-cache upx

WORKDIR /src
COPY . .

# XXX: Does it make sense to call strip(1) here as well?
RUN go build -ldflags "-s -w -extldflags '-static'" -o /bin/app .
RUN upx -q -9 /bin/app

FROM alpine@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1 AS runner

RUN apk add --no-cache bash curl git

RUN git clone --depth=1 --branch v3.0.0 https://github.com/tfutils/tfenv.git /home/user/tfenv
RUN ln -s /home/user/tfenv/bin/* /usr/local/bin
COPY .terraform-version /home/user/tfenv/version

RUN git clone --depth 1 --branch v1.2.1 https://github.com/tgenv/tgenv.git /home/user/tgenv
RUN ln -s /home/user/tgenv/bin/* /usr/local/bin
# XXX: https://github.com/tgenv/tgenv/pull/44
#COPY .terragrunt-version /home/user/tgenv/version
COPY .terragrunt-version /home/user/.terragrunt-version

COPY --from=builder --chown=0:0 /bin/app /usr/local/bin/terrahoot

ENTRYPOINT ["/usr/local/bin/terrahoot"]
