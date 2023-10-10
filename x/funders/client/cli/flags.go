package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagMoniker     = "moniker"
	FlagIdentity    = "identity"
	FlagWebsite     = "website"
	FlagContact     = "contact"
	FlagDescription = "description"
)

func flagSetFunderCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagMoniker, "", "The funder's name")
	fs.String(FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fs.String(FlagWebsite, "", "The funder's (optional) website")
	fs.String(FlagContact, "", "The funder's (optional) security contact email")
	fs.String(FlagDescription, "", "The funder's (optional) description")

	return fs
}
