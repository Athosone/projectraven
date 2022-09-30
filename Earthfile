VERSION 0.6
FROM golang:1.18
WORKDIR /workspace

ARG ENVTEST_K8S_VERSION=1.23
ARG BUILDNUMBER="dev"
ARG DEV="true"

# Set default values if dev
IF [ "$DEV" = "true" ]
  ENV DISTROLESS_TAG=debug
ELSE
  ENV DISTROLESS_TAG=nonroot
END

ENV LDFLAGS="-s -w"

RUN go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo
RUN go install -mod=mod github.com/axw/gocov/gocov
RUN go install -mod=mod github.com/AlekSi/gocov-xml
RUN go install -mod=mod github.com/golang/mock/mockgen
RUN go install -mod=mod sigs.k8s.io/controller-runtime/tools/setup-envtest@latest


DEPS:
  COMMAND
  ARG DEV="false"

  COPY go.mod go.sum ./

  # Output these back in case go mod download changes them.
  SAVE ARTIFACT go.mod AS LOCAL go.mod
  SAVE ARTIFACT go.sum AS LOCAL go.sum

GOBUILD:
  COMMAND
  ARG project_path
  ARG project_name

  COPY . .
  RUN go fmt ./... && go vet ./...
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o bin/$project_name $project_path/
  SAVE ARTIFACT bin/ /app/ AS LOCAL ./bin/

distrolessBase:
  ARG DISTROLESS_TAG
  FROM gcr.io/distroless/static:$DISTROLESS_TAG
  USER 65532:65532
  SAVE IMAGE --cache-hint

GOCOV:
  COMMAND
  ARG coverage_file
  ARG output_file
  RUN gocov convert $coverage_file | gocov-xml > $output_file

