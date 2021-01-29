## Launching this repository in a linux host through VSCode Remote Containers

### Requirements

- VSCode
- `ms-vscode-remote.remote-containers` extension in your host VSCode.
- Linux: Docker CE/EE 18.06+ and Docker Compose 1.21+

### Quickstart

From your host machine run:

```
./allow-ports.sh
```

Then using VSCode Command `Remote Containers: Open folder in container...` open this repository root folder.

You are ready to develop.

### Explanation

You can start the remote container extension either by clicking on the bottom left side icon or by opening the command palette and searching for `Remote Containers: Open folder in container...`

Select the root folder from this repository.

And voilà, if you open a new terminal inside VSCode you can check that you have our fully featured development environment already set up. For instance you can check the node version by running `node --version` or `npm version`. Maybe you could run also `go version`.

You can also `cd go/testutils` and run `./prepareMarketTests.sh` and you will see that the containers from that docker-compose will start.

Let's try to run some tests then `cd go/notary/consumer && go test ./... -v`.

Oops, timeout in localhost. If only it were that easy...

So, what's happening?

We are running our integration tests against the docker containers that we have started, but we are running them from inside the `crypto-soccer_devcontainer_dockerhost_1` container which doesn't know how to resolve http://localhost:{PORT} to reach those other containers, like the
PostgresDB or the ethereum local network.

If we were in MacOS or Windows we could replace `http://localhost:5000` by `http://host.docker.internal:5000` and everything would run smoothly. But as of today this feature is not implemented in Docker for Linux.

So, what do we do?

Yeah, we could run our tests from our host machine provided that we have `golang` environment installed, but that defeats the purpose of having your entire environment inside a Docker Container in order to have your host machine as light as a feather or maybe with another `golang` version, so... let's move on

On top of spinning up our custom development environment container called `crypto-soccer_devcontainer_dockerhost_1` we also spin up a container with the image `qoomon/docker-host` which will allow our requests from inside a container reach our host localhost where the other containers are exposed.

So, now every time that we wish to reach localhost we should point it to `http://crypto-soccer_devcontainer_dockerhost_1`

Onn top of that we need to enable firewall rules to allow dockerhost to connect to your host machine.

That could be done through UFW

```
sudo su
ipp=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' crypto-soccer_devcontainer_dockerhost_1) && ufw allow from $ipp to any port 8545 && ufw allow from $ipp to any port 5432

```

Or by running `.devcontainer/allow-ports.sh` which allows all the ports exposed in the containers started by test environments.

Once done you can safely run all the integration tests. And you are fully set up to build new features.

### TODO

- Configure git credentials (?)
