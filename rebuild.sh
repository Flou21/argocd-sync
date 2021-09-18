#!/bin/bash

docker build -t registry.flou.dev/argocd-sync .
docker push registry.flou.dev/argocd-sync
