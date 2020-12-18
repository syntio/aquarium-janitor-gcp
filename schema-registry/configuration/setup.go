package configuration

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"

	"gopkg.in/yaml.v3"
)

// Config structure is used to capture configuration parameters.
// Config structure is used to capture configuration parameters.
type Config struct {
	Topics struct {
		InvalidTopicJSON string `yaml:"invalidTopicJSON"`
		InvalidTopicCSV  string `yaml:"invalidTopicCSV"`
		ValidTopic       string `yaml:"validTopic"`
		DeadLetterTopic  string `yaml:"deadLetterTopic"`
		InputTopic       string `yaml:"inputTopic"`
	} `yaml:"topics"`

	Subscriptions struct {
		SubIdJSON       string `yaml:"subIdJSON"`
		SubIdCSV        string `yaml:"subIdCSV"`
		SubIdInput      string `yaml:"subIdInput"`
		SubIdValid      string `yaml:"subIdValid"`
		SubIdDeadletter string `yaml:"subIdDeadletter"`
	} `yaml:"subscriptions"`

	Functions struct {
		SchemaRegistryURL          string `yaml:"schemaRegistryURL"`
		CsvValidatorURL            string `yaml:"csvValidatorURL"`
		XmlValidatorURL            string `yaml:"xmlValidatorURL"`
		SchemaRegistryEvolutionURL string `yaml:"schemaRegistryEvolutionURL"`
	} `yaml:"functions"`

	Protoparam struct {
		TmpFilePath string `yaml:"tmpFilePath"`
		TmpFileName string `yaml:"tmpFileName"`
	} `yaml:"protoparam"`

	PullerCleanerJSON struct {
		TimeDurationSecondsAsync time.Duration `yaml:"timeDurationSecondsAsync"`
		MaxBatchSize             int           `yaml:"maxBatchSize"`
		MaxThroughput            int           `yaml:"maxThroughput"`
		LoggingEnabled           bool          `yaml:"loggingEnabled"`
	} `yaml:"pullercleanerjson"`

	PullerCleanerCSV struct {
		TimeDurationSecondsAsync time.Duration `yaml:"timeDurationSecondsAsync"`
		MaxBatchSize             int           `yaml:"maxBatchSize"`
		MaxThroughput            int           `yaml:"maxThroughput"`
		LoggingEnabled           bool          `yaml:"loggingEnabled"`
	} `yaml:"pullercleanercsv"`

	ContentType            string      `yaml:"contentType"`
	FileMode               os.FileMode `yaml:"fileMode"`
	FirebaseCollectionName string      `yaml:"firebaseCollectionName"`
}

// Function for obtaining configuration parameters values into an object.
func RetrieveConfig() (cfg Config) {
	bucketName := os.Getenv("BUCKET_NAME")
	fileName := os.Getenv("CONFIG_FILE")

	fileContent := ReadFromBucket(bucketName, fileName)

	if err := yaml.Unmarshal(fileContent, &cfg); err != nil {
		log.Printf("ERROR: Configuration object can't be obtained from storage. %v.\n", err)
	}
	return cfg
}

func ReadFromBucket(bucketName string, fileName string) []byte {
	ctx := context.Background()

	//using cloud function
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Println(err)
		return nil
	}

	rc, err := client.Bucket(bucketName).Object(fileName).NewReader(ctx)
	if err != nil {
		log.Printf("ERROR: Reader from cloud storage can't be obtained. Check environment variables. %v.\n", err)
		return nil
	}
	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		log.Printf("ERROR: Storage object is not valid. %v.\n", err)
		return nil
	}
	return slurp
}
