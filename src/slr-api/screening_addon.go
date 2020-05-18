package main

import (
	"fmt"
	"strings"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/gofrs/uuid"
	"github.com/wimspaargaren/slr-automation/src/slr-api/app"
)

type ClassType int

const (
	ClassTypeExclude ClassType = 0
	ClassTypeInclude ClassType = 1
)

// GetScreeningMediaForProject retrieve the predicted values for given article
func GetScreeningMediaForProject(projectID uuid.UUID, title, abstract string) *app.Articlescreening {
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	go model.OnlineLearn(errors)

	err := model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		fmt.Println("no model yet")
	}

	close(stream)

	for {
		err, more := <-errors
		if more {
			fmt.Printf("Error passed: %v", err)
		} else {
			// training is done!
			break
		}
	}
	res := &app.Articlescreening{}
	class, p := model.Probability(abstract)
	res.Abstract = &app.Textpredictmedia{
		Class:      int(class),
		Confidence: p,
		Text:       abstract,
	}

	class, p = model.Probability(title)
	res.Title = &app.Textpredictmedia{
		Class:      int(class),
		Confidence: p,
		Text:       abstract,
	}

	splittedAbstract := strings.Split(abstract, ".")
	for _, sentence := range splittedAbstract {
		class, p := model.Probability(sentence)
		fmt.Println("P", p)
		fmt.Println("Clas", class)
		res.Sentences = append(res.Sentences, &app.Textpredictmedia{
			Class:      int(class),
			Confidence: p,
			Text:       sentence,
		})
	}

	return res
}

// TrainModel trains the model
func TrainModel(projectID uuid.UUID, abstract, title string, include bool) error {
	stream := make(chan base.TextDatapoint, 100)
	errors := make(chan error)

	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	go model.OnlineLearn(errors)

	err := model.RestoreFromFile("screening-models/" + projectID.String())
	if err != nil {
		fmt.Println("no model yet")
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

	close(stream)

	for {
		err, more := <-errors
		if more {
			fmt.Printf("Error passed: %v", err)
		} else {
			// training is done!
			break
		}
	}
	err = model.PersistToFile("savemodel")
	if err != nil {
		return err
	}
	return nil
}
