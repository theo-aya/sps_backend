package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/Qwerci/sps_backend/database"
	"github.com/Qwerci/sps_backend/models"
	"gorm.io/gorm"
	"net/http"
)

// Upgrade websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Counter to assign unique IDs to connecting users
var userIDCounter uint = 1

// handle Connection
func HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Assuming that the user is identified in the authentication process
	// Associate the WebSocket connection with the authenticated user
	authenticatedUserID := GetNextUserID()
	// Create a User object and associate it with the WebSocket connection
	user := models.User{
		Model:      gorm.Model{ID: authenticatedUserID},
		Connection: ws,
	}

	// Handle incoming audio messages
	go HandleSpeechData(&user)

	// Handle received audio for the user
	go HandleReceivedAudio(&user)

	// Send initial data (e.g., chat history)
	SendChatHistory(&user)

	// Close the WebSocket connection when the function returns
	defer ws.Close()
}

// Function to get the next unique user ID (for POC purposes)
func GetNextUserID() uint {
	userIDCounter++
	return userIDCounter
}

// handle Speech Data
func HandleSpeechData(user *models.User) {
	for {
		// Read the audio message from the WebSocket connection
		_, audioData, err := user.Connection.ReadMessage()
		if err != nil {
			break
		}

		// Save the speech data to the database
		SaveSpeechData(user, audioData)

		// Send the speech data to other users
		SendSpeechToOtherUsers(user, audioData)
	}
}

// handle Received Speech(Audio)
func HandleReceivedAudio(user *models.User) {
	// Slice to store audio data
	var audioDataStorage [][]byte

	for{
		select {
		case audioData := <-user.ReceivedAudio:
			// Store the audio data
			audioDataStorage = append(audioDataStorage, audioData)

			// Process the received audio data (e.g., play the audio, display notification)
			storeReceivedAudio(user, audioDataStorage)
		}
	}
}

// processReceivedAudio processes the received audio data
func storeReceivedAudio(user *models.User, audioDataStorage [][]byte) {
	for _, audioData := range audioDataStorage {
		// Your database storage logic here
		// For example, use GORM to create a new Message record with the audio data
		message := models.Message{
			SenderID:    user.ID,
			RecipientID: 0, // Specify the recipient's ID,
			AudioData:   audioData,
		}
		database.DB.Create(&message)
	}
}

// save Speech Data
func SaveSpeechData(user *models.User, audioData []byte) {
	message := models.Message{
		SenderID:    user.ID,
		RecipientID: 0, // Specify the recipient's ID,
		AudioData:   audioData,
	}

	// Save the message to the database
	database.DB.Create(&message)
}

// send Speech Data to other users
func SendSpeechToOtherUsers(sender *models.User, audioData []byte) {
	// Get a list of connected users from the database
	var users []models.User
	database.DB.Find(&users)

	// Iterate through connected users and send the audio data
	for _, recipient := range users {
		// Skip sending to the sender
		if recipient.ID == sender.ID {
			continue
		}

		// Check if the recipient is connected
		if recipient.Connection != nil {
			// Send audio data to the recipient's WebSocket connection
			err := recipient.Connection.WriteMessage(websocket.BinaryMessage, audioData)
			if err != nil {
				// Handle error (e.g., log, disconnect user, etc.)
			}
		} else {
			// If the recipient is not connected, store the audio data in the ReceivedAudio channel
			select {
			case recipient.ReceivedAudio <- audioData:
				// Successfully stored audio data for the offline user
			default:
				// Handle the case when the channel is full (optional)
			}
		}
	}
}

// SendChatHistory sends the chat history to the user over the WebSocket connection
func SendChatHistory(user *models.User) {
	// Retrieve chat history for the user from the database
	chatHistory, err := getChatHistoryByUserID(user.Model.ID)
	if err != nil {
		// Handle the error (e.g., log, disconnect user, etc.)
		fmt.Println("Error retrieving chat history:", err)
		return
	}

	// Loop through the chat history and send messages to the user's WebSocket connection
	for _, message := range chatHistory {
		err := user.Connection.WriteMessage(websocket.BinaryMessage, message.AudioData)
		if err != nil {
			// Handle error (e.g., log, disconnect user, etc.)
			fmt.Println("Error sending chat history:", err)
			return
		}
	}
}

// getChatHistoryByUserID retrieves chat history from the database based on the user ID
func getChatHistoryByUserID(userID uint) ([]models.Message, error) {
	// Placeholder logic for retrieving chat history from the database
	// This might involve querying a database using an ORM like GORM

	// Assuming you have a Message model with an AudioData field
	var messages []models.Message
	result := database.DB.Where("recipient_id = ?", userID).Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}
