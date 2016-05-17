package gjwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type keychain struct {
	sync.Mutex
	keys map[string][]byte
}

var (
	mainKeychain    = &keychain{sync.Mutex{}, make(map[string][]byte)}
	flagSet         = flag.NewFlagSet("gjwt", flag.PanicOnError)
	certificatesURL = flagSet.String("certificates_url", "https://www.googleapis.com/oauth2/v1/certs", "Google signing certificates URL")
)

func init() {
	flag.Parse()
	loadPublicKeys(time.Second)
}

func loadPublicKeys(previousDelay time.Duration) {
	if previousDelay > 60.0*time.Second {
		previousDelay = time.Second
	}

	nextRun, err := updateKeys()
	if err != nil {
		nextRun = previousDelay * 2.0
	}

	// got here? no problems, schedule next run at expiry time
	time.AfterFunc(nextRun, func() {
		loadPublicKeys(nextRun)
	})
}

func updateKeys() (time.Duration, error) {
	req, err := http.NewRequest("GET", *certificatesURL, nil)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if res.StatusCode != 200 {
		return 0, errors.New(res.Status)
	}

	ageSeconds := res.Header.Get(http.CanonicalHeaderKey("age"))
	age, err := time.ParseDuration(fmt.Sprintf("%ss", ageSeconds))
	if err != nil {
		return 0, err
	}

	var maxAgeSeconds string
	cacheControl := res.Header[http.CanonicalHeaderKey("cache-control")]
	for _, v := range cacheControl {
		items := strings.Split(v, ",")
		for _, vv := range items {
			maxAgeItems := strings.Split(vv, "=")
			if len(maxAgeItems) == 2 {
				maxAgeSeconds = maxAgeItems[1]
			}
		}
	}

	maxAge, err := time.ParseDuration(fmt.Sprintf("%ss", maxAgeSeconds))
	if err != nil {
		return 0, err
	}

	var expiresIn = maxAge - age
	keysJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var certKeys map[string]string
	if err := json.Unmarshal(keysJson, &certKeys); err != nil {
		return 0, err
	}

	for x := range certKeys {
		// decode certificate
		block, _ := pem.Decode([]byte(certKeys[x]))
		var cert *x509.Certificate
		cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return 0, err
		}

		// marshal pure public key from certificate
		rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
		pubKey, err := x509.MarshalPKIXPublicKey(rsaPublicKey)
		if err != nil {
			return 0, err
		}

		// create pubkey pem representation
		pubKeyPem := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKey,
		})

		// update keychain
		mainKeychain.Lock()
		mainKeychain.keys[x] = pubKeyPem
		mainKeychain.Unlock()
	}

	return expiresIn, nil
}

func Validate(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		mainKeychain.Lock()
		key, ok := mainKeychain.keys[token.Header["kid"].(string)]
		mainKeychain.Unlock()
		if !ok {
			return nil, fmt.Errorf("Cannot find key")
		}
		return key, nil
	})

	if err != nil || !token.Valid {
		return err
	}

	return nil
}
