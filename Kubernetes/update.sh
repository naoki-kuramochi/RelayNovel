#!/bin/sh
kubectl update -f senteces-api/deployment.yaml && \
kubectl update -f novels-api/deployment.yaml && \
kubectl update -f mysql-proxy/deployment.yaml && \
kubectl update -f nginx-routing/deployment.yaml
