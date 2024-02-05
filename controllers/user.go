package controllers

import (
	"github.com/Qwerci/sps_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/Qwerci/sps_backend/database"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var validate = validator.New()

func RegisterUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to read the request body"})
		return
	}

	validationErr := validate.Struct(newUser)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	// Encrpt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to encrypt Password"})
		return
	}
	// Validate and hash the password (you may need to add password field in the newUser model)
	// ...
	newUser.Password =  string(hashedPassword)
	// For POC, a simple registration without authentication
	// Store user details in the database
	result := database.DB.Create(&newUser)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}


func LoginUser(c *gin.Context) {
	var loginData struct {
		PhoneNumber string `json:"phone_number"`
		Password	string `json:"password"`
	}

	c.BindJSON(&loginData)

	// Validate login credentials (you may need to compare with hashed password)
	validationErr := validate.Struct(loginData)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("phone_number = ?", loginData.PhoneNumber).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"Error":"user or password does not exist"})
		return
	}

	// Compare hashed password from the database with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"Error":"user or password does not exist"})
		return// Passwords do not match
	}

	// For POC, a simple login without authentication
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "data": existingUser})
}



// import contacts

func ImportContacts(c *gin.Context) {
	var importedContacts []models.Contact
	c.BindJSON(&importedContacts)

	// Assuming that the user is identified in the authentication process
	// You may need to associate the contacts with the authenticated user
	authenticatedUserID := uint(1)  // Replace with your authentication logic

	// Validate and associate the contacts with the authenticated user
	for _, contact := range importedContacts {
		// Check if a user with the same phone number exists
		existingUser := models.User{}
		result := database.DB.Where("phone_number = ?", contact.PhoneNumber).First(&existingUser)

		// If the user exists, associate the contact with the authenticated user
		if result.RowsAffected > 0 {
			contact.UserID = authenticatedUserID
			RegisteredContacts[contact.PhoneNumber] = true // Mark the contact as registered
		}

		// Store the contact in the database regardless of registration status
		database.DB.Create(&contact)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contacts imported successfully"})
}




// listing registed contacts
func GetRegisteredContacts(c *gin.Context) {
	// Assuming that the user is identified in the authentication process
	// You may need to get the authenticated user ID
	authenticatedUserID := uint(1) // Replace with your authentication logic

	// Retrieve registered contacts associated with the authenticated user
	var registeredContactList []models.Contact
	database.DB.Where("user_id = ?", authenticatedUserID).Find(&registeredContactList)

	// Return only phone numbers of registered contacts
	var registeredPhoneNumbers []string
	for _, contact := range registeredContactList {
		registeredPhoneNumbers = append(registeredPhoneNumbers, contact.PhoneNumber)
	}

	c.JSON(http.StatusOK, gin.H{"registered_contacts": registeredPhoneNumbers})
}

