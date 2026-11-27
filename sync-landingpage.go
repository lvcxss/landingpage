package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var postStrings = map[int]string{
	1: "experience",
	2: "write-up",
	3: "project",
	4: "guide",
}

type SyncPayload struct {
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
	PostType string `json:"post_type"`
	Content  string `json:"content"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	targetURL := "http://localhost:4000/api/admin/sync"
	secretToken := os.Getenv("ADMIN_TOKEN")

	if secretToken == "" {
		panic("ADMIN_TOKEN not defined (check .env)")
	}

	var title, filePath, totpCode string
	var typeInt int

	fmt.Print("Title: ")
	_, err = fmt.Scan(&title)
	if err != nil {
		panic("error reading title")
	}

	fmt.Print("Filepath: ")
	_, err = fmt.Scan(&filePath)
	if err != nil {
		panic("error reading filepath")
	}

	fmt.Print("Post Type ( 1 - experience | 2 - write-up | 3 - project | 4 - guide ) : ")
	_, err = fmt.Scan(&typeInt)
	if err != nil {
		panic("error reading typeInt")
	}

	fmt.Print("TOTP Code : ")
	_, err = fmt.Scan(&totpCode)
	if err != nil {
		panic("error reading TOTP")
	}

	mdContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("error reading file:", err)
	}

	payload := SyncPayload{
		Title:    title,
		FilePath: filePath,
		PostType: postStrings[typeInt],
		Content:  string(mdContent),
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Erro a criar pedido:", err)
	}

	req.Header.Set("X-Admin-Token", secretToken)
	req.Header.Set("X-TOTP-Code", totpCode)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Erro de rede ao contactar servidor:", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("sync failed! status: %d - reason: %s", resp.StatusCode, string(body))
	}

	fmt.Println("post synchronized successfuly")
}
