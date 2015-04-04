package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type Project struct {
	Name      string
	Author    string
	Directory *Directory
	Year      int
	License   string
	Files     []*File
}
type Directory struct {
	Files       []*File
	Directories []*Directory
	Name        string
	Parent      *Directory
}

func (d *Directory) AddFile(filename, fileContent string) *File {
	file := new(File)
	file.Parent = d
	file.Name = filename
	file.Content = fileContent
	d.Files = append(d.Files, file)
	return file
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
	project.Directory.AddFile("LICENSE", "")
	project.Directory.AddFile("README.md", fmt.Sprintf(readme, project.Name, project.Author, project.Year, project.License))
	project.Directory.AddFile(".gitignore", "cover.out")
	project.Directory.AddFile(".travis.yml", travis)
	project.Directory.AddFile(project.Name+".go", fmt.Sprintf(mainFileContent, project.Year, project.Author, project.License, project.Name))
	project.Directory.AddFile(project.Name+"_test.go", fmt.Sprintf(testFileContent, project.Name))
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
	err = InitGitRepository(project.Directory.Name)
	exitOnError(err)
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

// InitGitRepository create a git repository in the directory
func InitGitRepository(directory string) error {
	var (
		dir string
		err error
		cmd *exec.Cmd
	)
	dir, err = os.Getwd()
	exitOnError(err)
	err = os.Chdir(directory)
	exitOnError(err)
	cmd = exec.Command("git", "init", "-q")
	cmd.Stdout = os.Stdin
	defer os.Chdir(dir)
	err = cmd.Run()
	return err
}
