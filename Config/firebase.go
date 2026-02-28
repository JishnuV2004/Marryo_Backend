// package config

// import (
// 	"context"
// 	"log"
// 	"os"

// 	firebase "firebase.google.com/go/v4"
// 	"google.golang.org/api/option"
// )

// var FirebaseApp *firebase.App

// func InitFirebase() (*firebase.App, error){


// 	credPath := os.Getenv("FIREBASE_CREDENTIALS")
// 	if credPath == "" {
// 		credPath = "firebase-service-account.json"
// 	}

// 	projectID := os.Getenv("FIREBASE_PROJECT_ID")
// 	if projectID == "" {
// 		log.Fatal("FIREBASE_PROJECT_ID is not set")
// 	}

// 	opt := option.WithCredentialsFile(credPath)

// 	app, err := firebase.NewApp(context.Background(), &firebase.Config{
// 		ProjectID: projectID,
// 	}, opt)
// 	if err != nil {
// 		log.Fatal("Firebase init failed:", err)
// 	}

// 	return app, nil
// }
package config