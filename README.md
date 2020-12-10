
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
The following section includes instructions on deploying and setting up the environment for Janitor's functionalities. The Schema Registry component is put in a docker container and the container is deployed to the cloud using Google's Cloud Build service. The Central Consumer and Puller&Cleaner components are deployed using Cloud Functions.

Janitor can be deployed **via the deployment script** (through Cloud Build) or **manually** by deploying all of the components individually.

### Deployment via script
The deployment process includes a Cloud Build service which executes the bash script. The script creates all the needed resources and deploys the Janitor components to the Google Cloud Platform.

#### Prerequisites
Before the script starts the deployment process, it is necessary to do the following:
* [Enable the Cloud Build API](https://cloud.google.com/endpoints/docs/openapi/enable-api#console)
* Add the required roles to your Cloud Build service account
    * from the IAM & Admin page select the Cloud Build service account (ends with `cloudbuild.gserviceaccount.com`) and go in to edit
    * in the edit member section add the following roles:
        `Compute Admin, Compute Instance Admin, Cloud Functions Admin, Pub/Sub editor, Cloud Run Admin, Storage Admin, Cloud Tasks Admin and Cloud Scheduler Admin`.
This gives Cloud Build the permissions needed to create all the resources. 
* Customize `config.yaml` file to suit your needs, the description of the parameters in the file is as follows:

Parent field | Field (parameter) | Description | Type | Example
--- | --- | --- | --- |---
`topics` | `invalidTopicJSON` | ID of the Pub/Sub topic where the JSON invalid messages are sent | `string` | _"invalid-topic-json"_
<i></i> | `invalidTopicCSV` | ID of the Pub/Sub topic where the CSV invalid messages are sent | `string` | _"invalid-topic-csv"_
<i></i> | `validTopic` | ID of the Pub/Sub topic where the valid messages are sent | `string` | _"valid-topic"_
<i></i> | `deadLetterTopic` | ID of the Pub/Sub topic where the dead-letter messages are sent | `string` | _"dead-letter-topic"_
<i></i> | `inputTopic` | ID of the Pub/Sub topic where the publisher sends messages | `string` | _"input-topic"_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
`subscriptions` | `subIdInput` | ID of the Pub/Sub subscription with which the messages are retrieved from the input topic | `string` | _"input-topic-sub"_
<i></i> | `subIdValid` | ID of the Pub/Sub subscription with which the messages are retrieved from the valid topic | `string` | _"valid-topic-sub"_
<i></i> | `subIdDeadletter` | ID of the Pub/Sub subscription with which the messages are retrieved from the dead-letter topic | `string` | _"dead-letter-topic-sub"_
<i></i> | `subIdJSON` | ID of the Pub/Sub subscription with which the messages are retrieved from the JSON invalid topic | `string` | _"invalid-topic-json-sub"_
<i></i> | `subIdCSV` | ID of the Pub/Sub subscription with which the messages are retrieved from the CSV invalid topic | `string` | _"invalid-topic-csv-sub"_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
`functions` | `schemaRegistryURL` | URL of the Schema Registry's Cloud Run container | `string` | _"https://<i></i>schema-registry-ew.a.run.app"_
<i></i> | `csvValidatorURL` | URL of the auxiliary function that validates the CSV message against provided schema | `string` | _"https://<i></i>janitor-project.cloudfunctions.net/ csv-validator"_
<i></i> | `xmlValidatorURL` | URL of the auxiliary function that validates the XML message | `string` | _"https://<i></i>janitor-project.cloudfunctions.net/ xml-validator"_
<i></i> | `schemaRegistry  EvolutionURL` | URL of the Schema Registry's Cloud Run container for communication regarding the evolution | `string` | _"https://<i></i>schema-registry-ew.a.run.app/<br>schema/%s/evolution"_
<i></i> | `pullerCleaner  JsonURL` | HTTP trigger of the _JSON Puller&Cleaner_ | `string` | _"https://<i></i>janitor-project.cloudfunctions.net/ puller-cleaner-json"_
<i></i> | `pullerCleaner  CsvURL` | HTTP trigger of the _CSV Puller&Cleaner_ | `string` | _"https://<i></i>janitor-project.cloudfunctions.net/ puller-cleaner-csv"_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
`protoparam` | `tmpFilePath` | path to the `tmpFileName`'s directory | `string` | _"/tmp"_
<i></i> | `tmpFileName` | name of the auxiliary file for parsing the Protobuf message | `string` | _"tmp.txt"_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
`pullercleaner  json` | `timeDuration  Seconds` | execution time in seconds of the _JSON Puller&Cleaner_'s pulling part | `int` | _15_
<i></i> | `maxBatchSize` | maximum number of messages that the _JSON Puller&Cleaner_ can pull from the JSON/CSV invalid topic | `int` | _4000_
<i></i> | `maxThroughput` | maximum number of bytes that the _JSON Puller&Cleaner_ can pull from the JSON/CSV invalid topic | `int` | _26214400_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
`pullercleaner  csv` | `timeDuration  Seconds` | execution time in seconds of the _CSV Puller&Cleaner_'s pulling part | `int` | _30_
<i></i> | `maxBatchSize` | maximum number of messages that the _CSV Puller&Cleaner_ can pull from the JSON/CSV invalid topic | `int` | _500_
<i></i> | `maxThroughput` | maximum number of bytes that the _CSV Puller&Cleaner_ can pull from the JSON/CSV invalid topic | `int` | _26214400_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
_None_ | `pullerCleaner  TaskQueue` | name of the task queue that periodically runs the _Puller&Cleaner_ | `string` | _"puller-cleaner-queue"_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
_None_ | `contentType` | type of HTTP response | `string` |  _"application/json"_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
_None_ | `fileMode` | file's mode and permission bits | `int` | _0644_
<i></i> | <i></i> | <i></i> | <i></i> | <i></i>
_None_ | `firebase  CollectionName` | name of the collection in the FireStore database | `string` |  _"Registry"_


#### Deployment
Open the Cloud Shell and clone the GitHub repository:
```shell
git clone https://github.com/syntio/aquarium-janitor-gcp.git
```

Enter the created directory and run the following command to start the script execution:
```shell
gcloud builds submit --config=cloudbuild.yaml
```

The environment variables passed to the script and used in the deployment process are:
* `PROJECT_ID` - project name
* `REGION` - the world region where the components will be deployed to
* `BUCKET_NAME` - the name of the storage bucket where the configuration file and the Cloud Functions source code will be stored
* `CONFIG_FILE` - the name of the configuration file

The script will create the resources in the following order:
1) **Create the Storage bucket and upload the configuration file**
2) **Create topics and their subscriptions**
3) **Create a FireStore database**
4) **Zip the Cloud Functions source codes and upload them to the storage bucket**
5) **Build the image for Schema Registry component** - the app is put into a container and the image is deployed to the Container Registry
6) **Deploy the container image to Cloud Run** - the image in the Container Registry is deployed to Cloud Run
7) **Deploy the Central Consumer to Cloud Functions**
8) **Deploy the Puller & Cleaner component to Cloud Functions**
9) **Deploy the Puller Tasks function** - function triggered with Cloud Scheduler used to create Cloud Tasks that invoke the Puller & Cleaner function
10) **Deploy the helper functions** - functions for XML and CSV validation
11) **Create Cloud Task Queue** - queue for asynchronous task completion
12) **Create Cloud Scheduler jobs** - schedule Cloud Function invocation at regular intervals for task creation

### Manual deployment
Manual deployment of the components, if needed, can be done using the following commands:
1. **Schema Registry Component**
    * Create a Docker container image
        
        The following command builds the container image and deploys it in the Container Registry:
        ```shell
        gcloud builds submit --tag gcr.io/$PROJECT-ID/schema-registry
        ```
        *`$PROJECT_ID` is the name of your Google Cloud project*
    * Deploy the container image to Cloud run
        ```shell
        gcloud run deploy $SERVICE_NAME --runtime=go113 --image gcr.io/$PROJECT_ID/schema-registry --platform managed --region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,BUCKET_NAME=$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE
        ```

        *`$SERVICE_NAME` is your name of choice for the Cloud Run service, `$PROJECT_ID` the name of your project and `$REGION` is the name of the region you want to deploy the service to. Use the --set-env-vars parameter to set up the environment variables for the service in the form of key value pairs.*

    * You will be prompted for the service name: press Enter to accept the default name.
    * You will be prompted to allow unauthenticated invocations: respond n .

    Then wait a few minutes until the deployment is complete. On success, the command line displays the service URL. You can visit your deployed container by opening the service URL in a web browser.    

2. **Central Consumer** 

    The following command deploys the Cloud function which is triggered by a topic. Every time a message is sent to the topic, a new instance of the function is invoked.
    ```shell
    gcloud functions deploy central-consumer --source $SOURCE --runtime go113 --entry-point CentralConsumerHandler --trigger-topic $INPUT_TOPIC –-region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,BUCKET_NAME=$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE

    ```
3. **Puller & Cleaner**
    
    The following command deploys the Cloud function which is triggered by HTTP:
    ```shell
    gcloud functions deploy puller-cleaner-json --source $SOURCE --runtime go113 --entry-point PullerCleaner --trigger-http –-region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,BUCKET_NAME=$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE 
    ```
4. **Creating Cloud Scheduler jobs**

    Create a job that invokes the Cloud Function for Cloud Tasks creation in regular intervals.
    ```shell
    gcloud scheduler jobs create http puller-cleaner-invoker-json --schedule "* * * * *" --uri=$SCHEDULER_URI --http-method GET
    ```

    where `$SCHEDULER_URI` is the name of the Cloud Function that we want to invoke. The `schedule` parameter is used to set the frequency at which the function will be invoked, specified using the unix-cron format.

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
