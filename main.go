package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-resty/resty/v2"
)

type token struct {
	Value string `json:"value"`
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type answer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Result    string `json:"result"`
	GithubURL string `json:"githubUrl"`
}

func main() {
	client := resty.New()
	resp, err := client.R().
		SetHeader("content-type", "application/json-patch+json").
		SetBody(authRequest{
			Username: "omega",
			Password: "candidate",
		}).
		Post("http://omega-morse-service.eu-central-1.elasticbeanstalk.com/api/v1/auth")
	if err != nil {
		log.Fatal(err)
	}

	var tokenResponse token
	json.Unmarshal(resp.Body(), &tokenResponse)

	token, _, err := new(jwt.Parser).ParseUnverified(tokenResponse.Value, jwt.MapClaims{})
	if err != nil {
		log.Fatal(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Fatal("Can't convert token's claims to standard claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		log.Fatal("Can't convert expiration to float64")
	}

	text := strings.ToUpper("Vega IT Omega : " + strconv.Itoa(int(exp)))
	morseCode := []string{}

	for _, ch := range text {
		morseCode = append(morseCode, litteral2MorseCode[ch])
	}

	resp, err = client.R().
		SetHeader("content-type", "application/json-patch+json").
		SetAuthToken(tokenResponse.Value).
		SetBody(answer{
			FirstName: "Mrdjan",
			LastName:  "Mrksic",
			Email:     "mrdjanmrksic@gmail.com",
			Result:    strings.Join(morseCode, " "),
			GithubURL: "https://github.com/MrdjanMrksic/morse-code-challenge",
		}).
		Post("http://omega-morse-service.eu-central-1.elasticbeanstalk.com/api/v1/result")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(resp.Body()))

}

var litteral2MorseCode = map[rune]string{
	'A': ".-",
	'B': "-...",
	'C': "-.-.",
	'D': "-..",
	'E': ".",
	'F': "..-.",
	'G': "--.",
	'H': "....",
	'I': "..",
	'J': ".---",
	'K': "-.-",
	'L': ".-..",
	'M': "--",
	'N': "-.",
	'O': "---",
	'P': ".--.",
	'Q': "--.-",
	'R': ".-.",
	'S': "...",
	'T': "-",
	'U': "..-",
	'V': "...-",
	'W': ".--",
	'X': "-..-",
	'Y': ".",
	'Z': "--..",
	'1': ".----",
	'2': "..---",
	'3': "...--",
	'4': "....-",
	'5': ".....",
	'6': "-....",
	'7': "--...",
	'8': "---..",
	'9': "----.",
	'0': "-----",
	':': "---...",
	' ': "/",
}
