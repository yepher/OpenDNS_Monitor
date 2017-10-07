/**
* Written by Chris Wilson <github@yepher.com>
*
* This code is based on https://github.com/rcrowley/opendns-fetchstats
*
**/
package main

import (
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

// APIKey is used for the OpenDNS http queries
const APIKey = "F5DF5551AB0325FDBD6969F6920B33ED"

// APIBaseURL this is the common part of the url for all API requests
const APIBaseURL = "https://api.opendns.com/v1/"

// TokenResponse is the structure of the Token response API
type TokenResponse struct {
	Status   string `json:"status"`
	Response struct {
		Token string `json:"token"`
	} `json:"response"`
	Error        int    `json:"error"`
	ErrorMessage string `json:"error_message"`
}

// NetworkObject defines the content of the NetworkResponse
type NetworkObject struct {
	Dynamic   bool   `json:"dynamic"`
	Label     string `json:"label"`
	IPAddress string `json:"ip_address"`
}

// NetworksResponse carries the response for the networks API
type NetworksResponse struct {
	Status       string                   `json:"status"`
	Response     map[string]NetworkObject `json:"response"`
	Error        int                      `json:"error"`
	ErrorMessage string                   `json:"error_message"`
}

var (
	// Trace the trace logger
	Trace *log.Logger

	// Info the info logger
	Info *log.Logger

	// Warning the warning logger
	Warning *log.Logger

	// Error the error logger
	Error *log.Logger
)

// Init logging system
// See Logging setup https://www.goinggo.net/2013/11/using-log-package-in-go.html
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func initLogLeve(level int) {
	switch level {
	case 0:
		Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
	case 1:
		Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	case 2:
		Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stderr)
	case 3:
		Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	case 4:
		Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	default:
		Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	}
}

func main() {
	const DateFormat = "2006-01-02"

	// Process Command-line arguments
	usernamePtr := flag.String("username", "", "OpenDNS Account Username.")
	passwordPtr := flag.String("password", "", "OpenDNS Account Password.")
	networkIDPtr := flag.String("networkid", "", "The networkID to report on. If not specified a list of available networks will be listed.")
	outputFilePtr := flag.String("outputfile", "/tmp/dnsoutput.csv", "Where to write output csv")
	csv2console := flag.Bool("csv2console", true, "Write CSV data to the console")
	logLevelPtr := flag.Int("logLevel", 3, "0 - no logging, 1 - error, 2 - warn, 3 - info, 4 - verbose")
	fieldListPtr := flag.String("fieldList", "", "List of fields to report if set")
	showFilteredPtr := flag.Bool("showFiltered", true, "Write filtered DNS entries to the console")

	// SMTP Server Settings
	smtpUsername := flag.String("smtpUsername", "", "Email server username.")
	smtpPassword := flag.String("smtpPassword", "", "Email server password.")
	smtpHost := flag.String("smtpHost", "", "Email server hostname port example `smtp.example.com:587`.")
	smtpFrom := flag.String("smtpFrom", "", "Email from address.")
	smtpTo := flag.String("smtpTo", "", "Email to address.")

	// Default date to yesterday
	datePtr := flag.String("date", "yesterday", "Date to get results for. Valid values YYYY-MM-DD, yesterday, today")
	flag.Parse()

	initLogLeve(*logLevelPtr)

	if *usernamePtr == "" {
		Error.Println("Username is a required field\n\n ")
		flag.PrintDefaults()
		os.Exit(1)
	} else if *passwordPtr == "" {
		Error.Println("Password is a required field\n\n ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if strings.EqualFold(*datePtr, "today") {
		currentTime := time.Now().UTC()
		result := currentTime.Format(DateFormat)
		datePtr = &result
	} else if strings.EqualFold(*datePtr, "yesterday") {
		currentTime := time.Now().AddDate(0, 0, -1).UTC()
		result := currentTime.Format(DateFormat)
		datePtr = &result
	} else if *datePtr == "" {
		currentTime := time.Now().AddDate(0, 0, -1).UTC()
		result := currentTime.Format(DateFormat)
		datePtr = &result
	}

	// Basic setup

	loginURL := "https://login.opendns.com/?source=dashboard"
	csvURL := "https://dashboard.opendns.com"

	cookieJar, _ := cookiejar.New(nil)
	Trace.Printf("Cookies: %v\n", cookieJar)

	// Get the signin page's form token
	token := getFormToken(loginURL, cookieJar)
	Trace.Printf("Cookies: %v\n", cookieJar)
	Trace.Println("Token: " + token)

	// Sign into OpenDNS
	Trace.Println("Username: " + *usernamePtr)
	Trace.Println("Password: " + *passwordPtr)
	Trace.Println("Token: " + token)

	if !signIn(loginURL, token, *usernamePtr, *passwordPtr, cookieJar) {
		Error.Println("Login failed!")
		os.Exit(2)
	}
	Trace.Printf("Cookies: %v\n", cookieJar)
	Trace.Println("Login Success")

	if *networkIDPtr == "" {
		emptyCookieJar, _ := cookiejar.New(nil)
		signInResult := apiSignIn(*usernamePtr, *passwordPtr, emptyCookieJar)
		Trace.Println(signInResult)
		if signInResult.ErrorMessage != "" {
			Error.Printf("ERROR: Sign In failed - code:(%d) message:%s\n", signInResult.Error, signInResult.ErrorMessage)
			return
		}

		networksJSON := listNetworks(signInResult.Response.Token, emptyCookieJar)
		Trace.Println(networksJSON)
		if networksJSON.ErrorMessage != "" {

			Error.Printf("ERROR: Sign In failed - code:(%d) message:%s\n", networksJSON.Error, networksJSON.ErrorMessage)
			return
		}

		Info.Printf("Network ID not provdide so listing available networks:\n")
		for k, v := range networksJSON.Response {
			//fmt.Printf("key[%s] value[%v]\n", k, v)
			fmt.Printf("NetworkId - %s\n", k)
			isDynamic := "false"
			if v.Dynamic {
				isDynamic = "true"
			}
			fmt.Printf("\tLabel: %s\n", v.Label)
			fmt.Printf("\tDynamic: %s\n", isDynamic)
			fmt.Printf("\tIP Address: %s\n\n", v.IPAddress)
		}

		return
	}

	csvResult := fetchTopDomains(csvURL, *networkIDPtr, *datePtr, cookieJar)

	if *csv2console == true {
		fmt.Println(csvResult)
	}

	Trace.Println("Filter fields: " + *fieldListPtr)
	if len(*fieldListPtr) > 0 {
		result := processCSV(csvResult, fieldListPtr)

		if *showFilteredPtr && len(result) > 0 {
			fmt.Println(result)
		}

		if *smtpUsername != "" && len(result) > 0 {
			sendResultAsEmail(*smtpUsername, *smtpPassword, *smtpHost, *smtpFrom, *smtpTo, result)
		}
	}

	// write the whole body at once
	err := ioutil.WriteFile(*outputFilePtr, []byte(csvResult), 0644)
	if err != nil {
		Error.Println(err)
		check(err)
	}
}

// Get the signin page's form token
func getFormToken(loginURL string, cookieJar *cookiejar.Jar) string {

	client := &http.Client{
		Jar: cookieJar,
	}

	response, err := client.Get(loginURL)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		formTokenLine := findLine(string(contents), "formtoken")

		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		elements := strings.Split(formTokenLine, "value=\"")
		token := strings.Split(elements[1], "\"")[0]
		Trace.Printf("Found Token: %s\n", token)

		Trace.Printf("Inner Cookies: %v\n", cookieJar)
		return token
	}

	return ""
}

// Sign into OpenDNS
func signIn(loginURL string, formToken string, username string, password string, cookieJar *cookiejar.Jar) bool {

	client := &http.Client{
		Jar: cookieJar,
	}

	body := strings.NewReader("formtoken=" + formToken + "&username=" + username + "&password=" + password + "&sign_in_submit=foo")
	req, err := http.NewRequest("POST", os.ExpandEnv(loginURL), body)
	if err != nil {
		Error.Printf("%s\n", err)
		os.Exit(2)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		Error.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		if err != nil {
			Error.Printf("%s\n", err)
			os.Exit(2)
		} else {
			loggedIn := findLine(bodyString, "Logging you in")
			Trace.Println("Logged In: " + loggedIn)
			return len(loggedIn) > 0
		}
	} else {
		Error.Printf("Unexpected response code: %d\n", resp.StatusCode)
		os.Exit(2)
	}

	return false
}

// Fetch pages of Top Domains
func fetchTopDomains(url string, networkID string, date string, cookieJar *cookiejar.Jar) string {
	var buffer bytes.Buffer

	for page := 1; page < 100; page++ {
		domainURL := url + "/stats/" + networkID + "/topdomains/" + date + "/page" + strconv.Itoa(page) + ".csv"

		response := doGetRequest(domainURL, "", cookieJar)
		lines := strings.Split(response, "\n")
		Trace.Println("Lines: " + strconv.Itoa(len(lines)))
		if len(lines) <= 2 {
			return buffer.String()
		}

		for i := 0; i < len(lines); i++ {
			if page != 1 && i == 0 {
				// exclude header line of results after first page
			} else {
				if len(lines[i]) > 0 {
					buffer.WriteString(lines[i] + "\n")
				}
			}
		}

	}

	return buffer.String()
}

func doGetRequest(url string, bodyString string, cookieJar *cookiejar.Jar) string {
	Trace.Println("URL: " + url)
	client := &http.Client{
		Jar: cookieJar,
	}

	body := strings.NewReader(bodyString)

	method := "GET"
	if len(bodyString) > 0 {
		method = "POST"
	}

	req, err := http.NewRequest(method, os.ExpandEnv(url), body)
	if err != nil {
		Error.Printf("%s\n", err)
		os.Exit(2)
	}
	if len(bodyString) > 0 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)
	if err != nil {
		Error.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		if err != nil {
			Error.Printf("%s\n", err)
			os.Exit(2)
		} else {
			loggedIn := findLine(bodyString, "Logging you in")
			Trace.Println("Logged In: " + loggedIn)
			return bodyString
		}
	} else {
		Error.Printf("Unexpected response code: %d\n", resp.StatusCode)
		os.Exit(2)
	}

	return ""
}

// https://github.com/yepher/OpenDNS_Monitor/blob/master/Notes.md#get-networks
func apiSignIn(username string, password string, cookieJar *cookiejar.Jar) TokenResponse {
	bodyString := "api_key=" + APIKey + "&method=account_signin&username=" + username + "&password=" + password
	response := doGetRequest(APIBaseURL, bodyString, cookieJar)
	tokenResponse := TokenResponse{}
	json.Unmarshal([]byte(response), &tokenResponse)
	return tokenResponse
}

// https://github.com/yepher/OpenDNS_Monitor/blob/master/Notes.md#get-networks
func listNetworks(token string, cookieJar *cookiejar.Jar) NetworksResponse {
	domainURL := "https://api.opendns.com/v1/"
	bodyString := "api_key=" + APIKey + "&method=networks_get&token=" + token
	response := doGetRequest(domainURL, bodyString, cookieJar)
	networkResponse := NetworksResponse{}
	json.Unmarshal([]byte(response), &networkResponse)
	return networkResponse
}

func processCSV(in string, fieldListPtr *string) string {
	var headerRow []string
	fields := map[int]bool{}

	includeNames := map[string]bool{}

	fieldsList := strings.Split(*fieldListPtr, ",")
	for i := 0; i < len(fieldsList); i++ {
		includeNames[fieldsList[i]] = true
	}

	Trace.Printf("includeNames: %v", includeNames)

	r := csv.NewReader(strings.NewReader(in))
	r.Comma = ','
	r.Comment = '#'

	lineCount := 0
	var result bytes.Buffer
	for {
		record, err := r.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			Error.Println("Error:", err)
			check(err)
		}
		// record is an array of string so is directly printable
		Trace.Println("Record", lineCount, "is", record, "and has", len(record), "fields")

		if lineCount == 0 {
			headerRow = record
			for i := 0; i < len(record); i++ {
				if includeNames[record[i]] == true {
					fields[i] = true
				}
			}
		} else {
			hasHit := false
			var buffer bytes.Buffer
			buffer.WriteString(record[1] + ", " + record[2] + " [")
			for i := 3; i < len(record); i++ {
				if fields[i] && record[i] == "1" {
					hasHit = true
					buffer.WriteString(headerRow[i] + ";")
				}
			}

			if hasHit {
				val := buffer.String() + "]\n"
				fmt.Printf(val)
				result.WriteString(val)
			}
		}
		lineCount++
	}

	return result.String()
}

func sendResultAsEmail(smtpUsername string, smtpPassword string, smtpServer string, smtpFrom string, smtpTo string, body string) {
	// This code is based on https://gist.github.com/jim3ma/b5c9edeac77ac92157f8f8affa290f45
	Trace.Println("SMTP: Connect to: " + smtpServer)

	subject := "OpenDNS Monitor Alert"

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = smtpFrom
	headers["To"] = smtpTo
	headers["Subject"] = subject

	Trace.Printf("SMTP: %v", headers)

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body + "\r\n\r\n----\r\n Sent from OpenDNS_Monitor - https://github.com/yepher/OpenDNS_Monitor"

	servername := smtpServer

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	c, err := smtp.Dial(servername)
	if err != nil {
		Error.Println(err)
		log.Panic(err)
	}

	c.StartTLS(tlsconfig)

	if err = c.Auth(auth); err != nil {
		Error.Println(err)
		log.Panic(err)
	}

	if err = c.Mail(smtpFrom); err != nil {
		Error.Println(err)
		log.Panic(err)
	}

	if err = c.Rcpt(smtpTo); err != nil {
		Error.Println(err)
		log.Panic(err)
	}

	w, err := c.Data()
	if err != nil {
		Error.Println(err)
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		Error.Println(err)
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		Error.Println(err)
		log.Panic(err)
	}

	c.Quit()

}

////////////////////
// Utility Methods
////////////////////

// Percents escapes hex values that represent the input string
func encodeString(input string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(input); i++ {
		buffer.WriteString(fmt.Sprintf("%%%.2x", input[i]))
	}

	return buffer.String()
}

// Returns the first line that contains the given string
func findLine(bodyString string, value string) string {
	lines := strings.Split(bodyString, "\n")
	lineIndex := 0
	for _, line := range lines {
		Trace.Println("[", lineIndex, "]\t", line)
		lineIndex++
		if strings.Contains(line, value) {
			return line
		}
	}
	return ""
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
