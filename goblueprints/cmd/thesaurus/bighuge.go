package thesaurus

import (
	"encoding/json"
	"errors"
	"net/http"
)

type BigHuge struct {
	APIKey string
}

type synonyms struct {
	Noun *word `json:"noun"`
	Verb *word `json:"verb"`
}

type word struct {
	Syn []string `json:"syn"`
}

func (b *BigHuge) Synonyms(term string) ([]string, error) {
	var syns []string
	res, err := http.Get("http://words.bighugelabs.com/api/2/" + b.APIKey + "/" + term + "/json")
	if err != nil {
		return syns, errors.New("bighuge: Failed when looking for synonyms	for " + term + ": " + err.Error())
	}
	defer res.Body.Close()
	var data synonyms
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return syns, err
	}
	if data.Noun != nil {
		syns = append(syns, data.Noun.Syn...)
	}
	if data.Verb != nil {
		syns = append(syns, data.Verb.Syn...)
	}
	return syns, nil
}
