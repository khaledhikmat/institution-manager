BUILD_DIR="./build"
DIST_DIR="./dist"

clean_build:
	if [ -d "${BUILD_DIR}" ]; then rm -r ${BUILD_DIR}; fi

clean_dist:
	if [ -d "${DIST_DIR}" ]; then rm -r ${DIST_DIR}; fi; mkdir ${DIST_DIR}

test:
	echo "Invoking test cases..."

build: clean_dist clean_build test
	GOOS='linux' GOARCH='amd64' GO111MODULE='on' go build -o "${BUILD_DIR}/institution-manager-campaign" ./campaign/main.go
	GOOS='linux' GOARCH='amd64' GO111MODULE='on' go build -o "${BUILD_DIR}/institution-manager-externalizer" ./externalizer/main.go
	GOOS='linux' GOARCH='amd64' GO111MODULE='on' go build -o "${BUILD_DIR}/institution-manager-htmx" ./htmx/main.go
	GOOS='linux' GOARCH='amd64' GO111MODULE='on' go build -o "${BUILD_DIR}/institution-manager-institution" ./institution/main.go
	GOOS='linux' GOARCH='amd64' GO111MODULE='on' go build -o "${BUILD_DIR}/institution-manager-member" ./member/main.go
	GOOS='linux' GOARCH='amd64' GO111MODULE='on' go build -o "${BUILD_DIR}/institution-manager-membership" ./membership/main.go
	GOOS='linux' GOARCH='amd64' GO111MODULE='on' go build -o "${BUILD_DIR}/institution-manager-notifier" ./notifier/main.go

dockerize: clean_dist clean_build test build
	docker buildx build --platform linux/amd64 -t khaledhikmat/institution-manager-campaign:latest ./campaign -f ./campaign/Dockerfile
	docker buildx build --platform linux/amd64 -t khaledhikmat/institution-manager-externalizer:latest ..externalizer -f ./externalizer/Dockerfile
	docker buildx build --platform linux/amd64 -t khaledhikmat/institution-manager-htmx:latest ./htmx -f ./htmx/Dockerfile
	docker buildx build --platform linux/amd64 -t khaledhikmat/institution-manager-institution:latest ./institution -f ./institution/Dockerfile
	docker buildx build --platform linux/amd64 -t khaledhikmat/institution-manager-member:latest ./member -f ./member/Dockerfile
	docker buildx build --platform linux/amd64 -t khaledhikmat/institution-manager-membership:latest ./membership -f ./membership/Dockerfile
	docker buildx build --platform linux/amd64 -t khaledhikmat/institution-manager-notifier:latest ./notifier -f ./notifier/Dockerfile

start: clean_dist clean_build test build
	dapr run -f .

list: 
	dapr list

stop: 
	./stop.sh
	# For some reason, the awk command is not working in makefile
	# Ignore errors by placing `-` at the beginning of the line	
	# -dapr stop -f .
	# -(lsof -i:8080 | grep main) | awk '{print $2}' | xargs kill
	# -(lsof -i:8081 | grep main) | awk '{print $2}' | xargs kill
	# -(lsof -i:8082 | grep main) | awk '{print $2}' | xargs kill
	# -(lsof -i:8083 | grep main) | awk '{print $2}' | xargs kill
	# -(lsof -i:8084 | grep main) | awk '{print $2}' | xargs kill
	# -(lsof -i:8085 | grep main) | awk '{print $2}' | xargs kill
	# -(lsof -i:3000 | grep main) | awk '{print $2}' | xargs kill