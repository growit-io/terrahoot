FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk add --no-cache upx

WORKDIR /src
COPY . .

# XXX: Does it make sense to call strip(1) here as well?
RUN go build -ldflags "-s -w -extldflags '-static'" -o /bin/app .
RUN upx -q -9 /bin/app

FROM alpine AS runner

RUN apk add --no-cache bash curl git

COPY --from=builder --chown=0:0 /bin/app /usr/local/bin/terrahoot

RUN adduser -u 1000 -h /home/user -D user user

RUN git clone --depth=1 --branch v3.0.0 https://github.com/tfutils/tfenv.git /home/user/tfenv
RUN ln -s /home/user/tfenv/bin/* /usr/local/bin
COPY .terraform-version /home/user/tfenv/version

RUN git clone --depth 1 --branch v1.2.1 https://github.com/tgenv/tgenv.git /home/user/tgenv
RUN ln -s /home/user/tgenv/bin/* /usr/local/bin
# XXX: https://github.com/tgenv/tgenv/pull/44
#COPY .terragrunt-version /home/user/tgenv/version
COPY .terragrunt-version /home/user/.terragrunt-version

RUN chown -R user: /home/user
USER user

ENTRYPOINT ["/usr/local/bin/terrahoot"]
