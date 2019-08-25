# API documentation
Sona server is a server that uses a REST API to manage incidents. Below is the exposed API.

| Method | Url                                             | Description                             |
|--------|-------------------------------------------------|-----------------------------------------|
| POST   | /sona/v1/create                                 | Creates an incident.                    |
| PUT    | /sona/v1/{incidentId}/update                    | Updates an incident.                    |
| GET    | sona/v1/{incidentId}/attachments                | Gets an incidents attachments.          |
| POST   | /sona/v1/{incidentId}/attachment                | Uploads an attachment to an incident.   |
| GET    | /sona/v1/{incidentId}/attachment/{attachmentId} | Downloads an attachment.                |
| DELETE | /sona/v1/{incidentId}/attachment/{attachmentId} | Deletes an attachment from an incident. |
| GET    | /sona/v1/incidents                              | Gets incidents.                         |
| GET    | /sona/v1/incidents/{incidentId}                 | Gets an incident.                       |

## Creating in incident

> POST /sona/v1/create

### Body
