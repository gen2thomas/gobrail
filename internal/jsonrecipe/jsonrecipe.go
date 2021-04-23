package jsonrecipe

import (
	"fmt"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

// PrepareAndValidate make file an absolute path and validate the json style
func PrepareAndValidate(schema string, document string) (documentAbs string, err error) {
	documentAbs, err = filepath.Abs(document)
	if err != nil {
		return
	}
	var schemaAbs string
	schemaAbs, err = filepath.Abs(schema)
	if err != nil {
		return
	}
	err = validate("file://"+schemaAbs, "file://"+documentAbs)
	return
}

func validate(schema string, document string) error {

	schemaLoader := gojsonschema.NewReferenceLoader(schema)
	documentLoader := gojsonschema.NewReferenceLoader(document)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}

	for _, desc := range result.Errors() {
		fmt.Printf("- %s\n", desc)
	}
	return fmt.Errorf("The document is not valid. see errors above")
}
