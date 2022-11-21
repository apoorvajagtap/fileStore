
## File Store server
###### Work In Progress

## Current Status
Following tasks have been achieved:
- Add multiple text files to the server. Skips saving the file, if already exists.
- Delete the specified file.
- List all the filenames in the server.
- Update the specified file, or create new if it does not exist.
- Count the total words in all the files stored on the server. 

## Steps to Run Locally using OpenShift cluster.

*Pre-requisite*  
- You have a functional OpenShift cluster

*Start the server*
- Use OCP_manifests/template.yaml to create the required objects in the cluster.
```
$ oc create -f <OCP_manifests/template.yaml>
$ oc get pods
// must be running mongodb (2 replicas) & 1 filestore pod
```

- On local machine, clone the repo & build the CLI and test the following (copy the route-hostname of filestore server on OpenShift):
```
$ go build -o ./store client/client.go
$ ./store --url=<route_hostname> add <file_names>
$ ./store --url=<route_hostname> ls
$ ./store --url=<route_hostname> rm <file_name>
$ ./store --url=<route_hostname> wc
$ ./store --url=<route_hostname> update <file_name>
...
```