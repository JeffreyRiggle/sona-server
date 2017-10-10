package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBIncidentManager provides the ability to manage incidents in AWS DynamoDB
// The Region indicates what region the db will exist it.
// The IncidentTable indicates the name of the table to use for incidents.
// The AttachmentTable indicates the name of the table to use for attachments.
type DynamoDBIncidentManager struct {
	Region          *string
	IncidentTable   *string
	AttachmentTable *string
}

// Initialize setups up the DynamoDBIncidentManger.
// This will make sure we are able to connect to the region.
// It will also create the configured tables in that region if they do not already exist.
func (manager DynamoDBIncidentManager) Initialize() {
	logManager.LogPrintln("Initializing DynamoDB manager")

	svc := CreateService(*manager.Region)

	incidentInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(*manager.IncidentTable),
	}

	td, err := svc.DescribeTable(incidentInput)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				manager.createIncidentTable()
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogFatal(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogFatal(aerr.Error())
			}
		} else {
			logManager.LogFatal(err.Error())
		}
	} else {
		logManager.LogPrintf("Found table description %v", td)
	}

	attachmentInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(*manager.AttachmentTable),
	}

	td2, err2 := svc.DescribeTable(attachmentInput)

	if err2 != nil {
		if aerr2, ok := err2.(awserr.Error); ok {
			switch aerr2.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				manager.createAttachmentTable()
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogFatal(dynamodb.ErrCodeInternalServerError, aerr2.Error())
			default:
				logManager.LogFatal(aerr2.Error())
			}
		} else {
			logManager.LogFatal(err2.Error())
		}
	} else {
		logManager.LogPrintf("Found table description %v", td2)
	}
}

func (manager DynamoDBIncidentManager) createIncidentTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("type"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("type"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(*manager.IncidentTable),
	}

	svc := CreateService(*manager.Region)

	result, err := svc.CreateTable(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceInUseException, aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintln(aerr.Error())
			}
		} else {
			logManager.LogPrintln(err.Error())
		}

		return
	}

	logManager.LogPrintf("Table Created %v\n", result)
}

func (manager DynamoDBIncidentManager) createAttachmentTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("filename"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("filename"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(*manager.AttachmentTable),
	}

	svc := CreateService(*manager.Region)

	result, err := svc.CreateTable(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceInUseException, aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintln(aerr.Error())
			}
		} else {
			logManager.LogPrintln(err.Error())
		}

		return
	}

	logManager.LogPrintf("Table Created %v\n", result)
}

// AddIncident will add an incident to the configured DynamoDB incidents table.
// If the incident is unable to be added to dynamodb false will be returned.
func (manager DynamoDBIncidentManager) AddIncident(incident *Incident) bool {
	logManager.LogPrintln("Got add request in dynamodb manager")
	id, passed := manager.getNextId()

	if !passed {
		return false
	}

	incident.Id = id

	av, err := dynamodbattribute.MarshalMap(incident)
	if err != nil {
		logManager.LogPrintf("Unable to marshal incident, %v", err)
		return false
	}

	logManager.LogPrintln(av)

	svc := CreateService(*manager.Region)

	_, err2 := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(*manager.IncidentTable),
		Item:      av,
	})

	if err2 != nil {
		logManager.LogPrintf("Unable to put incident, %v", err2)
		return false
	}

	return true
}

func (manager DynamoDBIncidentManager) getNextId() (int64, bool) {
	incidents, err := manager.getAllIncidents()

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintf("Got generic aws error, %v\n", aerr.Error())
			}
		} else {
			logManager.LogPrintf("Got generic error, %v\n", err.Error())
		}

		return -1, false
	}

	if len(incidents) == 0 {
		logManager.LogPrintf("There are no records in the database creating first one")
		return 0, true
	}

	lastItem := incidents[len(incidents)-1]
	retVal := lastItem.Id

	retVal++
	logManager.LogPrintf("Found last id of %v next id is %v", lastItem.Id, retVal)
	return retVal, true
}

func (manager DynamoDBIncidentManager) getAllIncidents() ([]Incident, error) {
	var incidents []Incident

	svc := CreateService(*manager.Region)

	err := svc.ScanPages(&dynamodb.ScanInput{
		TableName: aws.String(*manager.IncidentTable),
	}, func(page *dynamodb.ScanOutput, last bool) bool {
		incs := []Incident{}

		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &incs)

		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal items, %v", err))
		}

		incidents = append(incidents, incs...)

		return true
	})

	return incidents, err
}

func (manager DynamoDBIncidentManager) getFilteredIncidents(filter *FilterRequest) ([]Incident, error) {
	var incidents []Incident

	svc := CreateService(*manager.Region)

	queryString, names, values := buildAWSFilterString(filter)

	logManager.LogPrintf("Attempting to query dynamodb with %v\n", queryString)

	err := svc.ScanPages(&dynamodb.ScanInput{
		TableName:                 aws.String(*manager.IncidentTable),
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
		FilterExpression:          aws.String(queryString),
	}, func(page *dynamodb.ScanOutput, last bool) bool {
		incs := []Incident{}

		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &incs)

		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal items, %v", err))
		}

		incidents = append(incidents, incs...)

		return true
	})

	return incidents, err
}

func buildAWSFilterString(filter *FilterRequest) (string, map[string]*string, map[string]*dynamodb.AttributeValue) {
	var buffer bytes.Buffer
	nameIter := 0
	attributeNames := make(map[string]*string, 0)
	attributeValues := make(map[string]*dynamodb.AttributeValue)

	for i, filter := range filter.Filters {
		if i != 0 {
			buffer.WriteString("and ")
		}
		for iter, complexFilter := range filter.Filter {
			if iter != 0 {
				buffer.WriteString("and ")
			}

			nIt := strconv.Itoa(nameIter)
			attributeNames["#name"+nIt] = aws.String(strings.ToLower(complexFilter.Property))
			attributeValues[":value"+nIt] = &dynamodb.AttributeValue{
				S: aws.String(complexFilter.Value),
			}

			buffer.WriteString(convertDynamoFilterExpression(complexFilter, "#name"+nIt, ":value"+nIt))
			buffer.WriteString(" ")
			nameIter++
		}
	}

	return buffer.String(), attributeNames, attributeValues
}

func convertDynamoFilterExpression(filter Filter, name string, value string) string {
	if isEqualsComparision(filter) {
		return name + " = " + value
	}

	if isNotEqualsComparision(filter) {
		return name + " <> " + value
	}

	return "contains( " + name + ", " + value + " )"
}

// GetIncident will attempt to get the requested incident out of dynamodb.
// If the attempt fails a empty incident will be returned along with a false.
// If the attempt fails the incident will be returned along with a true.
func (manager DynamoDBIncidentManager) GetIncident(incidentId int) (Incident, bool) {
	logManager.LogPrintln("Got Get request.")
	inc, pass := manager.getIncidentFromDataBase(incidentId)
	return *inc, pass
}

// UpdateIncident will attempt to update an incident in dynamodb.
// If the attempt fails a false will be returned.
func (manager DynamoDBIncidentManager) UpdateIncident(id int, update IncidentUpdate) bool {
	logManager.LogPrintln("Got update request.")
	inc, pass := manager.getIncidentFromDataBase(id)

	if !pass {
		return false
	}

	updated := updateIncident(inc, update)
	if !updated {
		return true
	}

	return manager.updateItemInDataBase(*inc)
}

func (manager DynamoDBIncidentManager) getIncidentFromDataBase(incidentId int) (*Incident, bool) {
	svc := CreateService(*manager.Region)

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"type": {
				S: aws.String("Incident"),
			},
			"id": {
				N: aws.String(strconv.Itoa(incidentId)),
			},
		},
		TableName: aws.String(*manager.IncidentTable),
	}

	result, err := svc.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintln(aerr.Error())
			}
		} else {
			logManager.LogPrintln(err.Error())
		}
		return nil, false
	}

	retVal := Incident{}
	for k, v := range result.Item {
		logManager.LogPrintf("Umarshaling %v", k)

		if k == "type" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.Type = umVal
		}
		if k == "id" {
			var umVal int64
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.Id = umVal
		}
		if k == "description" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.Description = umVal
		}
		if k == "reporter" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.Reporter = umVal
		}
		if k == "state" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.State = umVal
		}
		if k == "attributes" && v != nil {
			var umVal map[string]string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.Attributes = umVal
		}
	}

	if retVal.Attributes == nil {
		logManager.LogPrintln("No attributes found returning empty attributes")
		retVal.Attributes = make(map[string]string, 0)
	}

	return &retVal, true
}

// GetIncidents will attempt to get all incidents out of dynamodb.
// If the attempt fails an empty array and a false will be returned.
// If the attempt passes a array of incidents and a true will be returned.
func (manager DynamoDBIncidentManager) GetIncidents(filter *FilterRequest) ([]Incident, bool) {

	var (
		incidents []Incident
		err       error
	)

	if filter == nil {
		incidents, err = manager.getAllIncidents()
	} else {
		incidents, err = manager.getFilteredIncidents(filter)
	}

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintf("Got generic aws error, %v\n", aerr.Error())
			}
		} else {
			logManager.LogPrintf("Got generic error, %v\n", err.Error())
		}
	}

	return incidents, err == nil
}

func (manager DynamoDBIncidentManager) updateItemInDataBase(incident Incident) bool {
	svc := CreateService(*manager.Region)

	attMap, err := dynamodbattribute.MarshalMap(incident.Attributes)

	if err != nil {
		logManager.LogPrintln(err)
		return false
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#s": aws.String("state"),
			"#d": aws.String("description"),
			"#r": aws.String("reporter"),
			"#a": aws.String("attributes"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(incident.State),
			},
			":d": {
				S: aws.String(incident.Description),
			},
			":r": {
				S: aws.String(incident.Reporter),
			},
			":a": {
				M: attMap,
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"type": {
				S: aws.String("Incident"),
			},
			"id": {
				N: aws.String(strconv.FormatInt(incident.Id, 10)),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(*manager.IncidentTable),
		UpdateExpression: aws.String("SET #s = :s, #d = :d, #r = :r, #a = :a"),
	}

	result, err := svc.UpdateItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				logManager.LogPrintln(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintln(aerr.Error())
			}
		} else {
			logManager.LogPrintln(err.Error())
		}
		return false
	}

	logManager.LogPrintln(result)
	return true
}

// AddAttachment will attempt to add an association between an incident and an attachment.
// If the attempt fails a false will be returned.
func (manager DynamoDBIncidentManager) AddAttachment(incidentId int, attachment Attachment) bool {
	av, err := dynamodbattribute.MarshalMap(attachment)
	if err != nil {
		logManager.LogPrintf("Unable to marshal incident, %v", err)
		return false
	}

	av["id"] = &dynamodb.AttributeValue{
		N: aws.String(strconv.Itoa(incidentId)),
	}

	logManager.LogPrintln(av)

	svc := CreateService(*manager.Region)

	_, err2 := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(*manager.AttachmentTable),
		Item:      av,
	})

	if err2 != nil {
		logManager.LogPrintf("Unable to put attachment on incident, %v", err2)
		return false
	}

	return true
}

// GetAttachments will attempt to find all attachments for a given incident id.
// If the attempt fails an empty array an a false will be returned.
// If the attempt passes a array of attachments associated with the incident and a true will be returned.
func (manager DynamoDBIncidentManager) GetAttachments(incidentId int) ([]Attachment, bool) {
	var attachments []Attachment

	svc := CreateService(*manager.Region)

	var resp, err = svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(*manager.AttachmentTable),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						N: aws.String(strconv.Itoa(incidentId)),
					},
				},
			},
		},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintf("Got generic aws error, %v\n", aerr.Error())
			}
		} else {
			logManager.LogPrintf("Got generic error, %v\n", err.Error())
		}

		return make([]Attachment, 0), false
	}

	dynamodbattribute.UnmarshalListOfMaps(resp.Items, &attachments)
	logManager.LogPrintf("Found %v attachments\n", len(attachments))
	return attachments, true
}

// RemoveAttachment will find and remove an attachment associated with an incident.
func (manager DynamoDBIncidentManager) RemoveAttachment(incidentId int, fileName string) bool {
	svc := CreateService(*manager.Region)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.Itoa(incidentId)),
			},
			"filename": {
				S: aws.String(fileName),
			},
		},
		TableName: aws.String(*manager.AttachmentTable),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				logManager.LogPrintln(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				logManager.LogPrintln(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				logManager.LogPrintln(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				logManager.LogPrintln(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				logManager.LogPrintln(aerr.Error())
			}
		} else {
			logManager.LogPrintln(err.Error())
		}

		return false
	}

	logManager.LogPrintln("Removed attachment from dynamodb.")
	return true
}

// CleanUp will do any required cleanup actions on the incident manager.
func (manager DynamoDBIncidentManager) CleanUp() {
	// No op
}

// CreateService will create a new dynamodb.DynamoDB instance.
func CreateService(region string) *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return dynamodb.New(sess)
}
