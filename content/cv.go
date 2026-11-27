package content

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultCVPath = "docs/cv/cv.yaml"

type CV struct {
	Name        string
	Experiences []Experience
}

type Experience struct {
	Name        string
	Start       string
	End         string
	Description string
	ShowOnSite  *bool `yaml:"show_on_site"`
	Milestones  []Milestone
}

type Milestone struct {
	Label       string
	Title       string
	Period      string
	Description string
	Tags        []string
}

func ReadCV(path string) (CV, error) {
	var cv CV
	file, err := os.Open(path)
	if err != nil {
		return cv, err
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&cv); err != nil {
		return cv, fmt.Errorf("%s: %w", path, err)
	}
	if cv.Name == "" {
		return cv, fmt.Errorf("cv.yaml: missing name")
	}
	var visible []Experience
	for i, e := range cv.Experiences {
		if e.ShowOnSite != nil && !*e.ShowOnSite {
			continue
		}
		if e.Name == "" || e.Start == "" || e.Description == "" {
			return cv, fmt.Errorf("cv.yaml: experience %d needs a name, start and description", i+1)
		}
		visible = append(visible, e)
	}
	cv.Experiences = visible
	return cv, nil
}
