package main

import "strings"

func cleanWords(sentence string) string {
	words := strings.Split(sentence, " ")
	new_sentence := make([]string, len(words))
	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			new_sentence[i] = "****"
		default:
			new_sentence[i] = word
		}
	}
	return strings.Join(new_sentence, " ")
}