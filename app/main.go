package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"app/pubsub"

	ps "cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
)

type Service struct {
	pubsub *ps.Client
}

/*
📖 => Search for a line of code to uncomment
🎬 => An action you need to perform

== Part 1: Publish your first message and receive it ==
 1️⃣📖 Open http://localhost:8080/. It says "Welcome at..."
      Go to homeHandler to change the text and reload.

 2️⃣🎬 Open http://localhost:8080/submit. Read submitHandler to
	  see what is happening. The error is "Topic not found",
	  let's add a topic.

 3️⃣📖 Create a topic (uncomment lines)
   🎬 Open http://localhost:8080/submit. This time it should
      say "Message published, ID: 1". Refresh the page a few times.

 4️⃣📖 Add a push subscription (uncomment lines)
   🎬 Open http://localhost:8080/submit and check docker-compose logs.
      They should show: Decoded message: "Hello World"
	  What is the MessageID? Subscriptions start to receive only from
	  the moment they are created.

  5️⃣📖 Add a second push subscription to the topic, posting to the
       same PUSH_URL endpoint
	🎬 Open http://localhost:8080/submit and check docker-compose logs.
	   You should see the same message twice.

  🎤 Timeouts, retries, undeliverable, authentication with IAM,
     message retention, pull subscriptions.
*/

func main() {
	ctx := context.Background()
	client, err := ps.NewClient(ctx, os.Getenv("PUBSUB_PROJECT_ID"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create topic and subscription
	service := &Service{pubsub: client}
	// 3️⃣ Create a topic (uncomment next line)
	// service.createTopicIfNotExists("jobs")

	// 4️⃣ Add a push subscription (uncomment this block)
	// service.createPushSubscriptionIfNotExists(
	// 	client.Topic("jobs"),  // The topic you created
	// 	"push1",               // The name of the subscription
	// 	os.Getenv("PUSH_URL"), // The destination of the message
	// )

	// 5️⃣📖 Add a second push subscription, same topic, same URL
	// service.createPushSubscriptionIfNotExists(
	// // Fill this in yourself
	// )

	// Set up URL handlers
	r := mux.NewRouter()
	// The homepage
	r.HandleFunc("/", service.homeHandler)
	// Handler to publish a message
	r.HandleFunc("/submit", service.submitHandler)
	// Handler to receive a message
	r.HandleFunc("/push", service.pushHandler)

	r.Use(logging)
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Service) homeHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	// 1️⃣ Change the homepage text
	fmt.Fprintf(w, "Welcome at %s", hostname)
}

func (s *Service) submitHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	data, _ := json.Marshal("Hello World")
	publishResult := s.pubsub.Topic("jobs").Publish(
		ctx,
		&ps.Message{
			Data: data,
		},
	)
	result, err := publishResult.Get(ctx)
	if err != nil {
		sendErr(w, err)
		return
	}
	fmt.Fprintf(w, "Message published, ID: %v", result)
}

func (s *Service) pushHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := pubsub.DecodePubSubMessage(r.Body)
	log.Printf("Decoded message: %s", string(res))
}

func sendErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	fmt.Println(err)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header = nil
		r.Host = ""
		request, _ := httputil.DumpRequest(r, true)

		log.Println("=====")
		log.Println(string(request))
		next.ServeHTTP(w, r)

	})
}

func (s *Service) createTopicIfNotExists(name string) *ps.Topic {
	_, err := s.pubsub.CreateTopic(context.Background(), name)
	if err != nil && !strings.Contains(err.Error(), "AlreadyExists") {
		log.Fatalf("Failed to create topic: %v", err)
	}
	return s.pubsub.Topic("jobs")
}

func (s *Service) createPushSubscriptionIfNotExists(
	topic *ps.Topic,
	name string,
	url string) {
	_, err := s.pubsub.CreateSubscription(
		context.Background(),
		name,
		ps.SubscriptionConfig{
			Topic: topic,
			PushConfig: ps.PushConfig{
				Endpoint: url,
			},
		},
	)
	if err != nil && !strings.Contains(err.Error(), "AlreadyExists") {
		log.Fatalf("failed to create subscription: %s", err)
	}
}
