package main

import (
	//"encoding/json"
	"log"
	"net/http"
	//"os"

	//"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type ChatMessage struct {
	Text string `json:"text"`
}

/*var (
	rdb *redis.Client
)*/
var clients = make(map[*websocket.Conn]bool)
var broadcaster = make(chan ChatMessage)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true
	/*if rdb.Exists("chat_messages").Val() != 0 {
		sendPreviousMessages(ws)
	}*/

	for {
		var msg ChatMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		// send new message to the channel
		broadcaster <- msg
	}
}

/*func sendPreviousMessages(ws *websocket.Conn) {
	chatMessages, err := rdb.LRange("chat_messages", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	// send previous messages
	for _, chatMessage := range chatMessages {
		var msg ChatMessage
		json.Unmarshal([]byte(chatMessage), &msg)
		messageClient(ws, msg)
	}
}*/
func handleMessages() {
	for {

		// grab any next message from channel
		msg := <-broadcaster
		switch msg.Text {
		case "Veterans":
			msg.Text = "For More Community Experience: Link to Veterans@VMware POD"
		case "Pride":
			msg.Text = "For More Community Experience: Link to Pride@VMware POD"
		case "Disability":
			msg.Text = "For More Community Experience: Link to Disability@VMware POD"
		case "Black":
			msg.Text = "For More Community Experience: Link to Black@VMware POD"
		case "Women":
			msg.Text = "For More Community Experience: Link to Women@VMware POD"
		case "Asian":
			msg.Text = "For More Community Experience: Link to Asian@VMware POD"
		case "Latinos":
			msg.Text = "For More Community Experience: Link to Latinos@VMware POD"

		}
		//storeInRedis(msg)
		messageClients(msg)
	}
}

/*func storeInRedis(msg ChatMessage) {
	json, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	if err := rdb.RPush("chat_messages", json).Err(); err != nil {
		panic(err)
	}
}*/
func messageClients(msg ChatMessage) {
	// send to every client currently connected
	for client := range clients {
		messageClient(client, msg)
	}
}
func messageClient(client *websocket.Conn, msg ChatMessage) {

	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := "8080"
	//redisURL := os.Getenv("REDIS_URL")
	//opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	//rdb = redis.NewClient(opt)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/websocket", handleConnections)
	go handleMessages()

	log.Print("Server starting at localhost:" + port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
