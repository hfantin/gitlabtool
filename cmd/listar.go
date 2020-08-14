package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const regex = `^http[s]?:\/\/[-a-zA-Z0-9@:%._\+~#=]{1,256}\/([-a-zA-Z0-9@:%._\+~#=]{1,32})\/([a-zA-Z-]*)`

var projectRegex = regexp.MustCompile(regex)

type MergeRequest struct {
	State          string `json:"state"`
	Title          string `json:"title"`
	UserNotesCount int    `json:"user_notes_count"`
	WebUrl         string `json:"web_url"`
}

func ListarMergeRequest(urlApi, token, state string, minComents int) {
	var endpoint = obterEndpoint(state, urlApi)

	resp, err := executarGet(endpoint, token)
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
		validarResposta(body, minComents)
	} else {
		fmt.Printf("Não foi possível obter a lista de MRs: %s\n", string(body))
	}
}

func executarGet(endpoint, token string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("PRIVATE-TOKEN", token)
	return client.Do(req)
}

func obterEndpoint(state string, urlApi string) string {
	stateQuery := ""
	if len(state) > 0 {
		stateQuery = fmt.Sprintf("?state=%s", state)
	}
	return fmt.Sprintf("%s/%s%s", urlApi, "merge_requests", stateQuery)
}

func validarResposta(body []byte, minComents int) {
	var mergeRequests []MergeRequest
	err := json.Unmarshal(body, &mergeRequests)
	if err != nil {
		fmt.Printf("Falha ao decodificar a resposta: %v\n", err)
		return
	}
	found := false
	for _, mr := range mergeRequests {
		if mr.UserNotesCount > minComents {
			exibirResultado(&mr)
			found = true
		}
	}
	if !found {
		fmt.Println("Nenhum Merge Request encontrado.")
	}
}

func exibirResultado(mr *MergeRequest) {
	fmt.Print(fmt.Sprintf("\033[1;32m[%s]\033[0m ", obterNomeProjeto(mr.WebUrl)))
	if mr.State != "opened" {
		fmt.Print(fmt.Sprintf("\033[1;33m%s\033[0m ", mr.State))
	}
	fmt.Print(fmt.Sprintf("MR: %s ", mr.Title))
	fmt.Println(fmt.Sprintf("\033[1;36m[%s]\033[0m", mr.WebUrl))
}

func obterNomeProjeto(url string) string {
	match := projectRegex.FindStringSubmatch(url)
	if len(match) == 3 {
		return match[2]
	}
	return ""
}
