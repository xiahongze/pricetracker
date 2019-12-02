package models

import "cloud.google.com/go/datastore"

import "fmt"

type (
	// Options is the options for an entry
	Options struct {
		User       string  `json:"user"`
		CheckFreq  int16   `json:"checkFreq"` // in minutes
		AlertType  string  `json:"alertType"`
		Threshold  float32 `json:"threshold"`
		MaxRecords int16   `json:"maxRecords"`
		UseChrome  bool    `json:"useChrome"`
	}

	// CreateRequest defines the contract to add an entry
	CreateRequest struct {
		URL           string  `json:"url"`
		XPATH         string  `json:"xpath"`
		Name          string  `json:"name"`
		ExpectedPrice string  `json:"expectedPrice"`
		Options       Options `json:"options,omitempty"`
	}

	// UpdateRequest defines the contract to update an entry
	UpdateRequest struct {
		URL       string         `json:"url"`
		XPATH     string         `json:"xpath"`
		Name      string         `json:"name"`
		Key       *datastore.Key `json:"key"`
		UseChrome *bool          `json:"useChrome,omitempty"`
		Options   *Options       `json:"options,omitempty"`
	}

	// ReadOrDelRequest defines the contract to read/delete an entry
	ReadOrDelRequest struct {
		Key *datastore.Key `json:"key"`
	}
)

// Validate validates
func (r *CreateRequest) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("url is not set")
	}
	if r.XPATH == "" {
		return fmt.Errorf("xpath is not set")
	}
	if r.Name == "" {
		return fmt.Errorf("name is not set")
	}
	if r.ExpectedPrice == "" {
		return fmt.Errorf("ExpectedPrice is not set")
	}
	r.Options.setDefault()
	return nil
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
func (r *ReadOrDelRequest) Validate() error {
	if r.Key == nil {
		return fmt.Errorf("key is not given")
	}
	return nil
}

// Validate validates
func (r *UpdateRequest) Validate() error {
	if r.Key == nil {
		return fmt.Errorf("key is not given")
	}
	if r.Options == nil {
		return fmt.Errorf("options is not given")
	}
	return nil
}
