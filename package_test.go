package gjwt

import "testing"

type fixture struct {
	token string
	valid bool
}

var (
	fixtures = []fixture{
		{
			token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjUyNDNiNDI5ZGUxOGY0NTY4NTYwOTMwNDY3NDBlMDU2NjRjNDI5OTYifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdF9oYXNoIjoiVzlfazV6R2ZlVDBPdWZ3M2JLSG1WZyIsImF1ZCI6IjgwMzcwOTkwNjE3MS00anJmb3Yyc2E1dDRrOGlkZDV0cHI5ODRodWg4dXI4Zi5hcHBzLmdvb2dsZXVzZXJjb250ZW50LmNvbSIsInN1YiI6IjEwMzE3MjAzMTEzMzMyMjAyNTI3MCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhenAiOiI4MDM3MDk5MDYxNzEtNGpyZm92MnNhNXQ0azhpZGQ1dHByOTg0aHVoOHVyOGYuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJlbWFpbCI6ImZjdmFyZWxhQGdtYWlsLmNvbSIsImlhdCI6MTQ2MzQ4NTU0NywiZXhwIjoxNDYzNDg5MTQ3fQ.sBgCOnFocf2E0LXTFdqole6ZKYTC-LyvusWavgor4PW_gRXl-8YilPCTuStM8MOdCkOiVBoEGOXI368kg1z4AvyCVVSGd7TykXSGWWerCnkVsc5ZChi3imZheMAob7kO-jmiFRWZIaM-HzRpQBxwP6jrNkPk9VeVHaBxICNk7djJgp51usReNMf7dTKxKipBawZDBRQyi-CazTQBkYSJQxZeaeJ8TC4yyEXsmKD6m1mhjh6fkKlkv39k1C9-C83HgCjRhpaajOPjS5PW96Fbcse9tHaCPnT--PHqoNhH-cP30SOBzT2T7wediMHAe0PszPqruaQ3yT98sHrfCvznrA",
			valid: true,
		},
	}
)

func TestValidate(t *testing.T) {
	for _, f := range fixtures {
		_, err := Validate(f.token)

		// good token, expected bad one
		if err == nil && !f.valid {
			t.Error("Good token, expected bad")
			t.Fail()
		}

		if err != nil && f.valid {
			t.Error("Bad token, expected good")
			t.Fail()
		}
	}
}
