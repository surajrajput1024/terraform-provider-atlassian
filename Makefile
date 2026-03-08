.PHONY: build test install docs generate clean

build:
	go build -o terraform-provider-atlassian .

test:
	go test ./...

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/surajrajput1024/atlassian/0.0.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-atlassian ~/.terraform.d/plugins/registry.terraform.io/surajrajput1024/atlassian/0.0.0/$$(go env GOOS)_$$(go env GOARCH)/

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name atlassian

generate: docs

clean:
	rm -f terraform-provider-atlassian
