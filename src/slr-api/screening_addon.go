package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/gofrs/uuid"
	"github.com/reiver/go-porterstemmer"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"github.com/wimspaargaren/slr-automation/src/slr-api/models"
	"gopkg.in/jdkato/prose.v2"
)

const (
	AmountImportantWords int = 11
)

type ClassType int

const (
	ClassTypeExclude ClassType = 0
	ClassTypeInclude ClassType = 1
)

// VerifyModel verifies if the model is ready for automatic screening
func VerifyModel(projectID uuid.UUID) error {
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, base.OnlyWords)

	model.Output = ioutil.Discard
	go model.OnlineLearn(errors)
	err := model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		return fmt.Errorf("No model for article set yet")
	}

	close(stream)
	for {
		err, more := <-errors
		if more {
			log.Errorf("Error passed: %v", err)
		} else {
			// training is done!
			break
		}
	}
	for _, p := range model.Probabilities {
		if p == 0 {
			return fmt.Errorf("Both an included and excluded example should be provided to the model")
		}
	}
	return nil
}

// GetScreeningMediaForProject retrieve the predicted values for given article
func GetScreeningMediaForProject(projectID uuid.UUID, article *models.Article, abstractScreen bool) (*app.Articlescreening, error) {
	total, doc, err := getSanitizedText(article, abstractScreen)
	if err != nil {
		return nil, err
	}
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, base.OnlyWords)

	model.Output = ioutil.Discard
	go model.OnlineLearn(errors)
	sentenceTFIDF := ""
	sentenceTF := ""
	err = model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		log.Info("no model yet")
	} else {
		sentenceTFIDF = SentenceForTFIDF(model, total, doc)
		frequencies := TF(doc, AmountImportantWords)
		result := ""
		tfIDFWords := make(map[string]bool)
		for _, x := range frequencies {
			tfIDFWords[x.Word] = true
		}
		// Create sentence only of words occuring in idf set
		for _, tok := range doc.Tokens {
			sanitizedToken := SanitizeToken(tok)
			if sanitizedToken == "" {
				continue
			}
			if tfIDFWords[sanitizedToken] {
				result += sanitizedToken + " "
			}
		}
		sentenceTF = result
	}
	close(stream)
	for {
		err, more := <-errors
		if more {
			log.Errorf("Error passed: %v", err)
		} else {
			// training is done!
			break
		}
	}
	docuCount := model.DocumentCount
	class := uint8(0)
	p := float64(0)
	res := &app.Articlescreening{}
	if docuCount > 0 {
		class, p = model.Probability(sentenceTFIDF)
	}
	res.Tfidf = &app.Titleabstractpredictmedia{
		Class:      getClass(class),
		Confidence: p,
		Abstract:   article.Abstract,
		Title:      article.Title,
	}
	if docuCount > 0 {
		class, p = model.Probability(sentenceTF)

	}
	res.Tf = &app.Titleabstractpredictmedia{
		Class:      getClass(class),
		Confidence: p,
		Abstract:   article.Abstract,
		Title:      article.Title,
	}

	for _, sentence := range doc.Sentences {
		if docuCount > 0 {
			class, p = model.Probability(sentence.Text)
		}
		res.Sentences = append(res.Sentences, &app.Textpredictmedia{
			Class:      getClass(class),
			Confidence: p,
			Text:       sentence.Text,
		})
	}
	if docuCount == 0 {
		return res, nil
	}
	// nolint: govet
	tf := text.TFIDF(*model)
	frequencies := tf.MostImportantWords(article.Abstract, 10)
	for _, freq := range frequencies {
		res.MostImportantWords = append(res.MostImportantWords, &app.Mostimportantwordsmedia{
			Frequency: freq.Frequency,
			TfIdf:     freq.TFIDF,
			Word:      freq.Word,
		})
	}
	return res, nil
}

func TF(doc *models.ScreeningData, n int) text.Frequencies {
	tfinput := []string{}
	for _, tok := range doc.Tokens {
		sanitizedToken := SanitizeToken(tok)
		if sanitizedToken != "" {
			tfinput = append(tfinput, sanitizedToken)
		}
	}
	freq := text.TermFrequencies(tfinput)

	sort.Slice(freq[:], func(i, j int) bool {
		return freq[i].Frequency > freq[j].Frequency
	})

	if n > len(freq) {
		return freq
	}
	return freq[:n]
}

func SentenceForTFIDF(model *text.NaiveBayes, total string, doc *models.ScreeningData) string {
	// Calc TFIDF
	// nolint: govet
	tf := text.TFIDF(*model)
	result := ""
	temp := AmountImportantWords
	// Retrieve 11 most important words
	frequencies := tf.MostImportantWords(total, temp)
	if frequencies[len(frequencies)-1].TFIDF == frequencies[len(frequencies)-2].TFIDF {
		mem := frequencies[len(frequencies)-1].TFIDF
		for frequencies[len(frequencies)-1].TFIDF == mem {
			temp--

			if temp == 0 {
				return total
			}
			frequencies = tf.MostImportantWords(total, temp)

		}
	}
	tfIDFWords := make(map[string]bool)
	for _, x := range frequencies {
		tfIDFWords[x.Word] = true
	}
	// Create sentence only of words occuring in idf set
	for _, tok := range doc.Tokens {
		sanitizedToken := SanitizeToken(tok)
		if sanitizedToken == "" {
			continue
		}
		if tfIDFWords[sanitizedToken] {
			result += sanitizedToken + " "
		}
	}
	return result
}

// TrainModel trains the model
func TrainModel(projectID uuid.UUID, article *models.Article, screenAbstract, include bool) error {
	trainingSentence, doc, err := getSanitizedText(article, screenAbstract)
	if err != nil {
		return err
	}

	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)
	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	model.Output = ioutil.Discard
	go model.OnlineLearn(errors)
	sentenceTFIDF := trainingSentence
	err = model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		log.Info("no model yet")
	} else {
		sentenceTFIDF = SentenceForTFIDF(model, trainingSentence, doc)
	}
	identifier := ClassTypeExclude
	if include {
		identifier = ClassTypeInclude
	}

	stream <- base.TextDatapoint{
		X: strings.Trim(sentenceTFIDF, ""),
		Y: uint8(identifier),
	}

	close(stream)

	for {
		err, more := <-errors
		if more {
			log.Errorf("Error passed: %v", err)
		} else {
			// training is done!
			break
		}
	}
	err = model.PersistToFile("screening-models/" + projectID.String())
	if err != nil {
		return err
	}
	return nil
}

func getSanitizedText(article *models.Article, abstractScreen bool) (string, *models.ScreeningData, error) {
	var doc *models.ScreeningData
	var err error
	if abstractScreen {
		doc, err = article.GetAbstractDoc()
		if err != nil {
			return "", nil, err
		}
	} else {
		doc, err = article.GetFullTextDoc()
		if err != nil {
			return "", nil, err
		}
	}
	result := SanitizeText(doc)
	return result, doc, nil
}

func getClass(class uint8) string {
	switch int(class) {
	case 0:
		return "Exclude"
	case 1:
		return "Include"
	default:
		return "Unknown"
	}
}

// SanitizeText sanitizes text for training and screening
func SanitizeText(doc *models.ScreeningData) string {
	res := ""
	for _, tok := range doc.Tokens {
		sanitized := SanitizeToken(tok)
		if sanitized != "" {
			res += sanitized + " "
		}
	}
	return res
}

func SanitizeToken(tok prose.Token) string {
	if tok.Tag == "IN" ||
		tok.Tag == "RB" || tok.Tag == "RBR" || tok.Tag == "RBS" || tok.Tag == "RP" ||
		tok.Tag == "CC" ||
		tok.Tag == "CD" ||
		tok.Tag == "DT" ||
		tok.Tag == "PRP" ||
		tok.Tag == "." ||
		tok.Tag == "(" ||
		tok.Tag == ")" ||
		tok.Tag == "," ||
		tok.Tag == "TO" {
		return ""
	}
	stem := porterstemmer.StemString(tok.Text)
	return stem
}
