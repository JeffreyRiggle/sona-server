# API documentation
Sona server is a server that uses a REST API to manage incidents. Below is the exposed API.

| Method | Url                                             | Description                             |
|--------|-------------------------------------------------|-----------------------------------------|
| POST   | /sona/v1/incidents                              | Creates an incident.                    |
| PUT    | /sona/v1/incidents/{incidentId}                 | Updates an incident.                    |
| GET    | sona/v1/incidents/{incidentId}/attachments      | Gets an incidents attachments.          |
| POST   | /sona/v1/incidents/{incidentId}/attachment      | Uploads an attachment to an incident.   |
| GET    | /sona/v1/incidents/{incidentId}/attachment/{attachmentId} | Downloads an attachment.                |
| DELETE | /sona/v1/incidents/{incidentId}/attachment/{attachmentId} | Deletes an attachment from an incident. |
| GET    | /sona/v1/incidents                              | Gets incidents.                         |
| GET    | /sona/v1/incidents/{incidentId}                 | Gets an incident.                       |

## Creating in incident

> POST /sona/v1/incidents

### Body
| Property    | type                | Description                                  | Required |
|-------------|---------------------|----------------------------------------------|----------|
| Type        | string              | The type of incident                         | false    |
| Id          | number              | The id of the incident                       | false    |
| Description | string              | The description associated with the incident | false    |
| Reporter    | string              | The individual that reported the incident.   | true     |
| State       | string              | The state the incident is in                 | false    |
| Attributes  | Map<string, string> | Any additional attributes                    | false    |

## Updating an incident

> PUT /sona/v1/incidents/{incidentId}

### Body
| Property    | type                | Description                                  | Required |
|-------------|---------------------|----------------------------------------------|----------|
| Type        | string              | The type of incident                         | false    |
| Description | string              | The description associated with the incident | false    |
| Reporter    | string              | The individual that reported the incident.   | false    |
| State       | string              | The state the incident is in                 | false    |
| Attributes  | Map<string, string> | Any additional attributes                    | false    |

## Getting incident attachments

> GET sona/v1/incidents/{incidentId}/attachments

## Response

| Property    | type         | Description         |
|-------------|--------------|---------------------|
| Attachments | Attachment[] | List of attachments |

### Attachment

| Property | type   | Description                               |
|----------|--------|-------------------------------------------|
| FileName | string | The name of the file                      |
| Time     | string | UTC value for when the file was attached. |

## Adding an attachment to an incident

> POST sona/v1/{incidentId}/attachment

### Body
multipart/form-data upload file.

## Download an attachment

> GET sona/v1/incidents/{incidentId}/attachments/{attachmentId}

### Response

Attachment content


## Remove an attachment

> DELETE sona/v1/incidents/{incidentId}/attachments/{attachmentId}

## Get all incidents

> GET sona/v1/incidents

### Response

| Property  | type       | Description       |
|-----------|------------|-------------------|
| Incidents | Incident[] | List of incidents |

## Get specific incidents

> GET sona/v1/incidents/{incidentId}

| Property    | type                | Description                                  |
|-------------|---------------------|----------------------------------------------|
| Type        | string              | The type of incident                         |
| Id          | number              | The id of the incident                       |
| Description | string              | The description associated with the incident |
| Reporter    | string              | The individual that reported the incident.   |
| State       | string              | The state the incident is in                 |
| Attributes  | Map<string, string> | Any additional attributes                    |