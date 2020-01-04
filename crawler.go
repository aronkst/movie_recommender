package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func loadSite(url string) (*goquery.Document, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		message := fmt.Sprintf("Status code error: %d %s", response.StatusCode, response.Status)
		err := errors.New(message)
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func getValueFromSiteDocument(document *goquery.Document, selector string, attribute string) string {
	selection := document.Selection
	return getValueFromSiteSelection(selection, selector, attribute)
}

func getValueFromSiteSelection(selection *goquery.Selection, selector string, attribute string) string {
	var value string

	if attribute == "" {
		value = selection.Find(selector).Text()
	} else {
		value, _ = selection.Find(selector).Attr(attribute)
	}

	return clearString(value)
}

func getValueFromSiteSelectionInside(selection *goquery.Selection, attribute string) string {
	var value string

	if attribute == "" {
		value = selection.Text()
	} else {
		value, _ = selection.Attr(attribute)
	}

	return clearString(value)
}
