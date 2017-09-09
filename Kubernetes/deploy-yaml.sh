#!/bin/sh
sed -e "s/COMMIT_SHA/`git rev-parse origin/master`/g" status-api/deployment-template.yaml > status-api/deployment.yaml
sed -e "s/COMMIT_SHA/`git rev-parse origin/master`/g" mysql-proxy/deployment-template.yaml > mysql-proxy/deployment.yaml
sed -e "s/COMMIT_SHA/`git rev-parse origin/master`/g" nginx-routing/deployment-template.yaml > nginx-routing/deployment.yaml
sed -e "s/COMMIT_SHA/`git rev-parse origin/master`/g" novels-api/deployment-template.yaml > novels-api/deployment.yaml
sed -e "s/COMMIT_SHA/`git rev-parse origin/master`/g" senteces-api/deployment-template.yaml > senteces-api/deployment.yaml
sed -e "s/COMMIT_SHA/`git rev-parse origin/master`/g" novelists-api/deployment-template.yaml > novelists-api/deployment.yaml
