
* [OpenDNS_Monitor](https://github.com/yepher/OpenDNS_Monitor)- Source Code and Project Page
* [Releases](https://github.com/yepher/OpenDNS_Monitor/releases) - Pre Built Binary Files

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

## Usage Examples (Linux/ OSX)

For simplicity these examples here depend on the following environment variables to be set. You may specify them directly as an argument if that is preffered.

```bash
export OPENDNS_USERNAME=YOUR_OPENDNS_USERNAME
export OPENDNS_PASSWORD=YOUR_OPENDNS_PASSWORD
export $OPENDNS_NETWORKID=YOUR_OPENDNS_NETWORK_ID
```

This will output the compete list of DNS queries for yesterday:

`./OpenDNS_Monitor -username $OPENDNS_USERNAME -password $OPENDNS_PASSWORD -networkid $OPENDNS_NETWORKID`

This will list the complete list of DNS queries for today:

`./OpenDNS_Monitor -username $OPENDNS_USERNAME -password $OPENDNS_PASSWORD -networkid $OPENDNS_NETWORKID -date today`


This will list questionable DNS queries and report which values they DNS entry matches and how many times it was queried. It will also suppress the full list of DNS queries. If you add the SMTP setting it will email the list to the address your specify. Be careful not to include extra spaces in the `fieldList` values.


```
./OpenDNS_Monitor -username $OPENDNS_USERNAME -password $OPENDNS_PASSWORD -networkid $OPENDNS_NETWORKID \
    -date yesterday  -csv2console=false \
    -fieldList "Blacklisted,Blocked by Category,Blocked as Botnet,Blocked as Malware,Blocked as Phishing,\
    Academic Fraud,Adult Themes,Alcohol,Drugs,Hate/Discrimination,Lingerie/Bikini,Nudity,Pornography,\
    Sexuality,Tobacco"
```

**Sample Output**

```
api.example.com, 2 [Nudity;Pornography;]
counter.example.ru, 2 [Adult Themes;]
porn.example.com, 2 [Adult Themes;Lingerie/Bikini;Nudity;Pornography;Sexuality;]
```


## Username
        OpenDNS Account Username.

## Password
        OpenDNS Account Password.

## CSV to Console (csv2console)

If `-csv2console=true` a comma delimited list of all domains will be written to the console.

## Date

This specifies what data to query DNS queries for. Example `-date today` or `-date 2017-09-27`.

|Value | Description |
|---|---|---
| `today` | Will queries todays DNS entries |
| `yesterday` | Will queries yesterdays DNS entries |
| `YYYY-MM-DD` | Will queries the provided dates DNS entries|


### Field List

This field list comes from the columns of the CSV data that gets returned. By specifying a command delimited list of fields OpenDNS_Monitor will return DNS queries and which of the listed fields the query matches. A given query may trigger one or more of the fields listed. If `-filedList` is not specified queries are not filtered by category.

This is an example usage:

`-fieldList "Blacklisted,Blocked by Category,Blocked as Botnet,Blocked as Malware,Blocked as Phishing,Academic Fraud,Adult Themes,Alcohol,Drugs,Hate/Discrimination,Lingerie/Bikini,Nudity,Pornography,Sexuality,Tobacco"`


| Name  | Description  |
|---|---|
|`Blacklisted`| DNS Queries that were blocked  |
|`Blocked by Category`| The query is blocked by category  |
|`Blocked as Botnet`| The query resolves to an address that represents botnet  |
|`Blocked as Malware`| The query resolves to an address that represents malware  |
|`Blocked as Phishing`| The query resolves to an address that represents phishing  |
|`Resolved by SmartCache`| The query resolves to an address that represents resolved by smartcache  |
|`Academic Fraud`| The query resolves to an address that represents academic fraud  |
|`Adult Themes`| The query resolves to an address that represents adult themes  |
|`Adware`| The query resolves to an address that represents adware  |
|`Alcohol`| The query resolves to an address that represents alcohol  |
|`Anime/Manga/Webcomic`| The query resolves to an address that represents anime, manga and/or webcomics  |
|`Auctions`| The query resolves to an address that represents auctions  |
|`Automotive`| The query resolves to an address that represents automotive  |
|`Blogs`| The query resolves to an address that represents blogs  |
|`Business Services`| The query resolves to an address that represents business services  |
|`Chat`| The query resolves to an address that represents chat  |
|`Classifieds`| The query resolves to an address that represents classifieds  |
|`Dating`| The query resolves to an address that represents dating  |
|`Drugs`| The query resolves to an address that represents drugs  |
|`Ecommerce/Shopping`| The query resolves to an address that represents ecommerce and/or shopping  |
|`Educational Institutions`| The query resolves to an address that represents educational institutions  |
|`File Storage`| The query resolves to an address that represents file storage  |
|`Financial Institutions`| The query resolves to an address that represents financial institutions  |
|`Forums/Message boards`| The query resolves to an address that represents forums and/or message boards  |
|`Gambling`| The query resolves to an address that represents gambling  |
|`Games`| The query resolves to an address that represents games  |
|`German Youth Protection`| The query resolves to an address that represents German youth protection  |
|`Government`| The query resolves to an address that represents goverment  |
|`Hate/Discrimination`| The query resolves to an address that represents hate and/or discrimination  |
|`Health and Fitness`| The query resolves to an address that represents health and fitness  |
|`Humor`| The query resolves to an address that represents humor  |
|`Instant Messaging`| The query resolves to an address that represents instant messaging  |
|`Jobs/Employment`| The query resolves to an address that represents jobs and/or employment  |
|`Lingerie/Bikini`| The query resolves to an address that represents lingerie and or bikinis  |
|`Movies`| The query resolves to an address that represents movies  |
|`Music`| The query resolves to an address that represents music  |
|`News/Media`| The query resolves to an address that represents news and or media  |
|`Non-Profits`| The query resolves to an address that represents non-profits  |
|`Nudity`| The query resolves to an address that represents nudity  |
|`P2P/File sharing`| The query resolves to an address that represents p2p file sharing  |
|`Parked Domains`| The query resolves to an address that represents parked domains  |
|`Photo Sharing`| The query resolves to an address that represents photo sharing  |
|`Podcasts`| The query resolves to an address that represents podcasts  |
|`Politics`| The query resolves to an address that represents politics  |
|`Pornography`| The query resolves to an address that represents pornography  |
|`Portals`| The query resolves to an address that represents a portal  |
|`Proxy/Anonymizer`| The query resolves to an address that represents a proxy  |
|`Radio`| The query resolves to an address that represents radio content  |
|`Religious`| The query resolves to an address that represents religion  |
|`Research/Reference`| The query resolves to an address that represents research and/or references  |
|`Search Engines`| The query resolves to an address that represents a search engine  |
|`Sexuality`| The query resolves to an address that represents sexual content  |
|`Social Networking`| The query resolves to an address that represents social networking  |
|`Software/Technology`|The query resolves to an address that represents software and/or technology   |
|`Sports`| The query resolves to an address that represents sports  |
|`Tasteless`| The query resolves to an address that represents tasteless content  |
|`Television`| The query resolves to an address that represents television  |
|`Tobacco`| The query resolves to an address that represents tobacco  |
|`Travel`| The query resolves to an address that represents travel  |
|`Video Sharing`| The query resolves to an address that represents video sharing  |
|`Visual Search Engines`| The query resolves to a visual search engine  |
|`Weapons`| The query resolves to an address that represents weapons  |
|`Web Spam`| The query represents Web Spam  |
|`Webmail`| The query represents webmail  |


## Log Level

| Level | Description |
|---|---|
| `0` | No logging |
| `1` | Only log errors |
| `2` | Log Errors and Warnings |
| `3` (default) | Log Errors, Warnings, and Info |
| `4` | Log everything (verbose) |


## Network ID

OpenDNS Network ID. You can find this ID while your are on the OpenDNS stats page. It is the number found in the URL:

`https://dashboard.opendns.com/stats/12345678/start/`

In that URL the Network ID is `12345678`

## Output File-outputfile string

Where to write output CSV (default "/tmp/dnsoutput.csv"). Windows users must specify this value.

        
## SMTP Settings

By specifying SMTP settings OpenDNS_Monitor will email results that match your filter list. It maybe helpful to run this in a Crontab at midnight each night like this:

`crontab -e`

Then ad this line:

`00 00 * * * /usr/local/sbin/OpenDNS_Monitor [ALL THE OTHER ARGUMENTS YOU NEED]`


### SMTP From

Email from address.

### SMTP Host

Email server hostname port example `smtp.example.com:587`.

### SMTP Password

Email server password.
        
### SMTP To

Email to address.
        
### SMTP Username

Email server username.
        

