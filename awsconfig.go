package saml2aws

import (
	"errors"
	"os"
	"path/filepath"

	ini "gopkg.in/ini.v1"
)

var (
	// ErrCredentialsHomeNotFound returned when a user home directory can't be located.
	ErrCredentialsHomeNotFound = errors.New("user home directory not found")

	// ErrCredentialsFileNotFound returned when the required aws credentials file doesn't exist.
	ErrCredentialsFileNotFound = errors.New("aws credentials file not found")
)

// CredentialsProvider loads aws credentials file
type CredentialsProvider struct {
	Filename string
	Profile  string
}

// NewSharedCredentials helper to create the credentials provider
func NewSharedCredentials(profile string) *CredentialsProvider {
	return &CredentialsProvider{
		Profile: profile,
	}
}

// Exists verify that the credentials file exists
func (p *CredentialsProvider) Exists() error {
	filename, err := p.filename()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return ErrCredentialsFileNotFound
		}
		return err
	}

	return nil
}

// Save persist the credentials
func (p *CredentialsProvider) Save(id, secret, token string) error {
	filename, err := p.filename()
	if err != nil {
		return err
	}

	return saveProfile(filename, p.Profile, id, secret, token)
}

func (p *CredentialsProvider) filename() (string, error) {
	if p.Filename == "" {
		if p.Filename = os.Getenv("AWS_SHARED_CREDENTIALS_FILE"); p.Filename != "" {
			return p.Filename, nil
		}

		homeDir := os.Getenv("HOME") // *nix
		if homeDir == "" {           // Windows
			homeDir = os.Getenv("USERPROFILE")
		}
		if homeDir == "" {
			return "", ErrCredentialsHomeNotFound
		}

		p.Filename = filepath.Join(homeDir, ".aws", "credentials")
	}

	return p.Filename, nil
}

func saveProfile(filename, profile, id, secret, token string) error {
	config, err := ini.Load(filename)
	if err != nil {
		return err
	}
	iniProfile, err := config.NewSection(profile)
	if err != nil {
		return err
	}

	_, err = iniProfile.NewKey("aws_access_key_id", id)
	if err != nil {
		return err
	}

	_, err = iniProfile.NewKey("aws_secret_access_key", secret)
	if err != nil {
		return err
	}

	_, err = iniProfile.NewKey("aws_session_token", token)
	if err != nil {
		return err
	}

	return config.SaveTo(filename)
}
