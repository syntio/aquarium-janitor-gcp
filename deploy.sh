# Deployment script for the Janitor product

# Fetch the parameters from the environment variables
echo "Environment variables:"
echo "Project_id: " $PROJECT_ID
echo "Region: " $REGION
echo "Bucket_name: " $BUCKET_NAME
echo "Config_file: " $CONFIG_FILE

# Install the wget,zip and Go packages
apt-get update
apt-get install wget
apt install zip unzip

wget https://dl.google.com/go/go1.15.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.15.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install the yq package for parsing the YAML configuration file
GO111MODULE=on go get github.com/mikefarah/yq/v3
export PATH=$PATH:/builder/home/go/bin



# Parse the configuration file and get the config values
echo "Parsing the configuration file.."
INPUT_TOPIC=$(yq r $CONFIG_FILE topics.inputTopic)
VALID_TOPIC=$(yq r $CONFIG_FILE topics.validTopic)
INVALID_TOPIC_JSON=$(yq r $CONFIG_FILE topics.invalidTopicJSON)
INVALID_TOPIC_CSV=$(yq r $CONFIG_FILE topics.invalidTopicCSV)
DEADLETTER_TOPIC=$(yq r $CONFIG_FILE topics.deadLetterTopic)

INPUT_SUB=$(yq r $CONFIG_FILE subscriptions.subIdInput)
VALID_SUB=$(yq r $CONFIG_FILE subscriptions.subIdValid)
INVALID_SUB_JSON=$(yq r $CONFIG_FILE subscriptions.subIdJSON)
INVALID_SUB_CSV=$(yq r $CONFIG_FILE subscriptions.subIdCSV)
DEADLETTER_SUB=$(yq r $CONFIG_FILE subscriptions.subIdDeadletter)

echo "Topic for input messages: " $INPUT_TOPIC
echo "Topic for valid messages: " $VALID_TOPIC
echo "Topic for invalid JSON messages: " $INVALID_TOPIC_JSON
echo "Topic for invalid CSV messages: " $INVALID_TOPIC_CSV
echo "Topic for deadletter messages: " $DEADLETTER_TOPIC
echo "Subscription for input topic: " $INPUT_SUB
echo "Subscription for valid topic: " $VALID_SUB
echo "Subscription for invalid JSON topic: " $INVALID_SUB_JSON
echo "Subscription for invalid CSV topic: " $INVALID_SUB_CSV
echo "Subscription for deadletter topic: " $DEADLETTER_SUB

# Create a bucket for configuration file storing
echo "Creating the storage bucket.."
gsutil mb -l $REGION gs://$PROJECT_ID-$BUCKET_NAME

echo "Storing the configuration file on GCP bucket.."
gsutil cp $CONFIG_FILE gs://$PROJECT_ID-$BUCKET_NAME

# Create the topics:
# - input topic: the entry point for all the messages
# - valid_topic: topic for messages with valid schemas
# - invalid topic: topic for messages with invalid schemas
# - deadletter topic: topic for messages with missing schema structure or with errors
echo "Creating topics.."
gcloud pubsub topics create $INPUT_TOPIC $VALID_TOPIC $INVALID_TOPIC_JSON $INVALID_TOPIC_CSV $DEADLETTER_TOPIC

# Wait until all the topics have been created

# Create the subscriptions
echo "Creating subscriptions.."
gcloud pubsub subscriptions create $INPUT_SUB --topic $INPUT_TOPIC
gcloud pubsub subscriptions create $VALID_SUB --topic $VALID_TOPIC
gcloud pubsub subscriptions create $INVALID_SUB_JSON --topic $INVALID_TOPIC_JSON
gcloud pubsub subscriptions create $INVALID_SUB_CSV --topic $INVALID_TOPIC_CSV
gcloud pubsub subscriptions create $DEADLETTER_SUB --topic $DEADLETTER_TOPIC

# Create the Firestore data store for schema information storing
echo "Creating Firestore data store.."
gcloud alpha firestore databases create --region=$REGION --project=$PROJECT_ID

temp=$(gcloud iam service-accounts list | grep $PROJECT_ID@appspot.gserviceaccount.com)
tempList=($temp)
AppServiceAccount=${tempList[5]}
echo $AppServiceAccount
gcloud iam service-accounts keys create $PWD/key.json --iam-account=$AppServiceAccount

echo "Stroing the iam key on the GCP bucket"
gsutil cp $PWD/key.json gs://$PROJECT_ID-$BUCKET_NAME
# Create .zip source files for Cloud Functions
echo "Creating .zip source files.."
cd central-consumer;
zip -r ../central-consumer.zip *
cd ../puller-cleaner/puller-cleaner-json;
zip -r /workspace/puller-cleaner-json.zip *
cd ../puller-cleaner-csv;
zip -r /workspace/puller-cleaner-csv.zip *
cd  ../../puller-tasks/puller-tasks-json;
zip -r /workspace/puller-tasks-json.zip *
cd ../puller-tasks-csv
zip -r /workspace/puller-tasks-csv.zip *
cd ../../helper-functions/xml-validator;
zip -r /workspace/xml-validator.zip *
cd ../csv-validator
zip -r /workspace/csv-validator.zip *

cd /workspace

# Upload the source files to the storage bucket
echo "Uploading Cloud Functions source code to storage.."
for zip in `find . -iname \*.zip`
do
	gsutil cp $zip gs://$PROJECT_ID-$BUCKET_NAME
done

cd schema-registry

# Put the REST server application inside a docker container and deploy the image to the Container Registry
echo "Deploying Schema Registry component image in the Container Registry.."
gcloud builds submit --tag gcr.io/$PROJECT_ID/schema-registry

# Deploy the created image to Cloud Run
echo "Deploying the created image to Cloud Run.."
gcloud run deploy schema-registry --image gcr.io/$PROJECT_ID/schema-registry --platform managed --allow-unauthenticated --region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,BUCKET_NAME=$PROJECT_ID-$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE,SERVICE_ACCOUNT_KEY_FILE=key.json

cd ..

# Deploy the Central Consumer to Cloud Functions
echo "Deploying the Central Consumer component.."
gcloud functions deploy central-consumer --runtime go113 --timeout=540s --allow-unauthenticated --entry-point CentralConsumerHandler --source=gs://$PROJECT_ID-$BUCKET_NAME/central-consumer.zip --trigger-topic $INPUT_TOPIC --region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,BUCKET_NAME=$PROJECT_ID-$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE

# Deploy the Puller & Cleaner (JSON and CSV) to Cloud Functions
echo "Deploying the Puller & Cleaner component.."
gcloud functions deploy puller-cleaner-json --runtime go113 --timeout=540s --allow-unauthenticated --entry-point PullerCleaner  --source=gs://$PROJECT_ID-$BUCKET_NAME/puller-cleaner-json.zip --trigger-http --region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,BUCKET_NAME=$PROJECT_ID-$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE

gcloud functions deploy puller-cleaner-csv --runtime go113 --timeout=540s --allow-unauthenticated --entry-point PullerCleaner  --source=gs://$PROJECT_ID-$BUCKET_NAME/puller-cleaner-csv.zip --trigger-http --region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,BUCKET_NAME=$PROJECT_ID-$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE

# Deploy the function that creates Cloud Tasks to invoke Cloud Functions
echo "Deploying the Puller Tasks functions.."

gcloud functions deploy puller-tasks-json --runtime go113 --timeout=540s --allow-unauthenticated --entry-point HTTPPullerTasks --source=gs://$PROJECT_ID-$BUCKET_NAME/puller-tasks-json.zip --trigger-http --region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,REGION=$REGION,BUCKET_NAME=$PROJECT_ID-$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE

gcloud functions deploy puller-tasks-csv --runtime go113 --timeout=540s --allow-unauthenticated --entry-point HTTPPullerTasks --source=gs://$PROJECT_ID-$BUCKET_NAME/puller-tasks-csv.zip --trigger-http --region $REGION --set-env-vars PROJECT_ID=$PROJECT_ID,REGION=$REGION,BUCKET_NAME=$PROJECT_ID-$BUCKET_NAME,CONFIG_FILE=$CONFIG_FILE

# Deploy the helper functions for XML and CSV validation
echo "Deploying the helper Cloud Functions.."

gcloud functions deploy xml-validator --runtime python37 --timeout=540s --allow-unauthenticated --entry-point http_validation_handler --source=gs://$PROJECT_ID-$BUCKET_NAME/xml-validator.zip --trigger-http --region $REGION

gcloud functions deploy csv-validator --runtime java11 --timeout=540 --allow-unauthenticated --entry-point hr.syntio.handler.HttpHandler --source=gs://$PROJECT_ID-$BUCKET_NAME/csv-validator.zip --trigger-http --region $REGION

# Create Cloud Tasks queue
echo "Creating Cloud Tasks Queue.."
gcloud tasks queues create schema-registry-puller-queue-t

# Create Cloud Scheduler jobs
echo "Creating Cloud Scheduler jobs.."
export CF_NAME_JSON=puller-tasks-json
export CF_NAME_CSV=puller-tasks-csv
export SCHEDULER_URI_JSON=https://$REGION-$PROJECT_ID.cloudfunctions.net/$CF_NAME_JSON
export SCHEDULER_URI_CSV=https://$REGION-$PROJECT_ID.cloudfunctions.net/$CF_NAME_CSV

gcloud scheduler jobs create http puller-cleaner-invoker-json --schedule "* * * * *" \
	--uri $SCHEDULER_URI_JSON --message-body "{\"puller-number\":10}"

gcloud scheduler jobs create http puller-cleaner-invoker-csv --schedule "* * * * *" \
	--uri $SCHEDULER_URI_CSV --message-body "{\"puller-number\":10}"

echo "Deployment complete :)"
