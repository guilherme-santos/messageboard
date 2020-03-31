# Back's Message Board

This is a RESTFul API which allow unauthenticated users add new messages and authenticated users list, update and get details about a specific one.

We're using Makefile to make our life easier, you can check all available targets typing `make` or `make help`.

### Building

You can build the project using `make build` which will generate a image with `latest` tag, if you want to set a new tag you can use:

```shell
$ make build version=1.2.3
```

We're using multi-stage building, it means that we have one image called `buider` to build everything and also for development, and the final image which is an `alpine:3.11` with the our binary installed. Check [Dockerfile](./Dockerfile) for more details.

### Running

To run the services is analogous to build it, just type `make up`. If you need to stop them type `make down`.

A quick restart could be done with `make down up`

We also run a instance of [MongoDB](https://www.mongodb.com/), which was chose because I like NoSQL :) and we didn't have any need for a SQL database, or any specific feature, like transactions for example.

### Testing

If you need to run the tests you can type: `make test`, by default we run with the flag `-race`.

You can always override the default test command, here some examples:

```shell
$ make test flag=-v
$ make test flag=-v testcase="TestMessageBoardHandler"
$ make test testflag=-v # it'll remove -race
```

You can also check the code coverage typing `make test-coverage`

### Initial load

By default the file [messages.csv](./messages.csv) is loaded when the service start, if you need to load a different file, you need to map the new file to inside of the container, for example:

```yaml
# docker-compose.yml
    ...
    volumes:
      - <path-for-your-new-file>:/etc/messageboard/messages.csv
```

As you can see you can map a file in your local machine to inside of the container. By default we always use `/etc/messageboard/messages.csv` but this can also be changed, updating the environment variable `MONGODB_INITIAL_CSV` inside of your [docker-compose.yml](./docker-compose.yml).

***IMPORTANT*** every time the container goes up, we clear the whole database and load the CSV file.

### Accessing the API

Our API exports 4 endpoints:

- **POST /v1/messages**: create a new message (*public*)
- **GET /v1/messages**: list all messages, you can control pagination using `per_page` and `page` query strings (*private*)
- **GET /v1/messages/{id}**: get a specific message (*private*)
- **PUT /v1/messages/{id}**: update a specific message (*private*)

For the private endpoints we're using http basic auth, but other more secure ways should be implemented like JWT. The users available to the private endpoints could be configured in the [docker-compose.yml](./docker-compose.yml) as well. You have to change the environment variable `CREDENTIALS` which accept multiples users separated by comma. For example: `user1:pass-user-1,user2:pass-user-2,user3:pass-user-3`

### Developing

We provide a example of docker-compose.override to help during the development, it will allow you run the container once, change your code and run it again, making the development process way faster.

If you want to use our default docker-compose.override, type:

```shell
$ cp docker-compose.override.yml{.dist,}
```

Then you can build and run normally:

```shell
$ make build up
```

Once the container is running, you can open a shell inside of it, kill the `messageboard` process, as in this environment our service does not have the PID 1, the container will continue running, so you can change your code locally, and type `make run`.
