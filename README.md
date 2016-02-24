A docker API-compatible server for Cloud Foundry.

# Quickstart

## Target a CF, for example with MicroPCF

1. Follow the instructions at https://github.com/pivotal-cf/micropcf/
1. Start it, target it, log in etc. with the `cf` cli

## Start cf-docker-bridge (you can do this in your .profile if you like)

~~~~
cf-docker-bridge up # listens on socker.dock, forwards requests to cloud foundry via the CF cli
export DOCKER_HOST=unix://$PWD/socker.dock
~~~~

## Now run docker as normal, and use CF to manage your apps

~~~~
docker run -d busybox --name mydockerapp # prints created CF app name
cf apps # shows your docker container, now you can scale it and interact with it as normal with the cf CLI
~~~~

# Advice for Production Use

don't.
