package main

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
)

func decodeMessage(buffer []byte) (*MigrationPayload, error) {
	message := &MigrationPayload{}
	err := proto.Unmarshal(buffer, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func printOTPCodes(otpBuffer []byte) error {
	payload, err := decodeMessage(otpBuffer)
	if err != nil {
		return err
	}

	otpArray := payload.OtpParameters
	for _, otp := range otpArray {
		fmt.Println("Issuer: " + otp.Issuer)
		fmt.Println("Name: " + otp.Name)
		encodedSecret := base32.StdEncoding.EncodeToString(otp.Secret)
		fmt.Println("Secret: " + encodedSecret)
		fmt.Println("-----------------------------------")
	}
	return nil
}

func main() {
	files, err := os.ReadDir(".")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			content, err := os.ReadFile(file.Name())
			if err != nil {
				log.Printf("Error reading %s: %v\n", file.Name(), err)
				continue
			}

			otpURL := strings.TrimSpace(string(content))
			parsedURL, err := url.Parse(otpURL)
			if err != nil {
				log.Printf("Error parsing OTP URL in %s: %v\n", file.Name(), err)
				continue
			}

			dataParam := parsedURL.Query().Get("data")
			otpBuffer, err := base64.StdEncoding.DecodeString(dataParam)
			if err != nil {
				log.Printf("Error decoding data parameter in %s: %v\n", file.Name(), err)
				continue
			}

			fmt.Printf("\nOTP URL in %s:\n", file.Name())
			fmt.Println("********************************")
			if err := printOTPCodes(otpBuffer); err != nil {
				log.Printf("Error processing OTP URL in %s: %v\n", file.Name(), err)
			}
		}
	}
}
