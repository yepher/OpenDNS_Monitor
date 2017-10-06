# Other Potential Open DNS interfaces

* `kkowalczyk`
* `zweisoft`
* `GDataHTTPFetcher`

* `"echo \"%@\" | openssl enc -bf -d -pass pass:\"NojkPqnbK8vwmaJWVnwUq\" -salt -a`
* `https://api.opendns.com/v1/`
* `myip.opendns.com`
* `%@/nic/update?token=%@&api_key=%@&v=2&hostname=%@`
* `https://updates.opendns.com`
* `api_key=%@&method=network_dynamic_set&token=%@&network_id=%@&setting=on`
* `api_key=%@&method=networks_get&token=%@`
* `api_key=%@&method=account_signin&username=%@&password=%@`
* API Key: `F5DF5551AB0325FDBD6969F6920B33ED`

* ```while [ `ps -p %d > /dev/null; echo $?` -eq 0 ]; do sleep 0.1; done; /usr/bin/open ```



## Get Response Token

This gets a [TOKEN] that will be needed in future requests.


**Request**

```
POST /v1/ HTTP/1.1
Host: api.opendns.com
Accept-Encoding: gzip, deflate
Content-Type: application/x-www-form-urlencoded
Content-Length: 127
Accept-Language: en-us
Accept: */*
Connection: keep-alive
User-Agent: OpenDNS%20Updater/3.0 CFNetwork/811.5.4 Darwin/16.7.0 (x86_64)

api_key=F5DF5551AB0325FDBD6969F6920B33ED&method=account_signin&username=[YOUR_ACCOUNT_USERNAME]&password=[YOUR_ACCOUNT_PASSWORD]
```

**Success Response**

```
{
	"status": "success",
	"response": {
		"token": "8A66D2D192198E8EAF1234A5B67CD890"
	}
}
```

## Get Networks

**Request**

```
POST /v1/ HTTP/1.1
Host: api.opendns.com
Accept-Encoding: gzip, deflate
Content-Type: application/x-www-form-urlencoded
Content-Length: 99
Accept-Language: en-us
Accept: */*
Connection: keep-alive
User-Agent: OpenDNS%20Updater/3.0 CFNetwork/811.5.4 Darwin/16.7.0 (x86_64)

api_key=F5DF5551AB0325FDBD6969F6920B33ED&method=networks_get&token=[TOKEN]

```

**Success Response**

```
{
	"status": "success",
	"response": {
		"12345678": {
			"dynamic": true,
			"label": "NetworkLable",
			"ip_address": "11.222.333.444"
		}
	}
}

```

**Auth Error Response**

```
{
    "status": "failure",
    "error": 1004,
    "error_message": "Authentication required"
}
```

**Bad Token Error Response**

```
{
    "status": "failure",
    "error": 1002,
    "error_message": "Unknown API key"
}
```


## Update


*Request**

```
GET /nic/update?token=[TOKEN]&api_key=F5DF5551AB0325FDBD6969F6920B33ED&v=2&hostname=[NetworkLable] HTTP/1.1
Host: updates.opendns.com
Connection: keep-alive
Accept-Encoding: gzip, deflate
User-Agent: OpenDNS%20Updater/3.0 CFNetwork/811.5.4 Darwin/16.7.0 (x86_64)
Accept-Language: en-us
Accept: */*


```

**Success Response**

```
good 11.222.333.444
```



## Other Info

* `http://www.opendns.com/software/mac/dynip/about/`
* `https://www.opendns.com/dashboard`
* `http://www.opendns.com/software/mac/dynip/ip-taken/`
* `http://www.opendns.com/software/mac/dynip/ip-differs/`



