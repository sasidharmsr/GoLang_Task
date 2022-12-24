package handlers

import (
	"Task/db"
	"Task/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// creating a user
func PostUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		var user DummyUser
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		date, err := time.Parse("2006-01-02", user.DOB) // changing date format
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		newUser := model.User{
			ID:          primitive.NewObjectID(),
			Name:        user.Name,
			DOB:         date,
			Address:     user.Address,
			Description: user.Description,
			CreatedAt:   primitive.NewDateTimeFromTime(time.Now().UTC()),
		}

		res, err := Users.InsertOne(ctx, newUser) // Inserting into Database
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusCreated)
		response := model.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": res}}
		json.NewEncoder(w).Encode(response)

	}
}

// Showing All users in DB
func GetAllUsers(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var users []model.User
		results, err := Users.Find(ctx, bson.M{}) //Here we will get all users from DB

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		//reading users from the db
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser model.User
			if err = results.Decode(&singleUser); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
			}
			users = append(users, singleUser)
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
		json.NewEncoder(w).Encode(response)
	}
}

// Showing Single User
func GetSingleUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := path.Base(fmt.Sprint(r.URL))

		defer cancel()
		var user model.User
		objId, _ := primitive.ObjectIDFromHex(userId)
		err := Users.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(w).Encode(response)
	}
}

// Updating a User
func UpdateUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := path.Base(fmt.Sprint(r.URL))
		defer cancel()
		var user map[string]interface{}
		objId, _ := primitive.ObjectIDFromHex(userId)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		if user["dob"] != nil {
			user["dob"], _ = time.Parse("2006-01-02", fmt.Sprintf("%f", user["dob"]))
		}

		result, err := Users.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": user})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		var updatedUser model.User
		if result.MatchedCount == 1 {
			err := Users.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(w).Encode(response)
	}
}

// Delete a User from DB
func DeleteUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := Users.DeleteOne(ctx, bson.M{"id": objId}) // Deleting User from DB

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		// If user is not  Found in db
		if result.DeletedCount == 0 {
			w.WriteHeader(http.StatusNotFound)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}}
		json.NewEncoder(w).Encode(response)
	}
}

// Follow a User [Post Request]
func FollowUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		defer cancel()

		var Req FollowRequest
		if err := json.NewDecoder(r.Body).Decode(&Req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		var user model.User // Find a USER(X)
		objId, _ := primitive.ObjectIDFromHex(userId)
		err := Users.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		var followuser model.User // Find a User(Y) that the X User wants to Follow
		followobjId, _ := primitive.ObjectIDFromHex(Req.UserID)
		err = Users.FindOne(ctx, bson.M{"id": followobjId}).Decode(&followuser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		user.Following = append(user.Following, followuser)       // X is Following Y
		followuser.Followers = append(followuser.Followers, user) // X is Follower of Y

		result1, err := Users.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": user})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		result2, err := Users.UpdateOne(ctx, bson.M{"id": followobjId}, bson.M{"$set": followuser})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"user": result1, "following_user": result2}}
		json.NewEncoder(w).Encode(response)

	}
}

// Showing A Particlur User Followers
func GetFollowersofUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		defer cancel()

		var user model.User
		objId, _ := primitive.ObjectIDFromHex(userId)
		err := Users.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"Followers": user.Followers}}
		json.NewEncoder(w).Encode(response)

	}
}

// Showing A Particlur User Following Users
func GetFollowingofUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		defer cancel()

		var user model.User
		objId, _ := primitive.ObjectIDFromHex(userId)
		err := Users.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"Followers": user.Following}}
		json.NewEncoder(w).Encode(response)

	}
}

// Showing All USERS Near By A Particlur USer
func GetNearByUsers(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var users []NearestUsersResp
		results, err := Users.Find(ctx, bson.M{}) //Here we will get all users from DB
		params := mux.Vars(r)
		userId := params["id"]

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		//reading a user from the db
		var user model.User
		objId, _ := primitive.ObjectIDFromHex(userId)
		err = Users.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser model.User
			var singleresp NearestUsersResp
			if err = results.Decode(&singleUser); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
			}
			if singleUser.ID.Hex() != userId {
				singleresp.UserId = singleUser.ID.Hex()
				singleresp.UserName = singleUser.Name
				singleresp.Address = singleUser.Address.Addresss
				singleresp.Distance = DistancebetweenLoc(user.Address.Latitude, user.Address.Longitude, singleUser.Address.Latitude, singleUser.Address.Longitude, "K")
				users = append(users, singleresp)
			}
		}
		sort.Slice(users, func(i, j int) bool { return users[i].Distance < users[j].Distance })
		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
		json.NewEncoder(w).Encode(response)
	}
}

// This stuct is only for changing DOB data type
type DummyUser struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name"`
	DOB         string             `json:"dob"`
	Address     model.Address      `json:"address"`
	Description string             `json:"description"`
	CreatedAt   primitive.DateTime `json:"created_at"`
}

// This a Request when user send Follow Request
type FollowRequest struct {
	UserID string `json:"userid"`
}

// This for Nearest USER Response
type NearestUsersResp struct {
	Distance float64
	UserName string
	UserId   string
	Address  string
}
