#选择基础镜像
FROM golang:1.25-alpine
#FROM golang:alpine(自动拉取最新版本)

#设置工作目录
WORKDIR /app

#复制依赖文件
COPY go.mod go.sum ./
#下载依赖
RUN go mod download

#复制源码和配置文件到容器里
COPY . .

#编译Go应用,输出文件名为main
#RUN go build -o main main.go
#改为强制静态编译
Run CGO_ENABLED=0 GOOS=linux go build -o main main.go

#暴露端口，用8080
EXPOSE 8080

#启动
CMD ["./main"]