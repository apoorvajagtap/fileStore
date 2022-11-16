<< Work In Progress.>>
To check the add and ls functionalities for now, following steps could help.

$ git clone <repository>
$ go build -o ./store client/client.go

// start the server
$ go run server/server.go
$ ./store add <file_path_to_upload>
$ ./store ls

For now, the files will be uploaded to local filesystem. The server shall create a folder named as `./uploads` in PWD.