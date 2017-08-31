#!/bin/sh
kubectl replace -f senteces-api/deployment.yaml && \
kubectl replace -f novels-api/deployment.yaml && \
kubectl replace -f mysql-proxy/deployment.yaml && \
kubectl replace -f nginx-routing/deployment.yaml && \
kubectl replace -f status-api/deployment.yaml
