OUTPUT=cacheserver
GITVER=`git rev-parse --short HEAD`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-X github.com/n4mine/cacheserver/models.GitVer=${GITVER} -X github.com/n4mine/cacheserver/models.BuildTime=${BUILD_TIME}"

default:
	@echo "gopath $(GOPATH)"
	go build ${LDFLAGS} -o ${OUTPUT}
