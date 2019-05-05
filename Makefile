ECR_REPO=514200551837.dkr.ecr.us-west-1.amazonaws.com/go-builders
BUILDER_VERSION=0.1

docker.test.intuit:
	docker run  -v `pwd`:/go/src/github.com/prophesional/intuit \
    -e DB_SERVER_NAME=$(DB_SERVER_NAME) \
    -e DB_DATABASE_NAME=$(DB_DATABASE_NAME) \
    -e DB_DATABASE_USERNAME=$(DB_DATABASE_USERNAME) \
    -e DB_DATABASE_PASSWORD=$(DB_DATABASE_PASSWORD) \
    -e AWS_REGION=us-west-1 \
    -e DB_SERVER_TYPE=mysql \
	-w /go/src/github.com/prophesional/intuit \
	--entrypoint /bin/bash \
	$(ECR_REPO):$(BUILDER_VERSION) \
	-c " GO111MODULE=on PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go test -v "

docker.build.playlist:
# 	aws ecr get-login --registry-ids $(STREAMING_ECR_ID) --region $(STREAMING_ECR_REGION) --no-include-email | sh
	docker run -v $$HOME/.ssh:/root/.ssh -v `pwd`:/go/src/github.com/tunein/streaming-playlist/ \
	-w /go/src/github.com/tunein/streaming-playlist/ \
	--entrypoint /bin/bash \
	$(STREAMING_ECR_ID).dkr.ecr.$(STREAMING_ECR_REGION).amazonaws.com/builder-base:streaming-build-container-$(STREAMING_BUILD_CONTAINER_VERSION) \
	-c "glide install && PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go build -o playlist ./cmd/server"

docker.build.image.playlist:
	docker build . -t $(STREAMING_ECR_ID).dkr.ecr.$(STREAMING_ECR_REGION).amazonaws.com/$(APP):$(IMAGE_TAG) \
	--build-arg AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
    --build-arg AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
    --build-arg AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION)

docker.build.image.playlist.local:
	docker build . -t $(STREAMING_ECR_ID).dkr.ecr.$(STREAMING_ECR_REGION).amazonaws.com/$(APP):$(IMAGE_TAG) \
	--build-arg AWS_ACCESS_KEY_ID=`aws configure get aws_access_key_id` \
    --build-arg AWS_SECRET_ACCESS_KEY=`aws configure get aws_secret_access_key` \
    --build-arg AWS_DEFAULT_REGION=us-west-2 \
    --build-arg SHOCKWAVE_HOST=dev-shockwave.fns.tunein.com:9090

push.playlist:
	docker push $(STREAMING_ECR_ID).dkr.ecr.$(STREAMING_ECR_REGION).amazonaws.com/$(APP):$(IMAGE_TAG)

run.minikube.playlist.local:
	eval $$(minikube docker-env) && \
	kubectl run --rm -i playlist --image=$(STREAMING_ECR_ID).dkr.ecr.$(STREAMING_ECR_REGION).amazonaws.com/$(APP):$(IMAGE_TAG) --image-pull-policy=Never

push.playlist.latest:
	docker run -v $$HOME/.ssh:/root/.ssh -v `pwd`:/go/src/github.com/tunein/streaming-playlist/ \
	-w /go/src/github.com/tunein/streaming-playlist/ \
	--entrypoint /bin/bash tunein/streaming-build-container:0.11 \
	-c "glide install && PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go build -o playlist ./cmd/server" \
	&& docker build . -t $(STREAMING_ECR_ID).dkr.ecr.$(STREAMING_ECR_REGION).amazonaws.com/$(APP):latest \
	--build-arg AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
	--build-arg AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
	--build-arg AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) \
	--build-arg SQL_SERVER_USER=$(SQL_SERVER_USER) \
	&& docker push $(STREAMING_ECR_ID).dkr.ecr.$(STREAMING_ECR_REGION).amazonaws.com/$(APP):latest

build.image.playlist.local:
	docker build . -t $(LOCAL_DOCKER_REGISTRY)/$(APP):latest

push.image.playlist.local:
	docker push $(LOCAL_DOCKER_REGISTRY)/$(APP):latest

build.image.push.local:
	@make docker.build.playlist
	@make build.image.playlist.local
	@make push.image.playlist.local
