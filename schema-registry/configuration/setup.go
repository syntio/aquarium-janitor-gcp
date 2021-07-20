// Copyright 2020 Syntio Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	/*Functions struct {
		SchemaRegistryURL          string `yaml:"schemaRegistryURL"`
		CsvValidatorURL            string `yaml:"csvValidatorURL"`
		XmlValidatorURL            string `yaml:"xmlValidatorURL"`
		SchemaRegistryEvolutionURL string `yaml:"schemaRegistryEvolutionURL"`
		PullerCleanerJsonURL       string `yaml:"pullerCleanerJsonURL"`
		PullerCleanerCsvURL        string `yaml:"pullerCleanerCsvURL"`
	} `yaml:"functions"`*/

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

	ContentType             string      `yaml:"contentType"`
	FileMode                os.FileMode `yaml:"fileMode"`
	FirestoreCollectionName string      `yaml:"firestoreCollectionName"`
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
