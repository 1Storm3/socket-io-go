public:
	git checkout v4
	git pull
	git tag -a v4.1.0
	git push origin v4.1.0
	env GOPROXY=proxy.golang.org
	go list -m github.com/1Storm3/socket-io-go/v4