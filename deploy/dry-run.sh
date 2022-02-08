#!/bin/sh

kubectl apply -k . --dry-run=client
