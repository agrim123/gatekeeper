#/bin/bash

set -eoux pipefail

docker build . -t gatekeeper

docker rm gatekeeper || exit 0

docker run -d --name gatekeeper gatekeeper

docker cp gatekeeper:/opt/gatekeeper/gatekeeper .

docker kill gatekeeper
