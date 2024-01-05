# gRPC Hello World Delegate to Specified IP

## Run in Kind

* Set up a kind cluster with Istio.
* Deploy the backend server with `kubectl apply -f server-deployment.yaml`
* Base64 encode the IP address of one of the backend pods: `kubectl get pod -l app=grpc-server -o json | jq -jr '.items[0].status.podIP' | base64`
* Set this as the `--ip` argument in `client-job.yaml`

## Build Docs

Follow these setup to run the [quick start][] example:

 1. Get the code:

    ```console
    $ go get google.golang.org/grpc/examples/helloworld/greeter_client
    $ go get google.golang.org/grpc/examples/helloworld/greeter_server
    ```

 2. Run the server:

    ```console
    $ $(go env GOPATH)/bin/greeter_server &
    ```

 3. Run the client:

    ```console
    $ $(go env GOPATH)/bin/greeter_client
    Greeting: Hello world
    ```

For more details (including instructions for making a small change to the
example code) or if you're having trouble running this example, see [Quick
Start][].

[quick start]: https://grpc.io/docs/languages/go/quickstart
