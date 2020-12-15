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

package model

import (
	"time"
)

// Schema is a structure that defines how schemas will be saved in our Schema Registry.
type Schema struct {
	Id            string           `json:"id,omitempty" bson:"_id,omitempty" firestore:"ID"`
	SchemaType    string           `json:"schema-type" bson:"schema-type" firestore:"schema-type"`
	Autogenerated bool             `json:"autogenerated" bson:"autogenerated" firestore:"autogenerated"`
	Description   string           `json:"description" bson:"description" firestore:"description"`
	CreationDate  time.Time        `json:"creation-date" bson:"creation-date" firestore:"creation-date"`
	Name          string           `json:"name" bson:"name" firestore:"name"`
	SchemaDetails []*SchemaDetails `json:"schemas" bson:"schemas" firestore:"schemas"`
}

type SchemaDetails struct {
	Version       int32  `json:"version" bson:"version" firestore:"version"`
	Specification string `json:"specification" bson:"specification" firestore:"specification"`
	SchemaHash    string `json:"schema-hash" bson:"schema-hash" firestore:"schema-hash"`
}

// InsertInfo is a return value from a DB used when updating the DB
type InsertInfo struct {
	Id      string `json:"id,omitempty" bson:"_id,omitempty"`
	Version int32  `json:"version" bson:"version"`
}
