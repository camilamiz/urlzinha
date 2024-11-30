package urlzinha

import (
	"fmt"
	"net/http"
)

type GetUrlHandler struct {
}

// type GetUrlResponse struct {
// 	Url string `json:"url"`
// }

type GetUrlRequest struct {
	ShortUrl string `json:"short_url"`
}

func (h *GetUrlHandler) Handle(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.URL.Query().Get("short_url")


	// var requestUrl = r.Vars(r)["short_url"]

	fmt.Println("Request URL:", requestUrl)
	// storeUrl(body.Url, shortUrl)
	// w.Write([]byte(shortUrl))
}

// func checkExistingShortUrl(shortUrl string) string {
// 	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017"))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	defer client.Disconnect(ctx)

// 	collection := client.Database("admin").Collection("urls")
// 	var result bson.M

// 	err = collection.FindOne(context.Background(), bson.M{"short_url": shortUrl}).Decode(&result)
// 	if err != nil {
//     fmt.Println("Non existing short url, let's use it!")
// 		return ""
// 	}

// 	fmt.Println("Found existing URL:", result)
// 	return result["short_url"].(string)
// }
