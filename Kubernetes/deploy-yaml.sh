#!/bin/sh
sed -e "s/COMMIT_HASH/`git rev-parse origin/master`/g" errors-api/deployment-template.yaml > errors-api/deployment.yaml
sed -e "s/COMMIT_HASH/`git rev-parse origin/master`/g" mysql-proxy/deployment-template.yaml > mysql-proxy/deployment.yaml
sed -e "s/COMMIT_HASH/`git rev-parse origin/master`/g" nginx-routing/deployment-template.yaml > nginx-routing/deployment.yaml
sed -e "s/COMMIT_HASH/`git rev-parse origin/master`/g" novels-api/deployment-template.yaml > novels-api/deployment.yaml
sed -e "s/COMMIT_HASH/`git rev-parse origin/master`/g" senteces-api/deployment-template.yaml > senteces-api/deployment.yaml
