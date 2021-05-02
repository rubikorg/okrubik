FROM golang

# add all files from service folder to the container
ADD . ../cmd/{{ .Name }}

# okrubik dock {{ .Name }} || okrubik dock --up {{ .Name }} generates
# a target folder inside orchestration folder with the name of the 
# service. Add it to the docker container.

# the $RUBIK_ENV os the environment in which your server needs to run
# in and can be passed using `okrubik dock --env production server`
ADD . target/$RUBIK_ENV/{{ .Name }}

# install modules
RUN go mod tidy

# run the binary copied from the target/ folder
RUN ./{{ .Name }}

EXPOSE {{ .Port }}