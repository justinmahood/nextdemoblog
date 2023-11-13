package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Post struct {
	Title  string    `firestore:"title" json:"title"`
  Id     string    `firestore:ID", omitempty" json:"ID, omitempty"`
	Author string    `firestore:"author" json:"author"`
	Date   time.Time `firestore:"date" json:"date"`
	Body   string    `firestore:"body" json:"body"`
}

var t = template.Must(template.ParseGlob("./tmpl/*"))

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("incoming request to ", r.URL.Path)
	//w.Header().Set("Cache-Control", "public, max-age=10")
	path := r.URL.Path
	switch {
	case path[len("/"):] == "":
		indexHandler(w, r)
	case path[len("/post/"):] != "":
		postHandler(w, r)
	default:
		http.NotFound(w, r)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	posts, _ := queryFirestorePosts(r.Context(), 10)
	t.ExecuteTemplate(w, "index.html.tmpl", posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Path[len("/post/"):]
	p, err := getFirestorePost(r.Context(), postid)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	t.ExecuteTemplate(w, "post.html.tmpl", p)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request to", r.URL.Path)
	p, _ := loadPosts()
	for _, post := range p {
		client := createFirestoreClient(r.Context())
		_, _, _ = client.Collection("posts").Add(r.Context(), post)
	}
}

func loadPosts() ([]*Post, error) {
	entries, _ := os.ReadDir("./posts/")
	var posts []*Post
	for _, file := range entries {
		if strings.HasSuffix(file.Name(), ".json") {
			post := new(Post)
			postfile, _ := os.ReadFile(fmt.Sprintf("./posts/%s", file.Name()))
			json.Unmarshal(postfile, &post)
			posts = append(posts, post)
		}

	}
	return posts, nil
}

func createFirestoreClient(ctx context.Context) *firestore.Client {
	projectID := ""

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	return client
}

func queryFirestorePosts(ctx context.Context, count int) ([]*Post, error) {
	var posts []*Post
	client := createFirestoreClient(ctx)
	results, _ := client.Collection("posts").OrderBy("date", firestore.Desc).Limit(count).Documents(ctx).GetAll()
	for _, doc := range results {
		post := new(Post)
		post.Id = doc.Ref.ID
		doc.DataTo(&post)
		posts = append(posts, post)
	}
	return posts, nil
}

func getFirestorePost(ctx context.Context, id string) (*Post, error) {
	var post *Post
	client := createFirestoreClient(ctx)
	doc, err := client.Collection("posts").Doc(id).Get(ctx)
	doc.DataTo(&post)
	return post, err
}

func main() {
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request){
    http.ServeFile(w, r, fmt.Sprintf(".%s", r.URL.Path))
  })
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
