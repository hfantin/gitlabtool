package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	Normal = Teal
	Debug  = Green
	Info   = Teal
	Warn   = Yellow
	Fata   = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type MergeRequest struct {
	State          string `json:"state"`
	Title          string `json:"title"`
	UserNotesCount int    `json:"user_notes_count"`
	WebUrl         string `json:"web_url"`
}

func ListarMergeRequest(urlApi, token string) {
	// endpoint := fmt.Sprintf("%s%s", urlApi, "/projects/mov%2fmov-react-native/merge_requests?state=opened")
	endpoint := fmt.Sprintf("%s/%s", urlApi, "merge_requests")
	fmt.Println("endpoint ", endpoint)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("URL inválida: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Falha ao ler resposta: %v\n", err)
		return
	}
	if resp.StatusCode < http.StatusBadRequest {
		var data []MergeRequest
		err := json.Unmarshal(body, &data)
		if err != nil {
			fmt.Printf("Falha ao decodificar a resposta: %v\n", err)
			return
		}
		r := regexp.MustCompile(`^https:\/\/.*.intranet.bb.com.br\/(.{3})\/([a-zA-Z-]*)`)
		for i := 0; i < len(data); i++ {
			match := r.FindStringSubmatch(data[i].WebUrl)
			projeto := ""
			if len(match) == 3 {
				projeto = match[2]
			}
			if data[i].State == "opened" && data[i].UserNotesCount > 0 {
				// const (
				// 	InfoColor    = "\033[1;34m%s\033[0m"
				// 	NoticeColor  = "\033[1;36m%s\033[0m"
				// 	WarningColor = "\033[1;33m%s\033[0m"
				// 	ErrorColor   = "\033[1;31m%s\033[0m"
				// 	DebugColor   = "\033[0;36m%s\033[0m"
				// )

				fmt.Println(Debug("["+projeto+"]"), "MR: "+data[i].Title, data[i].Title, Normal("["+data[i].WebUrl+"]"))
			}
		}

	} else {
		fmt.Printf("falha: %s\n", string(body))
	}
}
