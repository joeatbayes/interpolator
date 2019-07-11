# \<Name of the API\>

\<API description in Plain English. It should be understandable by non-technical audience. Explain domain-specific words if needed.\>

> Sample CURL


```
<Put sample CURL example here>
```

\<This section should precede request section so that it align properly with right hand side code section in Slate\>

### Server URLs per Environment
* Test = https://test.rooseveltsolutions.com
* PROD = https://api.rooseveltsolutions.com

### URI Patterns


   * {HTTP Verb} {API URI}

### Service Tier

  * {Gateway, Service}

## Request
> \<Sample JSON Request\>


| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| request | body | request | Yes | [rqt-schema](#link to rqt-schema object)|

### Query Parameters
Field | Meaning 
---------- | -------
limit | This could be set to a positive value to indicate maximum number of  records to include in response. It defaults to 20 and could be set to  any value from 5 to 50 (both inclusive).
offset | This could be set to a positive value to indicate number of records to skip. By default, it is 0 meaning don't skip any results.

### Headers

Header Field | Meaning
---------- | -------
client_id | Application specific client identifier.
nonce | One time use number generated for session authentication to prevent replay attack.
id_token | Token returned by the authorization server with user authentication information.
x_csrf_token | Token to prevent CSRF attack.  


## Response
> \<Sample JSON Response\>

> \<HTTP status codes\>

| Name | Located in | Description | Schema |
| ---- | ---------- | ----------- | ---- |
| response | body | description | [[resp-schema](#link to resp-schema object)]|


## Errors

Error Code | Meaning
---------- | -------
400 | Bad Request -- Your request is invalid.
401 | Unauthorized -- Your token is wrong.
403 | Forbidden -- You need correct authorizations.
404 | Not Found -- The specified resource was not found.
405 | Method Not Allowed -- You tried to access with an invalid method.
406 | Not Acceptable -- You requested a format that isn't json.
410 | Gone -- The requested resource has been removed from our servers.
500 | Internal Server Error -- We had a problem with our server. Try again later.
503 | Service Unavailable -- We're temporarily offline for maintenance. Please try again later.


## Object Details

### \<Schema Object1 details in tabular format\>  
| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |

### \<Schema Object2 details in tabular format\>  
| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |

### \<Schema Object3 details in tabular format\>  
| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
