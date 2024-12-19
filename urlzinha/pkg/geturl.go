package urlzinha

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetUrlHandler struct {
}

// type GetUrlResponse struct {
// 	Url string `json:"url"`
// }

type GetUrlRequest struct {
	ShortUrl string `json:"short_url"`
}

type GetShortUrlResponse struct {
	ShortUrl  string `json:"short_url"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

func (h *GetUrlHandler) Handle(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]

	fmt.Println("Short URL:", shortUrl)
	existingUrl := getShortUrl(shortUrl)
	if existingUrl != nil {
		fmt.Println("URL already exists.")

		b, err := json.Marshal(&existingUrl)
		if err != nil {
			fmt.Println("Error marshalling response:", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func getShortUrl(shortUrl string) *GetShortUrlResponse {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer client.Disconnect(ctx)

	collection := client.Database("admin").Collection("urls")
	var result bson.M

	err = collection.FindOne(context.Background(), bson.M{"short_url": shortUrl}).Decode(&result)
	if err != nil {
		fmt.Println("Non existing url, let's store it!")
		return nil
	}

	fmt.Println("Found existing URL:", result)
	response := &GetShortUrlResponse{
		ShortUrl:  result["short_url"].(string),
		Url:       result["url"].(string),
		CreatedAt: result["created_at"].(string),
	}

	return response
}
