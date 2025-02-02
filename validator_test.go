package proxy

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type ValidatorTest struct {
	authEmailFile *os.File
	done          chan bool
	updateSeen    bool
}

func NewValidatorTest(t *testing.T) *ValidatorTest {
	vt := &ValidatorTest{}
	var err error
	vt.authEmailFile, err = ioutil.TempFile("", "test_auth_emails_")
	if err != nil {
		t.Fatal("failed to create temp file: " + err.Error())
	}
	vt.done = make(chan bool, 1)
	return vt
}

func (vt *ValidatorTest) TearDown() {
	vt.done <- true
	os.Remove(vt.authEmailFile.Name())
}

func (vt *ValidatorTest) NewValidator(domains []string,
	updated chan<- bool) func(string) bool {
	return newValidatorImpl(domains, vt.authEmailFile.Name(),
		vt.done, func() {
			if vt.updateSeen == false {
				updated <- true
				vt.updateSeen = true
			}
		})
}

// This will close vt.authEmailFile.
func (vt *ValidatorTest) WriteEmails(t *testing.T, emails []string) {
	defer vt.authEmailFile.Close()
	vt.authEmailFile.WriteString(strings.Join(emails, "\n"))
	if err := vt.authEmailFile.Close(); err != nil {
		t.Fatal("failed to close temp file " +
			vt.authEmailFile.Name() + ": " + err.Error())
	}
}

func TestValidatorEmpty(t *testing.T) {
	vt := NewValidatorTest(t)
	defer vt.TearDown()

	vt.WriteEmails(t, []string(nil))
	domains := []string(nil)
	validator := vt.NewValidator(domains, nil)

	if validator("foo.bar@example.com") {
		t.Error("nothing should validate when the email and " +
			"domain lists are empty")
	}
}

func TestValidatorSingleEmail(t *testing.T) {
	vt := NewValidatorTest(t)
	defer vt.TearDown()

	vt.WriteEmails(t, []string{"foo.bar@example.com"})
	domains := []string(nil)
	validator := vt.NewValidator(domains, nil)

	if !validator("foo.bar@example.com") {
		t.Error("email should validate")
	}
	if validator("baz.quux@example.com") {
		t.Error("email from same domain but not in list " +
			"should not validate when domain list is empty")
	}
}

func TestValidatorSingleDomain(t *testing.T) {
	vt := NewValidatorTest(t)
	defer vt.TearDown()

	vt.WriteEmails(t, []string(nil))
	domains := []string{"example.com"}
	validator := vt.NewValidator(domains, nil)

	if !validator("foo.bar@example.com") {
		t.Error("email should validate")
	}
	if !validator("baz.quux@example.com") {
		t.Error("email from same domain should validate")
	}
}

func TestValidatorMultipleEmailsMultipleDomains(t *testing.T) {
	vt := NewValidatorTest(t)
	defer vt.TearDown()

	vt.WriteEmails(t, []string{
		"xyzzy@example.com",
		"plugh@example.com",
	})
	domains := []string{"example0.com", "example1.com"}
	validator := vt.NewValidator(domains, nil)

	if !validator("foo.bar@example0.com") {
		t.Error("email from first domain should validate")
	}
	if !validator("baz.quux@example1.com") {
		t.Error("email from second domain should validate")
	}
	if !validator("xyzzy@example.com") {
		t.Error("first email in list should validate")
	}
	if !validator("plugh@example.com") {
		t.Error("second email in list should validate")
	}
	if validator("xyzzy.plugh@example.com") {
		t.Error("email not in list that matches no domains " +
			"should not validate")
	}
}

func TestValidatorComparisonsAreCaseInsensitive(t *testing.T) {
	vt := NewValidatorTest(t)
	defer vt.TearDown()

	vt.WriteEmails(t, []string{"Foo.Bar@Example.Com"})
	domains := []string{"Frobozz.Com"}
	validator := vt.NewValidator(domains, nil)

	if !validator("foo.bar@example.com") {
		t.Error("loaded email addresses are not lower-cased")
	}
	if !validator("Foo.Bar@Example.Com") {
		t.Error("validated email addresses are not lower-cased")
	}
	if !validator("foo.bar@frobozz.com") {
		t.Error("loaded domains are not lower-cased")
	}
	if !validator("foo.bar@Frobozz.Com") {
		t.Error("validated domains are not lower-cased")
	}
}

func TestValidatorIgnoreSpacesInAuthEmails(t *testing.T) {
	vt := NewValidatorTest(t)
	defer vt.TearDown()

	vt.WriteEmails(t, []string{"   foo.bar@example.com   "})
	domains := []string(nil)
	validator := vt.NewValidator(domains, nil)

	if !validator("foo.bar@example.com") {
		t.Error("email should validate")
	}
}
