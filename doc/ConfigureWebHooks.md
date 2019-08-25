# Web Hooks
Sona server allows you to configure webhooks. These webhooks can run at different times to allow you more automation potential. Web Hooks also support substitution so you can substitute in relevant data.

Web hooks can be broken down into 3 different stages.

1. When an incident is created.
2. When an incident is updated.
3. When an attachment is added to an incident.

## Simple example
The configuration is broken down into three sections, one for each different hook type.

A hook has a couple parameters.

* Method - The method to use (PUT/POST/DELETE/etc)
* Url - The url of your service
* Body - The content to send to the API.

This is a simple example configuration.

```json
{
    "webhooks": {
        "addedhooks":
        [
            {
		        "method": "PUT",
                "url": "https://mysite.com/my/api",
                "body":
		        {
		            "Test": "value"
		        }
            }
        ],
        "updatedhooks": 
        [
            {
		        "method": "PUT",
                "url": "https://mysite.com/my/api",
                "body":
		        {
		            "Test": "value"
		        }
            }
        ],
        "attachedhooks": 
        [
            {
		        "method": "PUT",
                "url": "https://mysite.com/my/api",
                "body":
		        {
		            "Test": "value"
		        }
            }
        ]
    }
}
```

## Substitution
A hook without context might not be very helpful so webhooks support substitution. When a webhook uses substitution values from the incident will be pulled into the request at runtime. You can provided substitution values by using `{{}}` notation.

Below is an example of a substituted webhook.
```json
"webhooks": {
        "addedhooks":
        [
            {
		        "method": "POST",
                "url": "http://mysite.com/email/send",
                "body":
		        {
		            "items":
		            [
		                {"key": "incident", "value": "id", "substitute": true},
		                {"key": "subject", "value": "Incident Created", "substitute": false},
                        {"key": "body", "value": "New Incident Created by {{reporter}} with description {{description}}.", "substitute": true},
                        {"key": "to", "value": "myemail@email.com", "substitute": false}
		            ]
		        }
            }
        ]
    }
```