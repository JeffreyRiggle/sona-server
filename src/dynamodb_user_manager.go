package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBUserManager provides the ability to manage incidents in AWS DynamoDB
// The Region indicates what region the db will exist it.
// The UsersTable indicates the name of the table to use for users.
type DynamoDBUserManager struct {
	Region             *string
	Endpoint           *string
	UsersTable         *string
	DefaultPermissions []string
}

// Initialize setups up the DynamoDBUserManager.
// This will make sure we are able to connect to the region.
// It will also create the configured tables in that region if they do not already exist.
func (manager DynamoDBUserManager) Initialize() {
	logManager.LogPrintln("Initializing DynamoDB user manager")

	svc := CreateService(*manager.Region, *manager.Endpoint)

	userInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(*manager.UsersTable),
	}

	td, err := svc.DescribeTable(userInput)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				manager.createUsersTable()
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
}

func (manager DynamoDBUserManager) createUsersTable() {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("emailAddress"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("emailAddress"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(*manager.UsersTable),
	}

	svc := CreateService(*manager.Region, *manager.Endpoint)

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

func (manager DynamoDBUserManager) AddUser(user *AddUser) (bool, User) {
	logManager.LogPrintln("Got add user request in dynamodb manager")
	id, passed := manager.getNextId()

	if !passed {
		return false, User{}
	}

	permissions := make([]string, len(manager.DefaultPermissions))
	usr := User{
		Id:           id,
		UserName:     user.UserName,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		EmailAddress: user.EmailAddress,
		Gender:       user.Gender,
		Permissions:  permissions,
	}

	av, err := dynamodbattribute.MarshalMap(usr)
	if err != nil {
		logManager.LogPrintf("Unable to marshal user, %v", err)
		return false, usr
	}

	logManager.LogPrintln(av)

	svc := CreateService(*manager.Region, *manager.Endpoint)

	_, err2 := svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(*manager.UsersTable),
		Item:      av,
	})

	if err2 != nil {
		logManager.LogPrintf("Unable to put user, %v", err2)
		return false, usr
	}

	manager.SetUserPassword(usr, createPasswordHash(usr, user.Password))

	return true, usr
}

func (manager DynamoDBUserManager) getNextId() (int64, bool) {
	users, err := manager.getAllUsers()

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

	if len(users) == 0 {
		logManager.LogPrintf("There are no records in the database creating first one")
		return 0, true
	}

	var lastIndex int64 = 0

	for _, u := range users {
		if u.Id > lastIndex {
			lastIndex = u.Id
		}
	}

	lastIndex++
	logManager.LogPrintf("Next id is %v", lastIndex)
	return lastIndex, true
}

func (manager DynamoDBUserManager) getAllUsers() ([]User, error) {
	var users []User

	svc := CreateService(*manager.Region, *manager.Endpoint)

	err := svc.ScanPages(&dynamodb.ScanInput{
		TableName: aws.String(*manager.UsersTable),
	}, func(page *dynamodb.ScanOutput, last bool) bool {
		usrs := []User{}

		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &usrs)

		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal items, %v", err))
		}

		users = append(users, usrs...)

		return true
	})

	logManager.LogPrintf("Found users %v", users)

	return users, err
}

func (manager DynamoDBUserManager) GetUser(userId int64) (User, bool) {
	logManager.LogPrintln("Got Get user request.")
	usr, pass := manager.getUserFromDataBase(userId)
	return *usr, pass
}

func (manager DynamoDBUserManager) UpdateUser(userId int64, user *User) bool {
	logManager.LogPrintln("Got update user request.")
	usr, pass := manager.getUserFromDataBase(userId)

	if !pass {
		return false
	}

	updated := updateUser(usr, *user)
	if !updated {
		return true
	}

	return manager.updateUserInDataBase(*usr)
}

func (manager DynamoDBUserManager) getUserFromDataBase(userId int64) (*User, bool) {
	svc := CreateService(*manager.Region, *manager.Endpoint)

	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				N: aws.String(strconv.FormatInt(userId, 10)),
			},
		},
		KeyConditionExpression: aws.String("id = :v1"),
		TableName:              aws.String(*manager.UsersTable),
	}

	result, err := svc.Query(input)

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
		return &User{}, false
	}

	retVal := User{}

	if len(result.Items) == 0 {
		return &User{}, false
	}

	for k, v := range result.Items[0] {
		logManager.LogPrintf("Umarshaling %v", k)

		if k == "emailAddress" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.EmailAddress = umVal
		}
		if k == "id" {
			var umVal int64
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.Id = umVal
		}
		if k == "userName" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.UserName = umVal
		}
		if k == "firstName" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.FirstName = umVal
		}
		if k == "lastName" {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.LastName = umVal
		}
		if k == "gender" && v != nil {
			var umVal string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}
			retVal.Gender = umVal
		}
		if k == "permissions" && v != nil {
			var umVal []string
			err2 := dynamodbattribute.Unmarshal(v, &umVal)

			if err2 != nil {
				logManager.LogPrintln(fmt.Sprintf("failed to unmarshal items, %v", err2))
			}

			retVal.Permissions = umVal
		}
	}

	return &retVal, true
}

func (manager DynamoDBUserManager) updateUserInDataBase(user User) bool {
	svc := CreateService(*manager.Region, *manager.Endpoint)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#u": aws.String("userName"),
			"#f": aws.String("firstName"),
			"#l": aws.String("lastName"),
			"#g": aws.String("gender"),
			"#p": aws.String("permissions"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":u": {
				S: aws.String(user.UserName),
			},
			":f": {
				S: aws.String(user.FirstName),
			},
			":l": {
				S: aws.String(user.LastName),
			},
			":g": {
				S: aws.String(user.Gender),
			},
			":p": {
				SS: aws.StringSlice(user.Permissions),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"emailAddress": {
				S: aws.String(user.EmailAddress),
			},
			"id": {
				N: aws.String(strconv.FormatInt(user.Id, 10)),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(*manager.UsersTable),
		UpdateExpression: aws.String("SET #u = :u, #f = :f, #l = :l, #g = :g, #p = :p"),
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

func (manager DynamoDBUserManager) RemoveUser(userId int64) bool {
	user, pass := manager.getUserFromDataBase(userId)

	if !pass {
		return false
	}

	svc := CreateService(*manager.Region, *manager.Endpoint)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"emailAddress": {
				S: aws.String(user.EmailAddress),
			},
			"id": {
				N: aws.String(strconv.FormatInt(user.Id, 10)),
			},
		},
		TableName: aws.String(*manager.UsersTable),
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

	logManager.LogPrintln("Removed user from dynamodb.")
	return true
}

func (manager DynamoDBUserManager) SetUserPassword(user User, password string) {
	// TODO
}

func (manager DynamoDBUserManager) SetPermissions(userId int64, permissions []string) bool {
	// TODO
	return true
}

func (manager DynamoDBUserManager) AuthenticateUser(user User, password string) (bool, TokenResponse) {
	// TODO
	return true, TokenResponse{}
}

func (manager DynamoDBUserManager) ValidateUser(token string) bool {
	// TODO
	return true
}

// CleanUp will do any required cleanup actions on the incident manager.
func (manager DynamoDBUserManager) CleanUp() {
	// No op
}
