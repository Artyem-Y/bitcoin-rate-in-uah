package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func GetEmailList(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var emailList []string
	for scanner.Scan() {
		email := scanner.Text()
		emailList = append(emailList, email)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return emailList, nil
}

func GetBitcoinRate() (float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data map[string]map[string]float64
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	rate, ok := data["bitcoin"]["uah"]
	if !ok {
		return 0, fmt.Errorf("The exchange rate of Bitcoin in UAH was not found")
	}

	return rate, nil
}

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,6}$`)
	return Re.MatchString(email)
}
