package cmd

import "testing"

func TestObterNomeProjeto(t *testing.T) {
	t.Log("quando informada url valida, entao retorna nome do projeto.")
	{
		var projeto = obterNomeProjeto("https://gitlab.com/hfantin/teste/-/merge_requests/1")
		if projeto == "teste" {
			t.Log("\nNome do projeto obtido com sucesso!")
		} else {
			t.Fatal("\tNome do projeto inválido.", projeto)
		}
	}
	t.Log("quando informada url invalida, entao retorna nome do projeto em branco.")
	{
		var projeto = obterNomeProjeto("https://gitlab.com/teste")
		if projeto == "" {
			t.Log("\nUrl é invalida, nome do projeto em branco!")
		} else {
			t.Fatal("\tNome do projeto inválido.", projeto)
		}
	}
}
