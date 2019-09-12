package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// MsIgniteAPIResponse is the JSON API response layout
type MsIgniteAPIResponse struct {
	Data   []SessionData `json:"data"`
	Facets Facets        `json:"facets"`
	Total  int           `json:"total"`
}

// SessionData are all the fields returned for a single session
type SessionData struct {
	Score                     float32   `json:"@search.score"`
	SessionID                 string    `json:"sessionId"`
	SessionInstanceID         string    `json:"sessionInstanceId"`
	SessionCode               string    `json:"sessionCode"`
	Title                     string    `json:"title"`
	SortTitle                 string    `json:"sortTitle"`
	SortRank                  int       `json:"sortRank"`
	Description               string    `json:"description"`
	RegistrationLink          string    `json:"registrationLink"`
	StartDateTime             string    `json:"startDateTime"`
	EndDateTime               string    `json:"endDateTime"`
	DuractionInMinutes        int       `json:"durationInMinutes"`
	SessionType               string    `json:"sessionType"`
	SessionTypeLogical        string    `json:"sessionTypeLogical"`
	LearningPath              []string  `json:"learningPath"`
	Level                     string    `json:"level"`
	Products                  []string  `json:"products"`
	Format                    string    `json:"format"`
	Topic                     string    `json:"topic"`
	SessionTypeID             string    `json:"sessionTypeId"`
	IsMandatory               bool      `json:"isMandatory"`
	VisibleInSessionListing   bool      `json:"visibleInSessionListing"`
	TechCommunityDiscussionID string    `json:"techCommunityDiscussionId"`
	SpeakerIDs                []string  `json:"speakerIds"`
	SpeakerNames              []string  `json:"speakerNames"`
	SpeakerCompanies          []string  `json:"speakerCompanies"`
	SessionLink               []string  `json:"sessionLinks"`
	MarketingCampaign         []string  `json:"marketingCampaign"`
	Links                     string    `json:"links"`
	LastUpdate                time.Time `json:"lastUpdate"`
	ChildModules              []string  `json:"childModules"`
	SiblingModules            []string  `json:"siblingModules"`
}

//Facets describes fields that can be filtered on
type Facets struct {
	DurationInMinutes FacetNum `json:"durationInMinutes"`
	SessionType       Facet    `json:"sessionType"`
	LearningPath      Facet    `json:"learningPath"`
	Level             Facet    `json:"level"`
	Products          Facet    `json:"products"`
	Format            Facet    `json:"format"`
	Topic             Facet    `json:"topic"`
	SessionTypeID     Facet    `json:"sessionTypeId"`
}

//Facet is a single facet with string filter values
type Facet struct {
	DisplayName string   `json:"displayName"`
	FacetName   string   `json:"facetName"`
	Visibility  bool     `json:"isVisible"`
	Filters     []Filter `json:"filters"`
}

//FacetNum is a single facet with integer filter values
type FacetNum struct {
	DisplayName string      `json:"displayName"`
	FacetName   string      `json:"facetName"`
	Visibility  bool        `json:"isVisible"`
	Filters     []FilterNum `json:"filters"`
}

//Filter applied to string values
type Filter struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

//FilterNum applied to integer values
type FilterNum struct {
	Value int `json:"value"`
	Count int `json:"count"`
}

func main() {

	sessions := PostSearchAPI()
	WriteSessionDataCSV(sessions)

}

// PostSearchAPI runs a search against the MS Ignite Search API and returns the JSON response
func PostSearchAPI() MsIgniteAPIResponse {
	url := "https://api-myignite.techcommunity.microsoft.com/api/session/search"

	payload := strings.NewReader("{\"itemsPerPage\": 600, \"searchText\": \"*\", \"searchPage\": 1, \"sortOption\": \"ASC\", \"searchFacets\": { \"facets\": [ { \"facetName\": \"sessionType\", \"displayName\": \"Breakout: 75 Minute\", \"names\": [ \"Breakout: 75 Minute\", \"Breakout: 45 Minute\", \"Theater: 20 Minute\" ] }, { \"facetName\": \"format\", \"displayName\": \"Session\", \"names\": [ \"Session\", \"Partner Led Session\", \"Panel Discussion\", \"Customer Showcase\" ] } ], \"personalizationFacets\": [], \"dateFacet\": [ { \"startDateTime\": \"2019-11-03T13:30:00.000Z\", \"endDateTime\": \"2019-11-03T23:59:59.000Z\" }, { \"startDateTime\": \"2019-11-04T13:30:00.000Z\", \"endDateTime\": \"2019-11-04T23:59:59.000Z\" }, { \"startDateTime\": \"2019-11-05T13:30:00.000Z\", \"endDateTime\": \"2019-11-05T23:59:59.000Z\" }, { \"startDateTime\": \"2019-11-06T13:30:00.000Z\", \"endDateTime\": \"2019-11-06T23:59:59.000Z\" }, { \"startDateTime\": \"2019-11-07T13:30:00.000Z\", \"endDateTime\": \"2019-11-07T23:59:59.000Z\" }, { \"startDateTime\": \"2019-11-08T13:30:00.000Z\", \"endDateTime\": \"2019-11-08T23:59:59.000Z\" } ] }, \"recommendedItemIds\": [], \"favoritesIds\": [], \"mustHaveOnDemandVideo\": false}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "api-myignite.techcommunity.microsoft.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal("Expect to receive 200 OK. Got ", res.Status)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var sessions MsIgniteAPIResponse
	err = json.Unmarshal(body, &sessions)
	if err != nil {
		log.Fatal(err)
	}

	if len(sessions.Data) != sessions.Total {
		fmt.Printf("Incomplete results. Received %d of %d search results. Increase items per page in search request", len(sessions.Data), sessions.Total)
	}

	return sessions
}

//WriteSessionDataCSV prints select fields from the API response to a CSV file
func WriteSessionDataCSV(sessions MsIgniteAPIResponse) {
	breakoutCSVFile, err := os.Create("ignite_breakout_sessions.csv")
	if err != nil {
		log.Fatal("Unable to create breakout session file. Received error ", err)
	}
	defer breakoutCSVFile.Close()
	theaterCSVFile, err := os.Create("ignite_theater_sessions.csv")
	if err != nil {
		log.Fatal("Unable to create theater session file. Received error ", err)
	}
	defer theaterCSVFile.Close()
	// w := csv.NewWriter(os.Stdout)
	breakoutW := csv.NewWriter(breakoutCSVFile)
	theaterW := csv.NewWriter(theaterCSVFile)

	header := []string{"Session ID", "Session Code", "Title", "Session Type", "Level", "Format", "Speaker Names", "Last Update"}

	if err = breakoutW.Write(header); err != nil {
		log.Fatalln("error writing header to breakout csv:", err)
	}
	if err = theaterW.Write(header); err != nil {
		log.Fatalln("error writing header to theater csv:", err)
	}

	for _, session := range sessions.Data {
		speakers := strings.Join(session.SpeakerNames, ";")
		switch {
		case strings.Contains(session.SessionType, "Breakout"):
			if err := breakoutW.Write([]string{session.SessionID, session.SessionCode, session.Title, session.SessionType, session.Level, session.Format, speakers, session.LastUpdate.Format(time.RFC3339)}); err != nil {
				log.Fatalln("error writing session info to csv:", err)
			}
			// Write any buffered data to the underlying writer (standard output).
			breakoutW.Flush()
		case strings.Contains(session.SessionType, "Theater"):
			if err := theaterW.Write([]string{session.SessionID, session.SessionCode, session.Title, session.SessionType, session.Level, session.Format, speakers, session.LastUpdate.Format(time.RFC3339)}); err != nil {
				log.Fatalln("error writing session info to csv:", err)
			}
			// Write any buffered data to the underlying writer (standard output).
			theaterW.Flush()
		default:
			fmt.Println("No output file defined for session type: ", session.SessionType)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	breakoutW.Flush()
	theaterW.Flush()

	if err = breakoutW.Error(); err != nil {
		log.Fatal(err)
	}
	if err = theaterW.Error(); err != nil {
		log.Fatal(err)
	}
}
