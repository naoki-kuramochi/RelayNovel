#!/bin/sh
kubectl autoscale deployment novels-api --cpu-percent=50 --min=1 --max=10
kubectl autoscale deployment mysql-proxy --cpu-percent=50 --min=1 --max=10
kubectl autoscale deployment sentences-api --cpu-percent=50 --min=1 --max=10
kubectl autoscale deployment nginx-routing --cpu-percent=50 --min=1 --max=10
kubectl autoscale deployment status-api --cpu-percent=50 --min=1 --max=10
