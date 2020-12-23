
# Aquarium GCP Janitor

## About Janitor
Janitor is a cloud-based schema management system for schema registration, versioning and recovery. It enables developers to define and manage standard schemas for various events by storing their versioned history.

Janitor is developed on Google Cloud Platform in Go programming language and is deployed in a docker container.

Product functionalities include message schema retrieval, registration, versioning and reconstruction (evolution).

**Technology stack**
- Programming language: Go
- Compute: GCP Cloud Functions, Cloud Run
- Big data: GCP Storage, Pub/Sub
- Database: FireStore
- Tools: Cloud Build, Container Registry, Cloud Scheduler, Cloud Tasks


Janitor has the following components:
1. `Schema Registry` component with a database and a RESTful interface - the main component which communicates with the server and database that holds the collection of all schemas and serves the purpose of retrieving, registering, validating and versioning schemas
2. `Central Consumer` component - a cloud function that pulls new messages and, depending on their matching with the schema, marks them as valid, invalid, or dead-letter; all with the help of communication with the Schema Registry
3. `Puller&Cleaner` component - a cloud function that pulls invalid-marked messages and tries to recover them by inferring new schema; all with the help of communication with the Schema Registry

**More details about Janitor can be found on**

* [Janitor Wiki Page](https://github.com/syntio/aquarium-janitor-gcp/wiki)

## Dependencies

**GCP:** 
- Existing or new project
- Enabled billing for your Cloud project
- Enabled Cloud Functions, Cloud Pub/Sub, Cloud Build, Cloud Build, Cloud Scheduler, Cloud Tasks, Cloud Run Admin, Cloud Firestore APIs

**Programming language:** Go 1.13

## Features
- **Schema registration** enables users to register a new schema, and retrieve them if needed. By providing ID and version of the schema users manipulate with schemas stored in the database.
- **Schema recovery** includes inferring schemas from invalid-marked messages (if possible), and then registering the new schema.

**Developed on** Google Cloud Platform

**Supports the following:**
* storage systems: Firestore
* messaging systems: Cloud Pub/Sub
* message formats:
    - schema registration: JSON, CSV, XML, AVRO, Protobuf
    - schema evolution: JSON, CSV 


## Repository structure
```bash
+---central-consumer
|       configuration
|       pubsub
|       registry
|       validator
|       function.go
|       go.mod
|
+---helper-functions
|       csv-validator
|       xml-validator
|
+---puller-cleaner
|       puller-cleaner-csv
|       puller-cleaner-json
|
+---puller-tasks
|       puller-tasks-csv
|       puller-tasks-json
|
+---schema-registry
|       business-logic
|       configuration
|       database
|       main
|       model
|       rest
|       schema_creation
|       util
|       Dockerfile
|       go.mod
|       jsonSchemaDynamicCreation.py
|   .gitignore
|   CODE_OF_CONDUCT.md
|   config.yaml
|   cloudbuild.yaml
|   CONTRIBUTING.md
|   deploy.sh
|   LICENSE.md
|   README.md
\   SUPPORT_INFORMATION.md
```

## Deployment and Configuration
Instructions on how to configure your GCP project and deploy the GCP janitor can be found [here](../../wiki/deployment-and-configuration).
## Usage
* Message schema registration
* Message schema retrieval
* Message schema evolution


## Links
Issue tracker: https://github.com/syntio/aquarium-janitor-gcp/issues

*In case of sensitive bugs like security vulnerabilities, please contact support@syntio.net directly instead of using issue tracker. We value your effort to improve the security and privacy of this project!*

## Contributing
Please refer to [CONTRIBUTING.md](./CONTRIBUTING.md)
## Developed by
The product is developed and maintained by Syntio Labs
## License
Licensed under the Apache License, Version 2.0
