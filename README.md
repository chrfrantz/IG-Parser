# IG-Parser
Parser for IG 2.0 Statements

## Deployment

* Prerequisites:
  * [Docker](https://docs.docker.com/engine/install/)
  * [docker-compose](https://docs.docker.com/compose/install/) (optional if environment, volume and port parameterization is done manually)
  * Quick installation of docker under Ubuntu LTS: `sudo snap install docker --classic`

* Deployment Guidelines
  * Clone this repository
  * Make `deploy.sh` executable (`chmod 740 deploy.sh`)
  * Run `deploy.sh` with superuser permissions (`sudo deploy.sh`)
    * This script deletes old versions of IG-Parser, before pulling the latest version and deploying it.
    * For manual start, run `sudo docker-compose up -d`. Run `sudo docker-compose down` to stop execution.
  
* Service Configuration
  * By default, the web service listens on port 4040, and logging is enabled in the subfolder `./logs`. 
  * The service automatically restarts if it crashes. 
  * Adjust the docker-compose.yml file to modify any of these characteristics.
