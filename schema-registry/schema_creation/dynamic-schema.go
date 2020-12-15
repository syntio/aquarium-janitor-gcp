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

package schema_creation

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
)

const (
	pythonJSONSchemaScriptName = "jsonSchemaDynamicCreation.py"
)

// JSONSchemaDynamicCreation takes serialized data of a message and converts it into a JSON-Schema.
// find out more about json standardization here http://json-schema.org/.
// The output of the Dynamic creation is a serialized schema representation, a flag indicating if the operation
// was a success and an error.
func JSONSchemaDynamicCreation(data []byte) ([]byte, bool, error) {
	script := exec.Command("python", pythonJSONSchemaScriptName)
	stdin, err := script.StdinPipe()
	if err != nil {
		return nil, false, fmt.Errorf("Script stdin error: %v", err)
	}

	go func() {
		defer stdin.Close()
		_, err = stdin.Write(data)
	}()
	if err != nil {
		return nil, false, fmt.Errorf("Writing to script stdin error: %v", err)
	}

	scriptResult, err := script.CombinedOutput()
	if err != nil {
		return nil, false, fmt.Errorf("Script running and getting result error: %v", err)
	}
	if len(scriptResult) == 0 {
		return scriptResult, false, nil
	}

	return scriptResult, true, nil
}

// CSVSchemaDynamicCreation takes serialized data of a csv file and converts it into a csv schema.
// The output of the Dynamic creation is a serialized schema representation, a flag indicating if the operation
// was a success and an error.
func CSVSchemaDynamicCreation(data []byte) ([]byte, bool, error) {
	del := []rune(",")
	r := bytes.NewReader(data)

	reader := csv.NewReader(r)
	reader.Comma = rune(del[0])

	lines, err := reader.ReadAll()
	if err != nil {
		return []byte{}, false, err
	}
	var b bytes.Buffer
	b.Write([]byte("version 1.1 \n"))
	b.Write([]byte("@totalColumns " + strconv.Itoa(len(lines[0])) + "\n"))

	for key := range lines[0] {
		b.Write([]byte(lines[0][key] + ": \n"))
	}

	return b.Bytes(), true, nil

}
