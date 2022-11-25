FROM golang:1.17.1-buster AS tilt

ENV GO111MODULE="on"
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

COPY bin/backend /usr/local/bin/backend
ENTRYPOINT /usr/local/bin/backend

FROM golang:1.17.1-buster AS common

ARG release="not-set"

ENV RELEASE=$release
ENV SCHEMA_PATH="./schema"
ENV PACKAGE_PATH="./backend"
ENV WORKING_DIR="/usr/src/"
ENV GO111MODULE="on"
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

WORKDIR ${WORKING_DIR}
COPY ${PACKAGE_PATH} ${WORKING_DIR}
COPY ${SCHEMA_PATH} ${WORKING_DIR}/schema


# --- builder image
FROM common AS builder

RUN apt-get -q -y update \
  && apt-get -q -y install bash g++ \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN go get \
  github.com/volatiletech/sqlboiler \
  github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql \
  github.com/99designs/gqlgen

# RUN go mod download
RUN bin/prepare.sh --build-gql \
  && bin/prepare.sh -sp cmd -pf main.go -op bin -an backend -bl

# --- Create minimal deployment image containing the binary
FROM common AS deployable

ARG deployment_url="https://app.perspex.us"

ENV PACKAGE_PATH="."
ENV WORKING_DIR="/usr/src/"
ENV DEPLOYMENT_URL=$deployment_url

RUN apt-get -q update -y \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder ${WORKING_DIR}/bin/backend /usr/local/bin/backend

CMD [ "backend" ]
