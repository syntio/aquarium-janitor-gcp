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
//
// Package task implements the Cloud Task creation
// Each task invokes the Puller & Cleaner Cloud Function that handles the schema recovery and evolution
package task

import (
	"context"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

// The main function that handles the Cloud Task creation
// Input parameters are the configuration values projectId and region,
// along with the URLs to the Cloud Task queue and the Puller Cleaner Cloud function
// neccessary for the task creation
func CreateHTTPTargetTask(projectID string, region string, pullerTaskQueue string, pullerCleanerURL string) error {
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}

	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, region, pullerTaskQueue)

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        pullerCleanerURL,
				},
			},
		},
	}

	_, err = client.CreateTask(ctx, req)
	if err != nil {
		return fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return nil
}
