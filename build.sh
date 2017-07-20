curl https://curl.haxx.se/ca/cacert.pem > cacert.pem
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
chmod +x main
docker build -t snorremd/nicolas-cage-bot:latest .