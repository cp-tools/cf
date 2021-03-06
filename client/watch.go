package cln

import (
	cfg "cf/config"
	pkg "cf/packages"

	"bytes"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	// Submission holds various data of
	// a particular submission
	Submission struct {
		ID, When, Name, Lang, Waiting,
		Verdict, Time, Memory string
	}
	// Problem holds problem solved status
	// based on current user session
	Problem struct {
		ID, Name, Status,
		Count string
	}
)

// WatchSubmissions finds all submissions in contest that matches query string
// query = problem to fetch all submissions in a particular problem (should be uppercase)
// query = submitID to fetch submission of given submission id
func WatchSubmissions(contest, query string, link url.URL) ([]Submission, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = pkg.RedirectCheck
	// fetch all submissions in contest
	link.Path = path.Join(link.Path, "my")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return nil, err
	} else if len(body) == 0 {
		// such page doesn't exist
		err = fmt.Errorf("Contest %v doesn't exist", contest)
		return nil, err
	}
	// to hold all submissions
	var data []Submission

	query = strings.ToUpper(query)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	sel := doc.Find("tr[data-submission-id]").Has("a[href*=\"/" + query + "\"]")
	sel.Each(func(_ int, row *goquery.Selection) {
		// select cell ...type(x) from row
		data = append(data, Submission{
			ID:      pkg.GetText(row, "td:nth-of-type(1)"),
			When:    pkg.GetText(row, "td:nth-of-type(2)"),
			Name:    pkg.GetText(row, "td:nth-of-type(4)"),
			Lang:    pkg.GetText(row, "td:nth-of-type(5)"),
			Waiting: pkg.GetAttr(row, "td:nth-of-type(6)", "waiting"),
			Verdict: pkg.GetText(row, "td:nth-of-type(6)"),
			Time:    pkg.GetText(row, "td:nth-of-type(7)"),
			Memory:  pkg.GetText(row, "td:nth-of-type(8)"),
		})
	})

	return data, nil
}

// WatchContest parses contest solved count status
func WatchContest(contest string, link url.URL) ([]Problem, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = pkg.RedirectCheck
	// fetch contest dashboard page
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return nil, err
	} else if len(body) == 0 {
		// such page doesn't exist
		err = fmt.Errorf("Contest %v doesn't exist", contest)
		return nil, err
	}
	// to hold all problems in contest
	var data []Problem

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find(".problems tr").Has("td").Each(func(_ int, row *goquery.Selection) {

		data = append(data, Problem{
			ID:     pkg.GetText(row, "td:nth-of-type(1)"),
			Name:   pkg.GetText(row, "td:nth-of-type(2) a"),
			Count:  pkg.GetText(row, "td:nth-of-type(4)"),
			Status: row.AttrOr("class", ""),
		})
	})
	return data, nil
}
