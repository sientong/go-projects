package handlers

import (
	"log"
	"websocket/chatting-app/utils"

	"encoding/json"
	"net/http"
)

type RequestChannel struct {
	ChannelName string `json:"channelName"`
}

type Channel struct {
	ID           int64
	ChannelName  string
	IsRestricted bool
	IsActive     bool
	CreatedAt    string
	UpdatedAt    string
}

func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Accessing channel handler")
	if r.Method == http.MethodPost {
		addNewChannel(w, r)
		return
	}

	if r.Method == http.MethodGet {
		channels, err := getChannelList()
		if err != nil {
			log.Printf("Error on fetch channels: %s", err)
			http.Error(w, "Failed to fetch channels list", http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(channels); err != nil {
			log.Printf("Error on encode users: %s", err)
			http.Error(w, "Failed to encode users", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPatch {
		updateChannel()
		return
	}
}

func getChannelList() ([]Channel, error) {
	channels := []Channel{}

	rows, err := utils.GetDatabaseConnection().Query("SELECT * FROM channels")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var channel Channel
		if err := rows.Scan(&channel.ID, &channel.ChannelName, &channel.IsActive, &channel.IsRestricted, &channel.CreatedAt, &channel.UpdatedAt); err != nil {
			log.Printf("Error on scanning channels list: %s", err)
			return nil, err
		}
		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error: %s", err)
		return nil, err
	}
	log.Println(len(channels))

	return channels, nil
}

func getIndividualChannel() {

}

func addNewChannel(w http.ResponseWriter, r *http.Request) {
	var newChannel Channel
	err := json.NewDecoder(r.Body).Decode(&newChannel)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO channels (channel_name, is_restricted) VALUES ($1, $2)"
	_, err = utils.GetDatabaseConnection().Exec(query, newChannel.ChannelName, false)
	if err != nil {
		log.Printf("Failed to save new channel with error %s", err)
		http.Error(w, "Failed to save new channel", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	defer r.Body.Close()
}

func updateChannel() {

}
