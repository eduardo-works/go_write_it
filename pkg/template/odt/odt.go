package odt

import (
	_ "embed"
	"os"

	"github.com/eduardo-works/writeit/pkg/project"
)

//go:embed _standalone.odt
var StandaloneODT []byte

//go:embed _chaptered.odt
var ChapteredDT []byte

//go:embed _parted.odt
var PartedODT []byte

//go:embed _other.odt
var OtherODT []byte

//go:embed _clean.odt
var CleanODT []byte

const Clean string = "clean"
const Other string = "other"

func PrepareTemplate(templateType any) (string, error) {
	switch v := templateType.(type) {
	case project.Project:
		if v.IsStandalone() {
			return referenceFile(StandaloneODT)
		}
		if v.IsChaptered() {
			return referenceFile(ChapteredDT)
		}
		return referenceFile(PartedODT)

	case string:
		if v == "clean" {
			return referenceFile(CleanODT)
		}
	}

	return referenceFile(OtherODT)
}

func referenceFile(reference []byte) (string, error) {
	file, err := os.CreateTemp("", "writeit-template.odt")
	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err := file.Write(reference); err != nil {
		os.Remove(file.Name())
		return "", err
	}

	return file.Name(), nil
}
