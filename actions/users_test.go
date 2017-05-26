package actions_test

import (
	"fmt"

	"github.com/as27/buffaloresourcetest/models"
)

// Following naming logic is implemented in Buffalo:
// Model: Singular (User)
// DB Table: Plural (Users)
// Resource: Plural (Users)
// Path: Plural (/users)
// View Template Folder: Plural (/templates/users/)
// While generation following props where added:
// FirstName
// LastName
// Email
func (as *ActionSuite) Test_UsersResource_List() {
	users := models.Users{
		{
			FirstName: "A string for FirstName",
			LastName:  "A string for LastName",
			Email:     "A string for Email",
		},
		{
			FirstName: "Another string for FirstName",
			LastName:  "Another string for LastName",
			Email:     "Another string for Email",
		},
	}
	for _, t := range users {
		err := as.DB.Create(&t)
		as.NoError(err)
	}
	res := as.HTML("/users").Get()
	body := res.Body.String()
	for _, t := range users {
		as.Contains(body, fmt.Sprintf("%s", t.FirstName))
	}
}

func (as *ActionSuite) Test_UsersResource_New() {
	res := as.HTML("/users/new").Get()
	as.Contains(res.Body.String(), "<h1>New <no value></h1>")
}

func (as *ActionSuite) Test_UsersResource_Create() {
	user := &models.User{
		FirstName: "A value for FirstName",
		LastName:  "A value for LastName",
		Email:     "A value for Email",
	}
	res := as.HTML("/users").Post(user)
	as.Equal(301, res.Code)
	as.Equal("/users", res.Location())

	err := as.DB.First(user)
	as.NoError(err)
	as.NotZero(user.ID)
	as.NotZero(user.CreatedAt)
	as.Equal("A value for FirstName", user.FirstName)
	as.Equal("A value for LastName", user.LastName)
	as.Equal("A value for Email", user.Email)
}

func (as *ActionSuite) Test_UsersResource_Create_Errors() {
	user := &models.User{}
	res := as.HTML("/users").Post(user)
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "FirstName can not be blank.")
	as.Contains(res.Body.String(), "LastName can not be blank.")
	as.Contains(res.Body.String(), "Email can not be blank.")
	c, err := as.DB.Count(user)
	as.NoError(err)
	as.Equal(0, c)
}

func (as *ActionSuite) Test_UsersResource_Update() {
	user := &models.User{
		FirstName: "A value for FirstName",
		LastName:  "A value for LastName",
		Email:     "A value for Email",
	}
	verrs, err := as.DB.ValidateAndCreate(user)
	as.NoError(err)
	as.False(verrs.HasAny())
	res := as.HTML("/users/%s", user.ID).Put(&models.User{
		ID:        user.ID,
		FirstName: "A value for FirstName",
		LastName:  "A value for LastName",
		Email:     "A value for Email",
	})
	as.Equal(200, res.Code)

	err = as.DB.Reload(user)
	as.NoError(err)
	as.Equal("A value for FirstName", user.FirstName)
	as.Equal("A value for LastName", user.LastName)
	as.Equal("A value for Email", user.Email)
}
