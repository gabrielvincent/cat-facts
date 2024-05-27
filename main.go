package main

import (
	components "catfacts/components"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type CatFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type CatFactsApiResponse struct {
	Data []CatFact `json:"data"`
}

func (c CatFact) Render(
	ctx context.Context,
	w http.ResponseWriter,
	other string,
	otherTwo int,
) {
	w.Write([]byte(c.Fact))
}

func getFacts() ([]CatFact, error) {
	// Simulate a slow API call
	time.Sleep(2 * time.Second)

	apiUrl := "https://catfact.ninja/facts?limit=10"
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var facts CatFactsApiResponse
	err = json.Unmarshal(body, &facts)
	if err != nil {
		return nil, err
	}

	return facts.Data, nil
}

func factsHandler(w http.ResponseWriter, r *http.Request) {
	flusher, canFlush := w.(http.Flusher)
	ctx := context.Background()

	if canFlush {
		// Se puder usar Flush, a gente já começa enviando a página estática.
		components.Index(nil).Render(ctx, w)
		flusher.Flush()
	}

	catFacts, err := getFacts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var facts []string
	for _, catFact := range catFacts {
		facts = append(facts, catFact.Fact)
	}

	if canFlush {
		// Envia somente o conteúdo dos fatos sobre gatos.
		components.Facts(facts).Render(ctx, w)
	} else {
		// Como Flush não está disponível, só nos resta enviar o conteúdo
		// completo após o carregamento dos fatos sobre gatos.
		components.Index(facts).Render(ctx, w)
	}

}

func main() {
	http.HandleFunc("/", factsHandler)
	http.HandleFunc(
		"/tailwind.css",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/css")
			http.ServeFile(w, r, "public/css/tailwind.css")
		},
	)
	http.HandleFunc(
		"/lazy-component.js",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/css")
			http.ServeFile(w, r, "public/js/lazy-component.js")
		},
	)
	http.ListenAndServe(":3000", nil)
}
