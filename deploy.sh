#!/bin/bash

# This file is part of the IG-Parser project.
# Maintainer: Christopher Frantz (cf@christopherfrantz.org)

# Clears current instance of IG-Parser and redeploys after pulling from git repository.
# Ensure to run the script with sudo!

echo "Initiating (re)deployment of latest version of IG-Parser ..."

# Create logs folder if it does not already exist
if [ ! -d ./logs ]; then
  mkdir -p ./logs;
  if [ $? -ne 0 ]; then
    echo "Error during folder creation." 
    exit 1
  fi
fi


# Tear down current version 
docker-compose down
if [ $? -ne 0 ]; then
  echo "Stopping of service." 
  exit 1
fi

# Pull latest version
git pull
if [ $? -ne 0 ]; then
  echo "Error when pulling latest version of repo." 
  exit 1
fi

# Create docker network if not already existing
docker network inspect tunnel_network > /dev/null
if [ $? -ne 0 ]; then
  echo "Network does not exist. Creating it ..."
  docker network create tunnel_network
  if [ $? -ne 0 ];  then
    echo "Network creation failed."
    exit 1
  fi
fi

# Deploy
docker-compose up -d --build
if [ $? -ne 0 ]; then
  echo "Error during deployment." 
  exit 1
fi

# Remove remaining build containers
docker image prune --filter label=stage=builder
if [ $? -ne 0 ]; then
  echo "Error when deleting build containers."
  exit 1
fi

echo "Latest version should be deployed (check output for errors)"
