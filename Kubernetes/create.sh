#!/bin/sh
kubectl create -f senteces-api/deployment.yaml && \
kubectl create -f novels-api/deployment.yaml && \
kubectl create -f mysql-proxy/deployment.yaml && \
kubectl create -f nginx-routing/deployment.yaml
#kubectl create -f senteces-api/service.yaml && \
#kubectl create -f novels-api/service.yaml && \
#kubectl create -f mysql-proxy/service.yaml && \
#kubectl create -f nginx-routing/service.yaml
