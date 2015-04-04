package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type Project struct {
	Name        string
	Author      string
	Directory   *Directory
	Year        int
	License     string
	Files       []*File
	Directories []*Directory
}
type Directory struct {
	Files       []*File
	Directories []*File
	Name        string
	Parent      *Directory
}
type File struct {
	Name    string
	Content string
	Parent  *Directory
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	var (
		projectName    string
		projectAuthor  string
		projectLicense string
		project        *Project
		err            error
	)

	flag.StringVar(&projectAuthor, "author", os.Getenv("USERNAME"), "project's author")
	flag.StringVar(&projectLicense, "license", "MIT", "project's license")
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) <= 0 || len(arguments[0]) < 2 {
		exitOnError(errors.New("missing project name as argument"))
	} else {
		projectName = strings.Replace(arguments[0], " ", "", -1)
	}
	rootDirectory, err := os.Getwd()
	exitOnError(err)
	project = new(Project)
	project.Author = projectAuthor
	project.Directory = new(Directory)
	project.License = projectLicense
	project.Name = projectName
	project.Directory.Name = path.Join(rootDirectory, project.Name)
	project.Year = time.Now().Year()
	readMeFile := new(File)
	readMeFile.Name = "README.md"
	readMeFile.Content = fmt.Sprintf(readme, project.Name, project.Author, project.Year, project.License)
	readMeFile.Parent = project.Directory
	gitignoreFile := new(File)
	gitignoreFile.Name = ".gitignore"
	gitignoreFile.Content = "cover.out"
	travisFile := new(File)
	travisFile.Name = ".travis.yml"
	travisFile.Content = travis
	mainFile := new(File)
	mainFile.Name = project.Name + ".go"
	mainFile.Content = fmt.Sprintf(mainFileContent, project.Year, project.Author, project.License, project.Name)
	testFile := new(File)
	testFile.Name = project.Name + "_test.go"
	testFile.Content = fmt.Sprintf(testFileContent, project.Name)
	project.Directory.Files = append(project.Directory.Files, readMeFile, gitignoreFile, travisFile, mainFile, testFile)
	fmt.Printf("Creating project %s \n", project.Name)
	fmt.Printf("Creating directory %s \n", project.Directory.Name)
	err = os.MkdirAll(project.Directory.Name, 0644)
	exitOnError(err)
	for _, file := range project.Directory.Files {
		fmt.Printf("Writing file")
		path := path.Join(project.Directory.Name, file.Name)
		fmt.Printf("Writing file %s\n", path)
		err = ioutil.WriteFile(path, []byte(file.Content), 0644)
	}

}

const readme string = `#%s
	
Author:  %s

Year: %d

License: %v
`

const travis string = "language: go"
const mainFileContent string = `// Copyright %d %s
// License %s

package %s
`
const testFileContent string = `package %s

import(
	"testing"
)

func Test(t *testing.T){
	t.Log("Test")
}
`
