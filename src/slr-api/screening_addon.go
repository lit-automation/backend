package main

import (
	"strings"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
)

type ClassType int

const (
	ClassTypeExclude ClassType = 0
	ClassTypeInclude ClassType = 1
)

// GetScreeningMediaForProject retrieve the predicted values for given article
func GetScreeningMediaForProject(projectID uuid.UUID, title, abstract string) *app.Articlescreening {
	abstract = enhanceAbstract(abstract)
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
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
		class, p = model.Probability(abstract)
	}
	res.Abstract = &app.Textpredictmedia{
		Class:      getClass(class),
		Confidence: p,
		Text:       abstract,
	}
	if docuCount > 0 {
		class, p = model.Probability(title)
	}
	res.Title = &app.Textpredictmedia{
		Class:      getClass(class),
		Confidence: p,
		Text:       title,
	}

	splittedAbstract := strings.Split(abstract, ".")
	for _, sentence := range splittedAbstract {
		if docuCount > 0 {
			class, p = model.Probability(sentence)
		}
		res.Sentences = append(res.Sentences, &app.Textpredictmedia{
			Class:      getClass(class),
			Confidence: p,
			Text:       sentence,
		})
	}

	return res
}

// TrainModel trains the model
func TrainModel(projectID uuid.UUID, abstract, title string, include bool) error {
	abstract = enhanceAbstract(abstract)
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	go model.OnlineLearn(errors)

	err := model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		log.Info("no model yet")
	}

	identifier := ClassTypeExclude
	if include {
		identifier = ClassTypeInclude
	}

	splittedAbstract := strings.Split(abstract, ".")
	for _, sentence := range splittedAbstract {
		stream <- base.TextDatapoint{
			X: strings.Trim(sentence, ""),
			Y: uint8(identifier),
		}
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

func enhanceAbstract(abstract string) string {
	return strings.ReplaceAll(abstract, "...", "")
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
