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
