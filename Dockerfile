FROM golang 1.13

RUN go get -u github.com/Saied74/graphs

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV APP_USER app
ENV APP_HOME /go/src/graphs
ENV GRAPHPATH /go/src/graphs

ARG GROUP_ID
ARG USER_ID

RUN groupadd --gid $GROUP_ID app && useradd -m -l --uid $USER_ID --gid $GROUP_ID $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

USER APP_USER
WORKDIR $APP_HOME

ESPOSE 8080

EBTRYPOINT ["go", "run", "." ]
