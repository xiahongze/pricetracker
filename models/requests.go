package models

import "cloud.google.com/go/datastore"

type (
	// Options is the options for an entry
	Options struct {
		User       string  `json:"user"`
		CheckFreq  int16   `json:"checkFreq"` // in minutes
		AlertType  string  `json:"alertType"`
		Threshold  float32 `json:"threshold"`
		MaxRecords int16   `json:"maxRecords"`
		UseChrome  *bool   `json:"useChrome"`
	}

	// CreateRequest defines the contract to add an entry
	CreateRequest struct {
		URL           string   `json:"url"`
		XPATH         string   `json:"xpath"`
		Name          string   `json:"name"`
		ExpectedPrice string   `json:"expectedPrice"`
		Options       *Options `json:"options"`
	}

	// UpdateRequest defines the contract to update an entry
	UpdateRequest struct {
		URL     string         `json:"url"`
		XPATH   string         `json:"xpath"`
		Name    string         `json:"name"`
		Key     *datastore.Key `json:"key"`
		Options *Options       `json:"options"`
	}

	// ReadOrDelRequest defines the contract to read/delete an entry
	ReadOrDelRequest struct {
		Key *datastore.Key `json:"key"`
	}
)

// Validate validates
func (r *CreateRequest) Validate() (string, bool) {
	if r.URL == "" {
		return "url is not set", false
	}
	if r.XPATH == "" {
		return "xpath is not set", false
	}
	if r.Name == "" {
		return "name is not set", false
	}
	if r.ExpectedPrice == "" {
		return "expectedPrice is not set", false
	}
	r.Options.setDefault()
	return "", true
}

func (o *Options) setDefault() {
	if o.CheckFreq == 0 {
		o.CheckFreq = 60 * 24
	}
	if o.MaxRecords == 0 {
		o.MaxRecords = 365
	}
	if o.AlertType == "" && o.Threshold == 0 {
		o.AlertType = "onChange"
	}
	if o.AlertType == "" && o.Threshold != 0 {
		o.AlertType = "threshold"
	}
}

// Validate validates
func (r *ReadOrDelRequest) Validate() (string, bool) {
	if r.Key == nil {
		return "key is not given", false
	}
	return "", true
}

// Validate validates
func (r *UpdateRequest) Validate() (string, bool) {
	if r.Key == nil {
		return "key is not given", false
	}
	if r.Options == nil {
		return "options is not given", false
	}
	return "", true
}
