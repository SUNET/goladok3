package ladoktypes

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/golden"
)

func TestFeedIDTrim(t *testing.T) {
	tts := []struct {
		name string
		have FeedID
		want FeedID
	}{
		{
			name: "OK",
			have: "urn:id:4856",
			want: "4856",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.have.trim())
		})
	}
}

func TestFeedIDInt(t *testing.T) {
	tts := []struct {
		name    string
		have    FeedID
		want    int
		wantErr bool
	}{
		{name: "valid", have: "4856", want: 4856},
		{name: "zero", have: "0", want: 0},
		{name: "invalid", have: "abc", wantErr: true},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.have.int()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestKontaktuppgifterEventParse(t *testing.T) {
	e := &KontaktuppgifterEvent{
		HandelseUID: "uid-123",
		EventContext: EventContext{
			AnvandareUID: "user-1",
			Anvandarnamn: "test@example.com",
			LarosateID:   "27",
		},
		Handelsetyp: "UPPDATERAD",
		Epostadress: "mail@example.com",
		StudentUID:  "student-1",
		Postadresser: []struct {
			Text             string `xml:",chardata"`
			Land             string `xml:"Land"`
			PostadressTyp    string `xml:"PostadressTyp"`
			Postnummer       string `xml:"Postnummer"`
			Postort          string `xml:"Postort"`
			Utdelningsadress string `xml:"Utdelningsadress"`
			CareOf           string `xml:"CareOf"`
		}{
			{
				Land:             "Sverige",
				PostadressTyp:    "POSTADRESS",
				Postnummer:       "10010",
				Postort:          "Stockholm",
				Utdelningsadress: "Gatan 1",
				CareOf:           "Name",
			},
		},
		Telefonnummer: "0701234567",
	}

	got := e.Parse("entry-1")

	assert.Equal(t, KontaktuppgifterEventName, got.EventTypeName)
	assert.Equal(t, "entry-1", got.EntryID)
	assert.Equal(t, "uid-123", got.HandelseUID)
	assert.Equal(t, "mail@example.com", got.Email)
	assert.Equal(t, "student-1", got.StudentUID)
	assert.Equal(t, "0701234567", got.Telefonnummer)
	assert.Len(t, got.Postadresser, 1)
	assert.Equal(t, "Sverige", got.Postadresser[0].Land)
	assert.Equal(t, "Name", got.Postadresser[0].CareOf)
}

func TestLokalStudentEventParse(t *testing.T) {
	e := &LokalStudentEvent{
		HandelseUID: "uid-1",
		EventContext: struct {
			Text         string `xml:",chardata"`
			AnvandareUID string `xml:"AnvandareUID"`
			Anvandarnamn string `xml:"Anvandarnamn"`
			LarosateID   string `xml:"LarosateID"`
		}{
			AnvandareUID: "user-1",
			Anvandarnamn: "test@example.com",
			LarosateID:   "27",
		},
		Handelsetyp:       "SKAPAD",
		Efternamn:         "Svensson",
		Fornamn:           "Erik",
		StudentUID:        "student-1",
		ExterntStudentUID: "ext-1",
		Fodelsedata:       "1990-01-01",
		Kon:               "2",
		Personnummer:      "199001011234",
	}

	got := e.Parse("entry-1")

	assert.Equal(t, LokalStudentEventName, got.EventTypeName)
	assert.Equal(t, "Svensson", got.Efternamn)
	assert.Equal(t, "Erik", got.Fornamn)
	assert.Equal(t, "199001011234", got.Personnummer)
	assert.Equal(t, "2", got.Kon)
}

func TestExternPartEventParse(t *testing.T) {
	e := &ExternPartEvent{
		HandelseUID: "uid-1",
		EventContext: EventContext{
			AnvandareUID: "user-1",
			Anvandarnamn: "test@example.com",
			LarosateID:   "-1",
		},
		EventTyp:          "SKAPAD",
		Giltighetsperiod:  "2021-01-01",
		ID:                "12345",
		Kod:               "CODE1",
		LandID:            "25",
		TypAvExternPartID: "1",
	}

	got := e.Parse("entry-1")

	assert.Equal(t, "ExternPartEvent", got.EventTypeName)
	assert.Equal(t, "12345", got.ID)
	assert.Equal(t, "CODE1", got.Kod)
	assert.Equal(t, "25", got.LandID)
}

func TestResultatEventParse(t *testing.T) {
	e := &ResultatEvent{
		HandelseUID: "uid-1",
		EventContext: struct {
			Text         string `xml:",chardata"`
			AnvandareUID string `xml:"AnvandareUID"`
			Anvandarnamn string `xml:"Anvandarnamn"`
			LarosateID   string `xml:"LarosateID"`
		}{
			AnvandareUID: "user-1",
			Anvandarnamn: "test@example.com",
			LarosateID:   "27",
		},
		Beslut: struct {
			Text              string `xml:",chardata"`
			BeslutUID         string `xml:"BeslutUID"`
			Beslutsdatum      string `xml:"Beslutsdatum"`
			Beslutsfattare    string `xml:"Beslutsfattare"`
			BeslutsfattareUID string `xml:"BeslutsfattareUID"`
		}{
			BeslutUID:         "beslut-1",
			Beslutsdatum:      "2021-10-01",
			Beslutsfattare:    "Teacher",
			BeslutsfattareUID: "teacher-1",
		},
		KursUID:          "kurs-1",
		KursinstansUID:   "ki-1",
		KurstillfalleUID: "kt-1",
		Resultat: struct {
			Text               string `xml:",chardata"`
			BetygsgradID       string `xml:"BetygsgradID"`
			BetygsskalaID      string `xml:"BetygsskalaID"`
			Examinationsdatum  string `xml:"Examinationsdatum"`
			GiltigSomSlutbetyg string `xml:"GiltigSomSlutbetyg"`
			OmfattningsPoang   string `xml:"OmfattningsPoang"`
			PrestationsPoang   string `xml:"PrestationsPoang"`
			ResultatUID        string `xml:"ResultatUID"`
		}{
			BetygsgradID:       "2302",
			BetygsskalaID:      "2",
			Examinationsdatum:  "2021-10-01",
			GiltigSomSlutbetyg: "true",
			OmfattningsPoang:   "7.5",
			PrestationsPoang:   "7.5",
			ResultatUID:        "res-1",
		},
		StudentUID:            "student-1",
		UtbildningsinstansUID: "utb-1",
	}

	got := e.Parse("ResultatPaModulAttesteratEvent", "entry-1")

	assert.Equal(t, "ResultatPaModulAttesteratEvent", got.EventTypeName)
	assert.Equal(t, "kurs-1", got.KursUID)
	assert.Equal(t, "beslut-1", got.Beslut.BeslutUID)
	assert.Equal(t, "7.5", got.Resultat.OmfattningsPoang)
	assert.Equal(t, "student-1", got.StudentUID)
}

func TestAnvandareEventParse(t *testing.T) {
	e := &AnvandareEvent{
		HandelseUID: "uid-1",
		EventContext: EventContext{
			AnvandareUID: "user-1",
			Anvandarnamn: "system@ladok.se",
			LarosateID:   "27",
		},
		AnvandareUID:   "anvandare-1",
		Anvandarnamnet: "test@school.se",
		Efternamn:      "Svensson",
		Fornamn:        "Erik",
	}

	got := e.Parse("AnvandareAndradEvent", "entry-1")

	assert.Equal(t, "AnvandareAndradEvent", got.EventTypeName)
	assert.Equal(t, "anvandare-1", got.AnvandareUID)
	assert.Equal(t, "test@school.se", got.Anvandarnamnet)
	assert.Equal(t, "Svensson", got.Efternamn)
}

func TestFeedParseAllEventTypes(t *testing.T) {
	feedXML := golden.Get(t, "feed_all_events.xml")

	f := &Feed{}
	err := xml.Unmarshal(feedXML, f)
	assert.NoError(t, err)

	superFeed, err := f.Parse()
	assert.NoError(t, err)
	assert.Equal(t, 100, superFeed.ID)
	assert.Len(t, superFeed.SuperEvents, 7)

	names := make([]string, len(superFeed.SuperEvents))
	for i, e := range superFeed.SuperEvents {
		names[i] = e.EventTypeName
	}
	assert.Contains(t, names, AnvandareAndradEventName)
	assert.Contains(t, names, AnvandareSkapadEventName)
	assert.Contains(t, names, ExternPartEventName)
	assert.Contains(t, names, KontaktuppgifterEventName)
	assert.Contains(t, names, ResultatPaModulAttesteratEventName)
	assert.Contains(t, names, ResultatPaHelKursAttesteratEventName)
	assert.Contains(t, names, LokalStudentEventName)
}

func TestFeedParseInvalidID(t *testing.T) {
	f := &Feed{ID: "urn:id:notanumber"}
	_, err := f.Parse()
	assert.Error(t, err)
}

func TestFeedParseEmptyEntries(t *testing.T) {
	f := &Feed{ID: "urn:id:42"}
	superFeed, err := f.Parse()
	assert.NoError(t, err)
	assert.Equal(t, 42, superFeed.ID)
	assert.Empty(t, superFeed.SuperEvents)
}
