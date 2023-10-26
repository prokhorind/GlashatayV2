.PHONY: dist-aws
## dist-aws: creates the bundle file for AWS Lambda deployments
dist-aws:
	rm -rf dist/aws && mkdir -p dist/aws

	cd cmd/aws ; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../dist/aws/glashatay .
	cd dist/aws ; \
				zip -m lambda.zip glashatay ;