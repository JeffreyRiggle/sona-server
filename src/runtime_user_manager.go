package main

type RuntimeUserManager struct {
	Users map[int]*User
	Passwords map[int]string
}


func (manager RuntimeUserManager) AddUser(user *User) bool {
	var id = len(manager.Users)
	user.Id = id
	manager.Users[id] = user

	return true
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