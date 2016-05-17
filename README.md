# gjwt

gjwt validates Google JWT tokens using the list of public keys published. This list is kept up-to-date automatically.

## Install

    go get -u github.com/fcvarela/gjwt
    
## Usage
	
    err := gjwt.Validate(tokenHere)
