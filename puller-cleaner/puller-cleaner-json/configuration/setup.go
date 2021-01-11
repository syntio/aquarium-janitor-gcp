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

// Package configuration handles extracting configuration parameters from a config file and environment variables.
package configuration

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"

	"gopkg.in/yaml.v2"
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

	Functions struct {
		SchemaRegistryURL          string `yaml:"schemaRegistryURL"`
		CsvValidatorURL            string `yaml:"csvValidatorURL"`
		XmlValidatorURL            string `yaml:"xmlValidatorURL"`
		SchemaRegistryEvolutionURL string `yaml:"schemaRegistryEvolutionURL"`
		PullerCleanerJsonURL       string `yaml:"pullerCleanerJsonURL"`
		PullerCleanerCsvURL        string `yaml:"pullerCleanerCsvURL"`
	} `yaml:"functions"`

	Protoparam struct {
		TmpFilePath string `yaml:"tmpFilePath"`
		TmpFileName string `yaml:"tmpFileName"`
	} `yaml:"protoparam"`

	PullerCleanerJSON struct {
		TimeDurationSeconds time.Duration `yaml:"timeDurationSeconds"`
		MaxBatchSize        int           `yaml:"maxBatchSize"`
		MaxThroughput       int           `yaml:"maxThroughput"`
	} `yaml:"pullercleanerjson"`

	PullerCleanerCSV struct {
		TimeDurationSeconds time.Duration `yaml:"timeDurationSeconds"`
		MaxBatchSize        int           `yaml:"maxBatchSize"`
		MaxThroughput       int           `yaml:"maxThroughput"`
	} `yaml:"pullercleanercsv"`

	PullerCleanerTaskQueue  string      `yaml:"pullerCleanerTaskQueue"`
	ContentType             string      `yaml:"contentType"`
	FileMode                os.FileMode `yaml:"fileMode"`
	FirestoreCollectionName string      `yaml:"firestoreCollectionName"`
}

// RetrieveConfig obtains configuration parameters from a
// config file (which is in GCS bucket) into an object.
//
// Output parameters are a Config struct which will be filled
// with configuration parameters, and a possible error occurred.
func RetrieveConfig() (cfg Config, err error) {
	bucketName := os.Getenv("BUCKET_NAME")
	fileName := os.Getenv("CONFIG_FILE")
	if bucketName == "" || fileName == "" {
		return cfg, fmt.Errorf("ERROR: Couldn't read environment variables.")
	}

	fileContent, err := readFromBucket(bucketName, fileName)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(fileContent, &cfg); err != nil {
		return cfg, fmt.Errorf("ERROR: Configuration object can't be obtained from storage. %v", err)
	}

	return cfg, nil
}

// readFromBucket reads file from a specified storage bucket.
//
// Input parameters are a bucket and a file to be read.
//
// Output parameters are file's content as byte array, and a possible error
// occurred while connection to GCS bucket.
func readFromBucket(bucketName string, fileName string) ([]byte, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Couldn't connect to GCS. %v", err)
	}

	rc, err := client.Bucket(bucketName).Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Reader from cloud storage can't be obtained. Check environment variables. %v", err)
	}

	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil, fmt.Errorf("ERROR: Storage object is not valid. %v", err)
	}

	return slurp, nil
}
