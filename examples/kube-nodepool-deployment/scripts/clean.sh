#!/bin/bash

# kill still running processes in case of errors
ps -ef| grep -i "sleep 10" | grep -v grep| awk '{print "kill "$2}' | sh
ps -ef| grep -i "create" | grep -v grep| awk '{print "kill "$2}' | sh

# Destroy/Remove created resources
terraform destroy -auto-approve

# Clean logs
> ../logs

