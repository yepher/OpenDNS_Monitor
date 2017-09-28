# OpenDNS Monitor

This tool is meant to simplify tracking inappropriate content access. 



## Usage

```
Usage of ./OpenDNS_Monitor:
  -csv2console
    	Write CSV data to the console (default true)
  -date string
    	Date to get results for. Defaults to yesterday. Valid values YYYY-MM-DD, yesterday, today (default "2017-09-27")
  -fieldList string
    	List of fields to report if set
  -logLevel int
    	0 - no logging, 1 - error, 2 - warn, 3(default) - info, 4 - verbose (default 3)
  -networkid string
    	OpenDNS Network ID.
  -outputfile string
    	Where to write output csv (default "/tmp/dnsoutput.csv")
  -password string
    	OpenDNS Account Password.
  -smtpFrom string
    	Email from address.
  -smtpHost smtp.example.com:587
    	Email server hostname port example smtp.example.com:587.
  -smtpPassword string
    	Email server password.
  -smtpTo string
    	Email to address.
  -smtpUsername string
    	Email server username.
  -username string
    	OpenDNS Account Username.

```