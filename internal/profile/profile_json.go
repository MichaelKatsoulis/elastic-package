// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package profile

import (
	"encoding/json"
	"io"
	"os"
	"os/user"
	"time"

	"github.com/pkg/errors"

	"github.com/elastic/go-resource"

	"github.com/elastic/elastic-package/internal/version"
)

// Metadata stores the data associated with a given profile
type Metadata struct {
	Name        string    `json:"name"`
	DateCreated time.Time `json:"date_created"`
	User        string    `json:"user"`
	Version     string    `json:"version"`
	Path        string    `json:"path"`
}

// profileMetadataContent generates the content of the profile.json file.
func profileMetadataContent(applyCtx resource.Context, w io.Writer) error {
	currentUser, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "error fetching current user")
	}

	profileName, found := applyCtx.Fact("profile_name")
	if !found {
		return errors.New("unknown profile name")
	}

	profilePath, found := applyCtx.Fact("profile_path")
	if !found {
		return errors.New("unknown profile path")
	}

	profileData := Metadata{
		profileName,
		time.Now(),
		currentUser.Username,
		version.CommitHash,
		profilePath,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	err = enc.Encode(profileData)
	if err != nil {
		return errors.Wrap(err, "error marshalling json")
	}

	return nil
}

func loadProfileMetadata(path string) (Metadata, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return Metadata{}, errors.Wrap(err, "error reading metadata file")
	}

	metadata := Metadata{}
	err = json.Unmarshal(d, &metadata)
	if err != nil {
		return Metadata{}, errors.Wrapf(err, "error checking profile metadata file %q", path)
	}
	return metadata, nil
}
