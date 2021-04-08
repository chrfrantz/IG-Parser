#!/bin/bash

# This file is part of the IG-Parser project.
# Maintainer: Christopher Frantz (cf@christopherfrantz.org)

# Clears current instance of IG-Parser and redeploys after pulling from git repository.
# Ensure to run with sudo

# Tear down current version 
docker-compose down

# Shortcut to clean up generated docker image
docker rmi ig-parser_web

# Pull latest version
git pull

# Deploy
docker-compose up -d

echo "Latest version should be deployed (check output for errors)"
