dep:
	dep ensure -vendor-only
init:
	cp configs/config.toml.dist configs/config.toml
build:
	 go build -o ./assisted_team_api cmd/assisted_team/main.go
run:
	./assisted_team_api