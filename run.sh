#!/bin/bash
echo "Building k8s-controller..."
go build -o k8s-controller main.go

echo "Starting server with controller-runtime..."
./k8s-controller serve
