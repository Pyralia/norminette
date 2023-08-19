package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"regexp"
)

var issues []string

func main() {
	// Vérifier si le chemin du dossier est fourni en argument
	if len(os.Args) != 2 {
		fmt.Println("Utilisation :", os.Args[0], "<chemin-du-dossier>")
		return
	}

	// Récupérer le chemin du dossier à partir des arguments
	dirPath := os.Args[1]

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Vérifier si le fichier est un fichier .java
		if strings.HasSuffix(info.Name(), ".java") {
			if err := analyzeJavaFile(path); err != nil {
				fmt.Printf("Erreur lors de l'analyse du fichier %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Erreur lors de la lecture du dossier.")
	}

	printIssues()
}

func analyzeJavaFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Vérifier l'espace entre "public void" et le nom de la méthode
		if matched, _ := regexp.MatchString(`\bvoid\s+[a-z]\w*\s*\(.*`, line); matched {
			parts := strings.SplitN(line, "void", 2)
			methodName := strings.TrimSpace(parts[1])
			if methodName == "" {
				continue
			}


			if !strings.Contains(line, "public") && !strings.Contains(line, "private"){
				issues = append(issues, fmt.Sprintf("Vous n'avez pas précisé le type de la méthode '%s': %s", filePath, line))
			}
			if !strings.Contains(line, ")") {
				issues = append(issues, fmt.Sprintf("La méthod n'est pas proprement définie '%s': %s", filePath, line))
			}
			if !strings.Contains(line, "{") {
				issues = append(issues, fmt.Sprintf("La méthode ne se termine pas par un bracket '%s': %s", filePath, line))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}


func printIssues() {
	if len(issues) > 0 {
		fmt.Println("Problèmes détectés :")
		for _, issue := range issues {
			fmt.Println(issue)
		}
	} else {
		fmt.Println("Aucun problème détecté.")
	}
}