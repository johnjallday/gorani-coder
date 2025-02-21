package docbuilder

import (
	"log"
	"os"
)

// BuildReadme gathers docs/installation.md, docs/requirements.md,
// and the logo image from docs/logo.png and combines them into README.md.
func BuildReadme() {
	logoMarkdown := `<img src="docs/logo.png" alt="Project Logo" width="300">` + "\n\n"

	installContent, err := os.ReadFile("docs/installation.md")
	if err != nil {
		log.Fatalf("failed reading docs/installation.md: %v", err)
	}

	reqContent, err := os.ReadFile("docs/requirements.md")
	if err != nil {
		log.Fatalf("failed reading docs/requirements.md: %v", err)
	}

	roadmapContent, err := os.ReadFile("docs/roadmap.md")
	if err != nil {
		log.Fatalf("failed reading docs/roadmap.md: %v", err)
	}

	/*
		logoData, err := ioutil.ReadFile("docs/logo.png")
		if err != nil {
			log.Fatalf("failed reading docs/logo.png: %v", err)
		}
		encodedLogo := base64.StdEncoding.EncodeToString(logoData)
		logoMarkdown = "![Project Logo](data:image/png;base64," + encodedLogo + ")\n\n"
	*/

	// Combine the contents into a single README.md content.
	combined := logoMarkdown + string(installContent) + "\n\n" + string(reqContent) + "\n\n" + string(roadmapContent)

	// Write the combined content to README.md in the project root.
	err = os.WriteFile("README.md", []byte(combined), 0644)
	if err != nil {
		log.Fatalf("failed writing README.md: %v", err)
	}

	log.Println("README.md generated successfully!")
}
