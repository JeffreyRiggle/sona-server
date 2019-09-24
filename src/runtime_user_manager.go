package main

type RuntimeUserManager struct {
	Users map[int]*User
	Passwords map[int]string
}


func (manager RuntimeUserManager) AddUser(user *AddUser) (bool, User) {
	var cuser, id = manager.convertAddUser(user)
	manager.Users[id] = cuser

	manager.SetUserPassword(id, user.Password)
	return true, *cuser
}

func (manager RuntimeUserManager) convertAddUser(user *AddUser) (*User, int) {
	var retVal User;
	var id = len(manager.Users)

	retVal.Id = id
	retVal.EmailAddress = user.EmailAddress
	retVal.UserName = user.UserName
	retVal.FirstName = user.FirstName
	
	if len(user.LastName) != 0 {
		retVal.LastName = user.LastName
	}

	if len(user.Gender) != 0 {
		retVal.Gender = user.Gender
	}

	return &retVal, id
}

func (manager RuntimeUserManager) GetUser(userId int) (User, bool) {
	if val, ok := manager.Users[userId]; ok {
		return *val, true
	}

	return User{}, false
}

func (manager RuntimeUserManager) UpdateUser(userId int, user *User) bool {
	manager.Users[userId] = user
	return true
}

func (manager RuntimeUserManager) RemoveUser(userId int) bool {
	delete(manager.Users, userId)
	return true
}


func (manager RuntimeUserManager) SetUserPassword(userId int, password string) {
	manager.Passwords[userId] = password
}

func (manager RuntimeUserManager) AuthenticateUser(userId int, password string) bool {
	return password == manager.Passwords[userId]
}