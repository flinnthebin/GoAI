package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func showHelp() {
	fmt.Println(`Usage: iago -s <tokens> -t <temperature> <input_file>

Options:
    -s <tokens>         Set number of tokens to generate
    -t <temperature>    Temperature for response generation
    -h                  Show help menu

Example Usage: iago -s 4096 -t 0.1 input.txt

------------

Tokens

------------

A token is considered to be 4 characters. The number of tokens set must be enough to include the length
of the provided input and expected output.

------------

Temperature

------------

The temperature parameter controls the randomness of the output. 

Low (0.1 - 0.3)
Predictable, deterministic output.

Medium (0.4 - 0.6)
An adaptable range balancing creativity and predictability

High (0.7 - 1.0)
Creative and unpredictable. Strong possibility of hallucination.`)
}

func generateOutFileName(inFile string) string {
	baseName := filepath.Base(inFile)
	extension := filepath.Ext(baseName)
	return fmt.Sprintf("output%s", extension)
}

func writeToFile(outFile, content string) error {
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type ChatCompletion struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func main() {

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	var (
		TOKEN int
		TEMP  float64
	)

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-s":
			if i+1 < len(os.Args) {
				TOKEN, _ = strconv.Atoi(os.Args[i+1])
			}
		case "-t":
			if i+1 < len(os.Args) {
				TEMP, _ = strconv.ParseFloat(os.Args[i+1], 64)
			}
		case "-h":
			showHelp()
			return
		}
	}

	if TOKEN == 0 || TEMP == 0 {
		fmt.Println("Error: Both tokens and temperature must be specified.")
		showHelp()
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Error: No argument for input file")
		showHelp()
		return
	}

	inFile := os.Args[len(os.Args)-1]
	outFile := generateOutFileName(inFile)

	file, err := os.Open(inFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	var messages []ChatMessage
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		message := ChatMessage{Role: "user", Content: line}
		messages = append(messages, message)
	}
	jsonData, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	cmd := exec.Command("curl",
		"-X", "POST",
		"https://api.openai.com/v1/chat/completions",
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("Authorization: Bearer %s", os.Getenv("OPENAI_API")),
		"--data", fmt.Sprintf(`{"model": "gpt-4-turbo-preview", "messages": %s, "max_tokens": %d, "temperature": %f}`, string(jsonData), TOKEN, TEMP),
	)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var completion ChatCompletion
	if err := json.Unmarshal(out, &completion); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	if len(completion.Choices) == 0 {
		fmt.Println("Error: No return response")
		return
	}

	if err := writeToFile(outFile, completion.Choices[0].Message.Content); err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}
	fmt.Println("I am not what I am")
}
