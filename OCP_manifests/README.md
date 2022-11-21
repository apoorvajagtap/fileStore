The `template.yaml` file would create the following:
- A custom SCC for mongodb application.
- Project named as `filestore-server`
- Statefulset that hosts mongodb, and svc to expose the same.
- Deployment for filestore server, svc, and route to expose the same.

##### *Steps to run the application on OpenShfit cluster:*
- Copy the content of `template.yaml` to a local file, and execute the following:

```
$ oc create -f <yaml_file_path>
```