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
	ShortUrl  string `json:"short_url"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

func (h *PostUrlHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var body PostUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingUrl := checkExistingUrl(body.Url)
	if existingUrl != nil {
		fmt.Println("URL already exists.")

		b, err := json.Marshal(&existingUrl)
		if err != nil {
			fmt.Println("Error marshalling response:", err)
			return
		}

		w.Write(b)
		return
	}

	shortUrl := generateShortUrl(body.Url)
	checkedNewShortUrl := checkExistingShortUrl(shortUrl)

	for checkedNewShortUrl != "" {
		shortUrl = generateShortUrl(body.Url)
		checkedNewShortUrl = checkExistingShortUrl(shortUrl)
	}

	fmt.Println("Short URL:", shortUrl)

	shortUrl, createdAt := storeUrl(body.Url, shortUrl)
	// w.Write([]byte(shortUrl))

	response := GetUrlResponse{
		ShortUrl:  shortUrl,
		Url:       body.Url,
		CreatedAt: createdAt,
	}
	b, err := json.Marshal(&response)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func generateShortUrl(url string) string {
	size := 8
	hash := md5.Sum([]byte(url))
	encodedHash := base62.Encode([]byte(hash[:]))
	return string(encodedHash[:size])
}

func storeUrl(url string, shortUrl string) (string, string) {
	fmt.Println("starting db")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer client.Disconnect(ctx)
	fmt.Println("db started")

	createdAt := time.Time.String(time.Now())
	collection := client.Database("admin").Collection("urls")
	_, err = collection.InsertOne(context.Background(), bson.M{"url": url, "short_url": shortUrl, "created_at": createdAt})
	fmt.Println("url:", url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("url stored successfully:", url, shortUrl)
	return shortUrl, createdAt
}

func checkExistingUrl(url string) *GetUrlResponse {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer client.Disconnect(ctx)

	collection := client.Database("admin").Collection("urls")
	var result bson.M

	err = collection.FindOne(context.Background(), bson.M{"url": url}).Decode(&result)
	if err != nil {
		fmt.Println("Non existing url, let's store it!")
		return nil
	}

	fmt.Println("Found existing URL:", result)
	response := &GetUrlResponse{
		ShortUrl:  result["short_url"].(string),
		Url:       result["url"].(string),
		CreatedAt: result["created_at"].(string),
	}

	return response
}

func checkExistingShortUrl(shortUrl string) string {
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
		fmt.Println("Non existing short url, let's use it!")
		return ""
	}

	fmt.Println("Found existing URL:", result)
	return result["short_url"].(string)
}
