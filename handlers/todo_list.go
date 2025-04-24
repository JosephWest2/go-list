package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"josephwest2.com/go-list/app"
	"josephwest2.com/go-list/components"
	"josephwest2.com/go-list/sqlc"
)

func TodoListPageHandler(appContext app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := app.ReadAuthToken(appContext, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userID := token.Claims.(*app.Claims).UserID
		q := sqlc.New(appContext.DBpool)
		items, err := q.GetUserItems(context.Background(), int32(userID))
		if err != nil {
			println(err)
		}
		RenderPage(components.TodoListPage(items), w)
	}
}

func CreateTodoListItemHandler(appContext app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := app.ReadAuthToken(appContext, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authentication Failed"))
			return
		}
		userID := token.Claims.(*app.Claims).UserID
		q := sqlc.New(appContext.DBpool)
		item, err := q.CreateItem(context.Background(), sqlc.CreateItemParams{
			UserID: int32(userID),
			Value:  r.FormValue("value"),
		})
		if err != nil {
			log.Fatal("Failed to create Item: ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"id": %d}`, item.ID)))
	}
}

func DeleteTodoListItemHandler(appContext app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := app.ReadAuthToken(appContext, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authentication Failed"))
			return
		}
		userID := token.Claims.(*app.Claims).UserID
		q := sqlc.New(appContext.DBpool)
		itemIdString := r.PathValue("id")
		itemId, err := strconv.ParseInt(itemIdString, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid item id"))
			return
		}
		item, err := q.GetItem(context.Background(), int32(itemId))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Item not found"))
			return
		}
		if item.UserID != int32(userID) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		err = q.DeleteItem(context.Background(), item.ID)
		if err != nil {
			log.Fatal("Failed to delete Item: ", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Item successfully deleted"))
	}
}
