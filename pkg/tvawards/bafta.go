package tvawards

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type NomineeStruct struct {
	FilmName  string
	ActorName string
	Winner    bool
}

func getMinBaftaYear(contentType string) int {
	switch contentType {
	case "television":
		return 1954
	case "film":
		return 1954
	case "games":
		return 2004
	case "tvcraft":
		return 1978
	case "childrens":
		return 1996
	case "cymru":
		return 1993
	default:
		return 0
	}
}

// Get BaftaAwards - Gets all BAFTA awards from a content type (e.g. television)
// Returns a map of years, to award names, to an array of nominees/winner(s)
func GetBaftaAwards(contentType string) (*map[int]map[string][]NomineeStruct, error) {
	minYear := getMinBaftaYear(contentType)

	// Unsupported content type
	if minYear == 0 {
		return nil, nil
	}

	finalResponse := make(map[int]map[string][]NomineeStruct)

	for i := minYear; i < 2022; i++ {
		resp, err := getBaftaByYear(i, contentType)

		if err != nil {
			continue
		}

		finalResponse[i] = resp
	}

	return &finalResponse, nil
}

// getBaftaByYear - returns a map of strings containing all award names,
// along with their nominees, and winners
func getBaftaByYear(year int, contentType string) (map[string][]NomineeStruct, error) {
	resp, err := http.Get("http://awards.bafta.org/award/" + strconv.Itoa(year) + "/" + contentType)

	if err != nil {
		return map[string][]NomineeStruct{}, err
	}

	defer resp.Body.Close()

	pageReader, _ := goquery.NewDocumentFromReader(resp.Body)

	sectionResult := make(map[string][]NomineeStruct)

	pageReader.Find(".search-result-wrapper").Each(func(u int, s *goquery.Selection) {
		sectionTitle := formatBaftaTitle(s.Find(".search-result-title").Text())

		headlines := make([]string, 0, 10)
		subtitle := make([]string, 0, 10)

		s.Find(".search-result-headline").Each(func(i int, s *goquery.Selection) {
			headlines = append(headlines, strings.TrimSpace(s.Text()))
		})

		s.Find(".search-result-subtitle").Each(func(i int, s *goquery.Selection) {
			subtitle = append(subtitle, strings.TrimSpace(s.Text()))
		})

		nomineeList := make([]NomineeStruct, 0, 5)

		for i, v := range headlines {
			if v == "" && subtitle[i] == "" {
				continue
			}

			nomineeList = append(nomineeList, NomineeStruct{
				FilmName:  v,
				ActorName: subtitle[i],
				Winner:    i == 0,
			})
		}

		sectionResult[sectionTitle] = nomineeList
	})

	return sectionResult, nil
}

var titleStartRegexp = regexp.MustCompile(`[Television|Film] \| `)
var yearAtEndRegexp = regexp.MustCompile(`( ?in ?)\d\d\d\d$`)

// formatBaftaTitle - removes the leading `television` or `film`
// and the ending `in YYYY` from the input string
func formatBaftaTitle(input string) string {
	matches := titleStartRegexp.FindAllStringIndex(input, 2)

	if matches != nil {
		input = input[matches[0][1]:]
	}

	matches = yearAtEndRegexp.FindAllStringIndex(input, 2)

	if matches != nil {
		input = input[:matches[0][0]]
	}

	return input
}
