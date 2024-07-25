package kvsclient

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/pauljubcse/kvs"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketStore(t *testing.T) {
	store := kvs.NewStore()
	server := httptest.NewServer(http.HandlerFunc(store.HandleWebSocket))
	defer server.Close()

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// Create Domain
	createDomainRequest := kvs.Request{
		Action: "create_domain",
		Domain: "test_domain",
	}
	conn.WriteJSON(createDomainRequest)
	var createDomainResponse kvs.Response
	conn.ReadJSON(&createDomainResponse)
	assert.Equal(t, "success", createDomainResponse.Status)

	// Set String
	setStringRequest := kvs.Request{
		Action: "set_string",
		Domain: "test_domain",
		Key:    "test_key",
		Value:  "test_value",
	}
	conn.WriteJSON(setStringRequest)
	var setStringResponse kvs.Response
	conn.ReadJSON(&setStringResponse)
	assert.Equal(t, "success", setStringResponse.Status)

	// Get String
	getStringRequest := kvs.Request{
		Action: "get_string",
		Domain: "test_domain",
		Key:    "test_key",
	}
	conn.WriteJSON(getStringRequest)
	var getStringResponse kvs.Response
	conn.ReadJSON(&getStringResponse)
	assert.Equal(t, "success", getStringResponse.Status)
	assert.Equal(t, "test_value", getStringResponse.Value)

	// Insert to SkipList
	insertSkipListRequest := kvs.Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
		Value:  "value1",
	}
	conn.WriteJSON(insertSkipListRequest)
	var insertSkipListResponse kvs.Response
	conn.ReadJSON(&insertSkipListResponse)
	assert.Equal(t, "success", insertSkipListResponse.Status)

	// Search in SkipList
	searchSkipListRequest := kvs.Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(searchSkipListRequest)
	var searchSkipListResponse kvs.Response
	conn.ReadJSON(&searchSkipListResponse)
	assert.Equal(t, "success", searchSkipListResponse.Status)
	assert.Equal(t, "value1", searchSkipListResponse.Value)

	// Delete from SkipList
	deleteSkipListRequest := kvs.Request{
		Action: "delete_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(deleteSkipListRequest)
	var deleteSkipListResponse kvs.Response
	conn.ReadJSON(&deleteSkipListResponse)
	assert.Equal(t, "success", deleteSkipListResponse.Status)

	// Confirm Deletion from SkipList
	searchSkipListRequest2 := kvs.Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(searchSkipListRequest2)
	var searchSkipListResponse2 kvs.Response
	conn.ReadJSON(&searchSkipListResponse2)
	assert.Equal(t, "error", searchSkipListResponse2.Status)
	assert.Equal(t, "key not found", searchSkipListResponse2.Message)

	// Insert Range in SkipList
	insertSkipListRequest1 := kvs.Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
		Value:  "value1",
	}
	insertSkipListRequest2 := kvs.Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "2",
		Value:  "value2",
	}
	insertSkipListRequest3 := kvs.Request{
		Action: "insert_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "3",
		Value:  "value3",
	}
	conn.WriteJSON(insertSkipListRequest1)
	conn.ReadJSON(&insertSkipListResponse)
	conn.WriteJSON(insertSkipListRequest2)
	conn.ReadJSON(&insertSkipListResponse)
	conn.WriteJSON(insertSkipListRequest3)
	conn.ReadJSON(&insertSkipListResponse)

	// Delete Range from SkipList
	deleteRangeSkipListRequest := kvs.Request{
		Action: "delete_range_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		MinKey: "1",
		MaxKey: "2",
	}
	conn.WriteJSON(deleteRangeSkipListRequest)
	var deleteRangeSkipListResponse kvs.Response
	conn.ReadJSON(&deleteRangeSkipListResponse)
	assert.Equal(t, "success", deleteRangeSkipListResponse.Status)

	// Confirm Range Deletion from SkipList
	searchSkipListRequest3 := kvs.Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "1",
	}
	conn.WriteJSON(searchSkipListRequest3)
	var searchSkipListResponse3 kvs.Response
	conn.ReadJSON(&searchSkipListResponse3)
	assert.Equal(t, "error", searchSkipListResponse3.Status)
	assert.Equal(t, "key not found", searchSkipListResponse3.Message)

	searchSkipListRequest4 := kvs.Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "2",
	}
	conn.WriteJSON(searchSkipListRequest4)
	var searchSkipListResponse4 kvs.Response
	conn.ReadJSON(&searchSkipListResponse4)
	assert.Equal(t, "error", searchSkipListResponse4.Status)
	assert.Equal(t, "key not found", searchSkipListResponse4.Message)

	searchSkipListRequest5 := kvs.Request{
		Action: "search_skiplist",
		Domain: "test_domain",
		SLKey:  "test_sl",
		Key:    "3",
	}
	conn.WriteJSON(searchSkipListRequest5)
	var searchSkipListResponse5 kvs.Response
	conn.ReadJSON(&searchSkipListResponse5)
	assert.Equal(t, "success", searchSkipListResponse5.Status)
	assert.Equal(t, "value3", searchSkipListResponse5.Value)
}


//Require server to be running else where
func TestIncrementDecrement(t *testing.T) {
	server, err := kvs.StartServer("ws://localhost:9000/ws")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	fmt.Println("Started server...")

	
	client, err := NewClient("ws://localhost:9000/ws")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer client.Close()

	// Create a domain
	domain := "test_domain"
	err = client.CreateDomain(domain)
	if err != nil {
		t.Fatalf("Error creating domain: %v", err)
	}

	// Set a key with initial value
	key := "counter"
	initialValue := "0"
	err = client.SetString(domain, key, initialValue)
	if err != nil {
		t.Fatalf("Error setting string: %v", err)
	}

	// Increment the key
	err = client.Increment(domain, key)
	if err != nil {
		t.Fatalf("Error incrementing counter: %v", err)
	}

	// Check the incremented value
	value, err := client.GetString(domain, key)
	if err != nil {
		t.Fatalf("Error getting string: %v", err)
	}
	expectedValue := "1"
	if value != expectedValue {
		t.Fatalf("Expected value after increment is %s, but got %s", expectedValue, value)
	}

	// Decrement the key
	err = client.Decrement(domain, key)
	if err != nil {
		t.Fatalf("Error decrementing counter: %v", err)
	}

	// Check the decremented value
	value, err = client.GetString(domain, key)
	if err != nil {
		t.Fatalf("Error getting string: %v", err)
	}
	expectedValue = "0"
	if value != expectedValue {
		t.Fatalf("Expected value after decrement is %s, but got %s", expectedValue, value)
	}
	
	fmt.Println("Shutting down server...")
	if err := server.CloseServer(); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server gracefully stopped")

}