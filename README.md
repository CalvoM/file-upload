# File-Upload

## How to run

### Generate the environment variables
```sh
go generate ./...
```

Fill the **.env** file with the corresponding environment variables.

**NB** *If opting to use the docker containers*
Please fill the **.docker_env** file with postgres container variables
Create docker network for postgres to connect

### Run the application
```sh
go run .
```


