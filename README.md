
## File Store server
###### Work In Progress

## Current Status
Following tasks have been achieved:
- Add multiple text files to the server. Skips saving the file, if already exists.
- Delete the specified file.
- List all the filenames in the server.
- Update the specified file, or create new if it does not exist.
- Count the total words in all the files stored on the server. 

## Steps to Run Locally.

*Pre-requisite*  
- You have a mongodb container running on port '27017' locally.
```
$ docker run -d -p 27017:27017 --name fileStore-mongo mongo:4.0.4
$ docker ps
```

*Start the server*
- Clone the repository locally
```
$ git clone git@github.com:apoorvajagtap/fileStore.git
$ cd fileStore
$ go run server/server.go           
```

- On another terminal, build the CLI and test the following:
```
$ go build -o ./store client/client.go
$ ./store add <file_names>
$ ./store ls
$ ./store rm <file_name>
$ ./store wc
$ ./store update <file_name>
...
```