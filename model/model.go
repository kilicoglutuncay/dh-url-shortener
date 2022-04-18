package model

type RedirectionData struct {
	OriginalURL string
	Hits        int
}

type ListData struct {
	Hash        string `json:"hash"`
	OriginalURL string `json:"original_url"`
	Hits        int    `json:"hits"`
}
