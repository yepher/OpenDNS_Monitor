LOGLEVEL=3
#go build

./OpenDNS_Monitor -username $OPENDNS_USERNAME -password $OPENDNS_PASSWORD -networkid $OPENDNS_NETWORKID -logLevel $LOGLEVEL -csv2console=false -date today -fieldList "Blacklisted,Blocked by Category,Blocked as Botnet,Blocked as Malware,Blocked as Phishing,Academic Fraud,Adult Themes,Alcohol,Drugs,Hate/Discrimination,Lingerie/Bikini,Nudity,Pornography,Sexuality,Tobacco"
