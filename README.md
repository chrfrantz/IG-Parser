# IG-Parser
Parser for IG 2.0 Statements

## Deployment

* Prerequisites:
  * [Docker](https://docs.docker.com/engine/install/)
  * [docker-compose](https://docs.docker.com/compose/install/) (optional if environment, volume and port parameterization is done manually)

* Deployment Guidelines
  * Start the IG Script Parser service by running `docker-compose up -d`. Run `docker-compose down` to stop execution.
  * By default, the web service listens on port 4040, and logging is enabled in the subfolder `./logs`. 
  * The service automatically restarts if it crashes. 
  * Adjust the docker-compose.yml file to modify any of these characteristics.
