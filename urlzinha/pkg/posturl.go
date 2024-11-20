package urlzinha

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jxskiss/base62"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostUrlHandler struct {
}

type PostUrlRequest struct {
	Url string `json:"url"`
}

type GetUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

func (h *PostUrlHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var body PostUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := generateShortUrl(body.Url)
	fmt.Println("Short URL:", shortUrl)

	storeUrl(body.Url, shortUrl)
	w.Write([]byte(shortUrl))
}

func generateShortUrl(url string) string {
	size := 8
	hash := md5.Sum([]byte(url))
	encodedHash := base62.Encode([]byte(hash[:]))
	return string(encodedHash[:size])
}

func storeUrl(url string, shortUrl string) {
	fmt.Println("starting db")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer client.Disconnect(ctx)
	fmt.Println("db started")

	collection := client.Database("admin").Collection("urls")
	_, err = collection.InsertOne(context.Background(), bson.M{"url": url, "short_url": shortUrl})
	fmt.Println("url:", url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("url stored successfully:", url, shortUrl)
}
