#!/bin/bash

# This file is part of the IG Parser project.
# Maintainer: Christopher Frantz (cf@christopherfrantz.org)

# Clears current instance of IG Parser and redeploys after pulling from git repository.
# Ensure to run the script with sudo!

echo "Initiating (re)deployment of latest version of the IG Parser ..."

# Create logs folder if it does not already exist
if [ ! -d ./logs ]; then
  echo "Creating './logs' folder ..."
  if [ "$(mkdir -p ./logs)" -ne 0 ]; then
    echo "Error during folder creation." 
    exit 1
  fi
fi

# Tear down current version 
echo "Undeploying running IG Parser instance ..."
if [ "$(docker-compose down)" -ne 0 ]; then
  echo "Stopping of service." 
  exit 1
fi

# Pull latest version
echo "Retrieving latest version ..."
if [ "$(git pull)" -ne 0 ]; then
  echo "Error when pulling latest version of repo." 
  exit 1
fi

# Create docker network if not already existing
echo "Checking network setup ..."
if [ "$(docker network inspect tunnel_network > /dev/null)" -ne 0 ]; then
  echo "Dedicated docker network does not exist. Creating 'tunnel_network' ..."
  if [ "$(docker network create tunnel_network)" -ne 0 ];  then
    echo "Network creation failed."
    exit 1
  fi
fi

# Deploy
echo "Deploying latest version of IG Parser ..."
if [ "$(docker-compose up -d --build)" -ne 0 ]; then
  echo "Error during deployment." 
  exit 1
fi

# Remove remaining build containers
echo "Removing all intermediate images created during build process ..."
if [ "$(docker image prune --filter label=stage=builder -f)" -ne 0 ]; then
  echo "Error when deleting intermediate images created during build."
  exit 1
fi

echo "Latest version should be deployed (check output for potential errors)"
