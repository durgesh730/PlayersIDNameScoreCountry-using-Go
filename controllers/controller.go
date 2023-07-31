package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/durgesh730/authenticationInGo/database"
	"github.com/durgesh730/authenticationInGo/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// create players

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ply models.Player
	err := json.NewDecoder(r.Body).Decode(&ply)

	if err != nil {
		http.Error(w, "Fail to save data", http.StatusInternalServerError)
	}

	// check id is exist in data base or not
	filter := bson.M{"id": ply.Id}
	count, err := database.SaveData.CountDocuments(context.Background(), filter)

	// fmt.Println(filter)
	//if exist than return error
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if count > 0 {
		http.Error(w, "Id already exist", http.StatusConflict)
		return
	} else {

		// Check if id is empty
		if ply.Id == "" {
			http.Error(w, "ID field is required", http.StatusBadRequest)
			return
		}

		// Check if name is empty or exceed 15 characters
		if ply.Name == "" {
			http.Error(w, "Name field is required", http.StatusBadRequest)
			return
		} else if len(ply.Name) > 15 {
			http.Error(w, "Name field should not exceed 15 characters", http.StatusBadRequest)
			return
		}

		// Check if country is not exactly 2 characters long
		if len(ply.Country) != 2 {
			http.Error(w, "Country field should contain exactly 2 characters", http.StatusBadRequest)
			return
		}

		// Check if score is less than 0
		if ply.Score < 0 {
			http.Error(w, "Score not be less than 0", http.StatusBadRequest)
			return
		}

		database.SaveData.InsertOne(context.Background(), ply)
		json.NewEncoder(w).Encode(ply)
	}
}

// update data

func UpdateNameandScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updated models.Player
	err := json.NewDecoder(r.Body).Decode(&updated)

	if err != nil {
		log.Fatal(err)
	}

	// Get the student ID from the URL parameter

	params := mux.Vars(r)
	filter := bson.M{"id": params["id"]}

	update := bson.M{"$set": bson.M{
		"name":  updated.Name,
		"score": updated.Score,
	}}

	fmt.Println(filter, update)

	data, err := database.SaveData.UpdateOne(context.Background(), filter, update)

	if err != nil {
		http.Error(w, "Failed to update data in MongoDB", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// delete palyers

func DeleteOnePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	filter := bson.M{"id": params["id"]}

	deleted, err := database.SaveData.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Player deleted", deleted)
	json.NewEncoder(w).Encode(deleted)
}

func GetRandomPlayer(w http.ResponseWriter, r *http.Request) {
	pipeline := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}

	cursor, err := database.SaveData.Aggregate(context.Background(), pipeline)
	if err != nil {
		http.Error(w, "Error while getting a random player from DB", http.StatusInternalServerError)
		return
	}

	var players []models.Player
	if err = cursor.All(context.Background(), &players); err != nil {
		http.Error(w, "Error while decoding player", http.StatusInternalServerError)
		return
	}

	// Check if any player was found
	if len(players) == 0 {
		http.Error(w, "No players in DB", http.StatusInternalServerError)
		return
	}

	// Return the first (and only) player
	json.NewEncoder(w).Encode(players[0])
}

func GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	//descending order by score.
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"score", -1}})

	cursor, err := database.SaveData.Find(context.Background(), bson.D{{}}, findOptions)
	if err != nil {
		http.Error(w, "Error while sorting player", http.StatusInternalServerError)
		return
	}

	var players []models.Player
	if err = cursor.All(context.Background(), &players); err != nil {
		http.Error(w, "Error while decoding player", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}



func GetRankedPlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["val"]
	rank, err := strconv.Atoi(id)

	if err != nil || rank < 1 {
		http.Error(w, "Invalid rank it should not be less than 1", http.StatusInternalServerError)
		return
	}

	pipeline := []bson.M{
		{"$sort": bson.M{"score": -1}}, // Sort by score in descending order
		{"$skip": rank - 1},            // Skip over the players with a higher rank
		{"$limit": 1},                  // Limit to 1 player
	}

	cursor, err := database.SaveData.Aggregate(context.Background(), pipeline)
	if err != nil {
		http.Error(w, "Error while getting player from DB", http.StatusInternalServerError)
		return
	}

	var players []models.Player
	if err = cursor.All(context.Background(), &players); err != nil {
		http.Error(w, "Error while decoding player", http.StatusInternalServerError)
		return
	}

	// Check if any player was found
	if len(players) == 0 {
		http.Error(w, "No player found at the requested rank", http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(players[0])
}
