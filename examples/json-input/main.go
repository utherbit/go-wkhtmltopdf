package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	pdf "github.com/utherbit/go-wkhtmltopdf"
)

// For the full list of options, see pkg.go.dev/github.com/utherbit/go-wkhtmltopdf.
// NOTE: pdf.ConverterOpts and pdf.ObjectOpts also support YAML unmarshalling.
var jsonInput = strings.NewReader(`{
	"converterOpts": {
		"title": "google.com",
		"paperSize": "A4",
		"orientation": "Portrait",
		"marginLeft": "10mm",
		"marginRight": "10mm"
	},
	"objectOpts": {
		"location": "https://google.com",
		"footer": {
			"contentCenter": "[page]",
			"fontSize": 14
		}
	}
}`)

type inputData struct {
	ConverterOpts *pdf.ConverterOpts `json:"converterOpts"`
	ObjectOpts    *pdf.ObjectOpts    `json:"objectOpts"`
}

func main() {
	// Initialize library.
	if err := pdf.Init(); err != nil {
		log.Fatal(err)
	}
	defer pdf.Destroy()

	// Set default options. Any option fields specified in the JSON
	// input data will overwrite the defaults.
	input := &inputData{
		ConverterOpts: pdf.NewConverterOpts(),
		ObjectOpts:    pdf.NewObjectOpts(),
	}

	if err := json.NewDecoder(jsonInput).Decode(input); err != nil {
		log.Fatal(err)
	}

	// Create object.
	object, err := pdf.NewObjectWithOpts(input.ObjectOpts)
	if err != nil {
		log.Fatal(err)
	}

	// Create converter.
	converter, err := pdf.NewConverterWithOpts(input.ConverterOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer converter.Destroy()

	// Add object to the converter.
	converter.Add(object)

	// Convert object and save the output PDF document.
	outFile, err := os.Create("out.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	if err := converter.Run(outFile); err != nil {
		log.Fatal(err)
	}
}
