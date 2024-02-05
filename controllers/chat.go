package controllers

import(
	"github.com/Qwerci/sps_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/Qwerci/sps_backend/database"
	"net/http"
)

// Map to store registered phone numbers
var RegisteredContacts map[string]bool 

type ChatRequest struct {
	RecipientPhoneNumber string `json:"recipient_phone_number"`
}

func init() {
	RegisteredContacts = make(map[string]bool)
}


func InitiateChat(c *gin.Context) {
	var chatRequest ChatRequest
	c.BindJSON(&chatRequest)

	// Assuming that the user is identified in the authentication process
	// You may need to get the authenticated user ID
	authenticatedUserID := GetNextUserID() // Replace with your authentication logic

	// Find the authenticated user
	var authenticatedUser models.User
	database.DB.First(&authenticatedUser, authenticatedUserID)

	// Find the recipient user based on the provided phone number
	var recipientUser models.User
	result := database.DB.Where("phone_number = ?", chatRequest.RecipientPhoneNumber).First(&recipientUser)

	// Check if the recipient user exists
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	// Check if the authenticated user and recipient are registered contacts
	if !RegisteredContacts[chatRequest.RecipientPhoneNumber] {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Recipient is not a registered contact"})
		return
	}

	// Check if a chat between the authenticated user and recipient already exists
	var existingChat models.Message
	result = database.DB.Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		authenticatedUserID, recipientUser.ID, recipientUser.ID, authenticatedUserID).
		First(&existingChat)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Chat already exists"})
		return
	}

	// Create a new chat message to initialize the chat
	chatMessage := models.Message{
		SenderID:    authenticatedUserID,
		RecipientID: recipientUser.ID,
		AudioData:   []byte{}, // You can include an initial message if needed
	}

	// Save the chat message to the database
	database.DB.Create(&chatMessage)

	c.JSON(http.StatusOK, gin.H{"message": "Chat initialized successfully"})
}

