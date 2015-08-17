package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"
)

var (
	textdir  = flag.String("text", "text", "directory containing text")
	tmplfile = flag.String("tmpl", "base.tmpl", "template file name")
	outfile  = flag.String("out", "resume.html", "file to generate the resume to")
)

func main() {
	flag.Parse()

	resume, err := loadFiles(*textdir)
	if err != nil {
		return
	}

	tmpl, err := template.New("resume").ParseFiles(*tmplfile)
	if err != nil {
		log.Fatalf("could not parse template file: %s", err)
	}

	out, err := os.Create(*outfile)
	if err != nil {
		log.Fatalf("could not open out file: %s", err)
	}

	err = tmpl.ExecuteTemplate(out, *tmplfile, resume)
	if err != nil {
		log.Fatalf("could not create resume: %s", err)
	}
}

func loadFiles(dirname string) (*Resume, error) {
	personalF, err := os.Open(filepath.Join(dirname, "personal.yaml"))
	if err != nil {
		log.Printf("error opening personal file: %s", err)
		return nil, err
	}
	defer personalF.Close()
	personal, err := ioutil.ReadAll(personalF)
	if err != nil {
		log.Printf("error reading personal file: %s", err)
		return nil, err
	}

	educationF, err := os.Open(filepath.Join(dirname, "education.yaml"))
	if err != nil {
		log.Printf("error opening education file: %s", err)
		return nil, err
	}
	defer educationF.Close()
	education, err := ioutil.ReadAll(educationF)
	if err != nil {
		log.Printf("error reading education file: %s", err)
		return nil, err
	}

	workF, err := os.Open(filepath.Join(dirname, "work.yaml"))
	if err != nil {
		log.Printf("error opening work file: %s", err)
		return nil, err
	}
	defer workF.Close()
	work, err := ioutil.ReadAll(workF)
	if err != nil {
		log.Printf("error reading work file: %s", err)
		return nil, err
	}

	resume := Resume{}
	err = yaml.Unmarshal(personal, &resume.Personal)
	if err != nil {
		log.Printf("error parsing personal.yaml: %s", err)
		return nil, err
	}
	yaml.Unmarshal(work, &resume.WorkExperience)
	if err != nil {
		log.Printf("error parsing work.yaml: %s", err)
		return nil, err
	}
	yaml.Unmarshal(education, &resume.Education)
	if err != nil {
		log.Printf("error parsing education.yaml: %s", err)
		return nil, err
	}
	return &resume, nil
}

// Resume defines your entire resume
type Resume struct {
	Personal       Personal
	Education      []Education
	WorkExperience []Work
}

// Personal defines all your personal information
type Personal struct {
	Name      string
	Website   string
	Email     string
	Telephone string
	About     string
}

// Education defines an individual course you took
type Education struct {
	University string
	College    string
	From       int
	To         int
	Degree     string
	Subject    string
	Course     string
}

// Work defines an individual work experience
type Work struct {
	From     int
	To       int
	Position string
	Company  string
	Place    string
	KeyWords []string
	Summary  string
	Detailed string
}
