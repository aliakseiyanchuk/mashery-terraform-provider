TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=aliakseiyanchuk
NAME=mashery
BINARY=terraform-provider-${NAME}
VERSION=0.5
BUILD_PRERELEASE=alpha.3
OS_ARCH=darwin_arm64
DOCKER_IMAGE=nexus:5001

LOCAL_MIRROR_PATH=${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

compile_docker_binaries:
	rm -rf ./docker/dist
	GOOS=linux GOARCH=arm64 		go build -o ./docker/dist/${LOCAL_MIRROR_PATH}/linux_arm/${BINARY}			main.go
	GOOS=linux GOARCH=arm GOARM=6 	go build -o ./docker/dist/${LOCAL_MIRROR_PATH}/linux_armv6/${BINARY} 		main.go
	GOOS=linux GOARCH=amd64 		go build -o ./docker/dist/${LOCAL_MIRROR_PATH}/linux_amd64/${BINARY}		main.go
	GOOS=linux GOARCH=386 			go build -o ./docker/dist/${LOCAL_MIRROR_PATH}/linux_386/${BINARY}			main.go

push_dev_container: compile_docker_binaries
	docker build ./docker -t nexus:5001/terraform/mashery-terraform-container:v${VERSION}-${BUILD_PRERELEASE}
	docker push nexus:5001/terraform/mashery-terraform-container:v${VERSION}-${BUILD_PRERELEASE}

push_release_container: compile_docker_binaries
	docker build ./docker -t lspwd2/terraform-provider-mashery:v${VERSION} -t lspwd2/terraform-provider-mashery:latest
	docker image push --all-tags lspwd2/terraform-provider-mashery