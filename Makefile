ECR_REPO=514200551837.dkr.ecr.us-west-1.amazonaws.com/go-builders
#APP_ECR_REPO=514200551837.dkr.ecr.us-west-1.amazonaws.com/interviews/intuit
APP_ECR_REPO=prophesional/builds
APP_TAG=1.1
BUILDER_VERSION=0.2
path=/c/Users/proph/go-playground/src/github.com/prophesional/intuit/

docker.test.intuit:
	docker run  -v $(path):/go/src/github.com/prophesional/intuit \
    -e DB_SERVER_NAME=$(DB_SERVER_NAME) \
    -e DB_DATABASE_NAME=$(DB_DATABASE_NAME) \
    -e DB_DATABASE_USERNAME=$(DB_DATABASE_USERNAME) \
    -e DB_DATABASE_PASSWORD=$(DB_DATABASE_PASSWORD) \
    -e AWS_REGION=us-west-1 \
    -e DB_SERVER_TYPE=mysql \
	-w /go/src/github.com/prophesional/intuit \
	--entrypoint /bin/bash \
	$(ECR_REPO):$(BUILDER_VERSION) \
	-c "GOOS=linux GO111MODULE=on PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go test -v "

docker.build.intuit:
# 	aws ecr get-login --region us-west-1 --no-include-email | sh
	docker run  -v $(path):/go/src/github.com/prophesional/intuit \
	-w /go/src/github.com/prophesional/intuit \
    --entrypoint /bin/bash \
   	$(ECR_REPO):$(BUILDER_VERSION) \
	-c "GOOS=linux GO111MODULE=on  PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go build -o server ./cmd/server"

docker.build.image.intuit:
	docker build . -t $(APP_ECR_REPO):$(APP_TAG) \
	--build-arg DB_SERVER_NAME=$(DB_SERVER_NAME) \
    --build-arg DB_DATABASE_NAME=$(DB_DATABASE_NAME) \
    --build-arg DB_DATABASE_USERNAME=$(DB_DATABASE_USERNAME) \
    --build-arg DB_DATABASE_PASSWORD=$(DB_DATABASE_PASSWORD) \
    --build-arg AWS_REGION=us-west-1 \
    --build-arg DB_DATABASE_TYPE=mysql \
	--build-arg AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
    --build-arg AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
    --build-arg AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION)

push.image.intuit:
	docker push $(APP_ECR_REPO):$(APP_TAG)
