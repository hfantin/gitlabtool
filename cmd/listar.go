package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// 	Black   = Color("\033[1;30m%s\033[0m")
// 	Red     = Color("\033[1;31m%s\033[0m")
// 	Green   = Color("\033[1;32m%s\033[0m")
// 	Yellow  = Color("\033[1;33m%s\033[0m")
// 	Purple  = Color("\033[1;34m%s\033[0m")
// 	Magenta = Color("\033[1;35m%s\033[0m")
// 	Teal    = Color("\033[1;36m%s\033[0m")
// 	White   = Color("\033[1;37m%s\033[0m")

type MergeRequest struct {
	State          string `json:"state"`
	Title          string `json:"title"`
	UserNotesCount int    `json:"user_notes_count"`
	WebUrl         string `json:"web_url"`
}

func ListarMergeRequest(urlApi, token, state string, minComents int) {
	// endpoint := fmt.Sprintf("%s%s", urlApi, "/projects/mov%2fmov-react-native/merge_requests?state=opened")
	stateQuery := ""
	if len(state) > 0 {
		stateQuery = fmt.Sprintf("?state=%s", state)
	}
	endpoint := fmt.Sprintf("%s/%s%s", urlApi, "merge_requests", stateQuery)
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
		found := false
		for _, mr := range data {
			match := r.FindStringSubmatch(mr.WebUrl)
			projeto := ""
			if len(match) == 3 {
				projeto = match[2]
			}
			if mr.UserNotesCount > minComents {
				fmt.Print(fmt.Sprintf("\033[1;32m[%s]\033[0m ", projeto))
				if state != "opened" {
					fmt.Print(fmt.Sprintf("\033[1;33m%s\033[0m ", mr.State))
				}
				fmt.Print(fmt.Sprintf("MR: %s ", mr.Title))
				fmt.Println(fmt.Sprintf("\033[1;36m[%s]\033[0m", mr.WebUrl))
				found = true
			}
		}
		if !found {
			fmt.Println("Nenhum Merge Request encontrado.")
		}

	} else {
		fmt.Printf("Não foi possível obter a lista de MRs: %s\n", string(body))
	}
}
