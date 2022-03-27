#!/bin/bash
kubectl create secret generic dotenv --from-env-file=../.env
kubectl get secret dotenv
