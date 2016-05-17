# gjwt

gjwt validates Google JWT tokens using the list of public keys published. This list is kept up-to-date automatically.

## Install

    go get -u github.com/fcvarela/gjwt
    
## Usage
	
    if token, err := gjwt.Validate(tokenHere); err == nil {
    	clientid, ok := token.Claims["aud"].(string)
    }
    
