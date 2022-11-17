<< Work In Progress.>>
Pre-requisite: MongoDB container running on port 27017
Following steps could be followed to test the CRUD functionality.

$ docker run -d -p 27017:27017 --name test-mongo mongo:4.0.4
$ docker ps

$ git clone <repository>
$ go build -o ./store client/client.go

// start the server
$ go run server/server.go
$ ./store add <file_path_to_upload>
$ ./store ls

For now, the files will be uploaded to local filesystem. The server shall create a folder named as `./uploads` in PWD.