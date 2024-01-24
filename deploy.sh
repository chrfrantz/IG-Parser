#!/bin/bash

# This file is part of the IG Parser project (URL: https://github.com/chrfrantz/IG-Parser).
# Maintainer: Christopher Frantz (cf@christopherfrantz.org)

# Clears current instance of IG Parser and redeploys after pulling from git repository.
# Ensure to run the script with sudo (e.g., 'sudo ./deploy.sh')!

echo "Initiating (re)deployment of latest version of the IG Parser ..."

# Create logs folder if it does not already exist
if [ ! -d ./logs ]; then
  echo "Creating './logs' folder ..."
  mkdir -p ./logs;
  if [ $? -ne 0 ]; then
    echo "Error during folder creation. Ensure you are running the script with superuser permissions (e.g., sudo)." 
    exit 1
  fi
fi

# Tear down current version 
echo "Undeploying running IG Parser instance ..."
docker-compose down
if [ $? -ne 0 ]; then
  echo "Error when stopping service. Ensure you are running the script with superuser permissions (e.g., sudo)." 
  exit 1
fi

# Delete runtime image generated as part of previous build process
echo "Deleting image of undeployed IG Parser instance ..."
docker image prune --filter label=stage=runner -f
if [ $? -ne 0 ]; then
  echo "Error when deleting runner images following undeploying. Ensure you are running the script with superuser permissions (e.g., sudo)."
  exit 1
fi

# Pull latest version
echo "Retrieving latest version ..."
git pull
if [ $? -ne 0 ]; then
  echo "Error when pulling latest version of repo. Ensure you are running the script with superuser permissions (e.g., sudo)." 
  exit 1
fi

# Create docker network if not already existing
echo "Checking network setup ..."
docker network inspect tunnel_network > /dev/null
if [ $? -ne 0 ]; then
  echo "Dedicated docker network does not exist. Creating 'tunnel_network' ..."
  docker network create tunnel_network
  if [ $? -ne 0 ];  then
    echo "Network creation failed. Ensure you are running the script with superuser permissions (e.g., sudo)."
    exit 1
  fi
fi

# Deploy
echo "Deploying latest version of IG Parser ..."
docker-compose up -d --build
if [ $? -ne 0 ]; then
  echo "Error during deployment. Ensure you are running the script with superuser permissions (e.g., sudo)." 
  exit 1
fi

# Remove remaining build containers
echo "Removing all intermediate images created during build process ..."
docker image prune --filter label=stage=builder -f
if [ $? -ne 0 ]; then
  echo "Error when deleting intermediate images created during build. Ensure you are running the script with superuser permissions (e.g., sudo)."
  exit 1
fi

echo "Completed deployment of latest IG Parser version (check output for potential errors)"
