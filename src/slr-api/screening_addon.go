package main

import (
	"io/ioutil"
	"strings"
	"sync"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
	"gopkg.in/jdkato/prose.v2"
)

type ClassType int

const (
	ClassTypeExclude ClassType = 0
	ClassTypeInclude ClassType = 1
)

// GetScreeningMediaForProject retrieve the predicted values for given article
func GetScreeningMediaForProject(projectID uuid.UUID, title, abstract string) (*app.Articlescreening, error) {
	wg := &sync.WaitGroup{}
	var doc *prose.Document
	wg.Add(2)
	var abstractSanitized string
	var titleSanitized string
	go func() {
		var err error
		abstractSanitized, doc, err = SanitizeText(abstract)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	go func() {
		var err error
		titleSanitized, _, err = SanitizeText(title)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	wg.Wait()
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	model.Output = ioutil.Discard
	go model.OnlineLearn(errors)

	err := model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		log.Info("no model yet")
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
		class, p = model.Probability(abstractSanitized)
	}
	res.Abstract = &app.Textpredictmedia{
		Class:      getClass(class),
		Confidence: p,
		Text:       abstract,
	}
	if docuCount > 0 {
		class, p = model.Probability(titleSanitized + " " + abstractSanitized)
	}
	res.AbstractAndTitle = &app.Textpredictmedia{
		Class:      getClass(class),
		Confidence: p,
		Text:       abstract,
	}
	if docuCount > 0 {
		class, p = model.Probability(titleSanitized)
	}
	res.Title = &app.Textpredictmedia{
		Class:      getClass(class),
		Confidence: p,
		Text:       title,
	}

	for _, sentence := range doc.Sentences() {
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
	tf := text.TFIDF(*model)
	frequencies := tf.MostImportantWords(abstract, 10)
	for _, freq := range frequencies {
		res.MostImportantWords = append(res.MostImportantWords, &app.Mostimportantwordsmedia{
			Frequency: freq.Frequency,
			TfIdf:     freq.TFIDF,
			Word:      freq.Word,
		})
	}
	return res, nil
}

// TrainModel trains the model
func TrainModel(projectID uuid.UUID, abstract, title string, include bool) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		var err error
		abstract, _, err = SanitizeText(abstract)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	go func() {
		var err error
		title, _, err = SanitizeText(title)
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()
	wg.Wait()

	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)
	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	model.Output = ioutil.Discard
	go model.OnlineLearn(errors)

	err := model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		log.Info("no model yet")
	}

	identifier := ClassTypeExclude
	if include {
		identifier = ClassTypeInclude
	}
	// stream <- base.TextDatapoint{
	// 	X: strings.Trim(title, "") + " " + strings.Trim(abstract, ""),
	// 	Y: uint8(identifier),
	// }

	stream <- base.TextDatapoint{
		X: strings.Trim(abstract, ""),
		Y: uint8(identifier),
	}

	stream <- base.TextDatapoint{
		X: strings.Trim(title, ""),
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
// Removes:
//	- conjunction, subordinating or preposition (IN)
//	- infinitival to (TO)
//	- determiner (DT)
//	- conjunction, coordinating (CC)
func SanitizeText(text string) (string, *prose.Document, error) {
	text = strings.ReplaceAll(text, "...", "")
	doc, err := prose.NewDocument(text)
	if err != nil {
		return "", nil, err
	}
	res := ""
	for _, tok := range doc.Tokens() {
		if tok.Tag == "IN" {
			continue
		}
		if tok.Tag == "TO" {
			continue
		}
		if tok.Tag == "RB" || tok.Tag == "RBR" || tok.Tag == "RBS" || tok.Tag == "RP" {
			continue
		}
		if tok.Tag == "DT" {
			continue
		}
		if tok.Tag == "CC" {
			continue
		}
		if tok.Tag == "CD" {
			continue
		}
		// stem := porterstemmer.StemString(tok.Text) + " "
		// res += stem
		res += tok.Text + " "
	}
	return res, doc, nil
}
