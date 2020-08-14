package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/hfantin/gitlabtool/cmd"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

type Env struct {
	UrlApi string
	Token  string
}

var env *Env

func main() {
	loadEnv()
	exibirMenu()
	opcao := exibirPrompt("Opção")
	executarComando(opcao)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Arquivo .env não encontrado")
	}
	env = &Env{Token: os.Getenv("gitlab_access_token"), UrlApi: os.Getenv("gitlab_url_api")}
	if env.Token == "" {
		fmt.Println("Token não encontrado, favor configurar a variavel gitlab_access_token .env ou no ambiente.")
		os.Exit(1)
	}
	if env.UrlApi == "" {
		fmt.Println("Url da api não encontrado, favor configurar a variavel gitlab_url_api no arquivo .env ou no ambiente.")
		os.Exit(1)
	}
}

func exibirPrompt(label string) int {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validar,
	}
	opcao, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	opcaoInt, err := strconv.Atoi(opcao)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return opcaoInt
}

func exibirMenu() {
	fmt.Println("Selecione a opção:")
	fmt.Println(" 1) Listar MRs abertos com comentários")
	fmt.Println(" 2) Listar MRs abertos")
	fmt.Println(" 3) Exibir todos MRs")
	fmt.Println(" 9) Sair")
	fmt.Println("")
}

func executarComando(opcao int) {
	switch opcao {
	case 1:
		cmd.ListarMergeRequest(env.UrlApi, env.Token, "opened", 0)
	case 2:
		cmd.ListarMergeRequest(env.UrlApi, env.Token, "opened", -1)
	case 3:
		cmd.ListarMergeRequest(env.UrlApi, env.Token, "all", -1)
	case 9:
		fmt.Println("Falou campeão!")
	default:
		fmt.Printf("Opção inválida: %v\n.", opcao)
	}
}
func validar(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("Número inválido")
	}
	return nil
}
