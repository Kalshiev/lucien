package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/Kalshiev/lucien/internal/database"
	"github.com/Kalshiev/lucien/pkg/app"
	"github.com/Kalshiev/lucien/pkg/middleware"
	"github.com/Kalshiev/lucien/pkg/models"
	"github.com/Kalshiev/lucien/pkg/rest"
	"github.com/google/uuid"
)

// LendBook lends a book to a borrower
func LendBook(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		borrowerName := rest.GetPathParam(r, "borrowerName")
		bookID := rest.GetPathParam(r, "bookID")

		bookUUID, err := uuid.Parse(bookID)
		if err != nil {
			fmt.Println("Bad Book ID")
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		parsedUserID, err := uuid.Parse(validatedUser)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		loan, err := a.DB.CreateLoan(r.Context(), database.CreateLoanParams{
			ID:       uuid.New(),
			Lender:   parsedUserID,
			Borrower: borrowerName,
			Book:     bookUUID,
		})
		if err != nil {
			fmt.Println("Loan failed: ", err)
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = a.DB.UpdateBook(r.Context(), database.UpdateBookParams{
			IsAvailable: false,
			Borrower:    sql.NullString{String: borrowerName, Valid: true},
		})
		if err != nil {
			fmt.Println("Updating Book Status Failed: ", err)
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, loan)
	}
}

// ReturnBook returns a borrowed book
func ReturnBook(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := rest.GetPathParam(r, "bookID")
		bookUUID, err := uuid.Parse(bookID)
		if err != nil {
			fmt.Println("Bad Book ID")
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		validatedUser := middleware.GetUserID(r)
		fmt.Println("User logged in: ", validatedUser)

		loanRecord, err := a.DB.ReturnLoan(r.Context(), bookUUID)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println("Loan record updated: ", loanRecord)

		_, err = a.DB.UpdateBook(r.Context(), database.UpdateBookParams{
			IsAvailable: true,
			Borrower:    sql.NullString{Valid: false},
		})
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.RespondJSON(w, http.StatusOK, loanRecord)
	}
}

// GetLoanHistory gets the loan history for a book
func GetLoanHistory(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatedUser := middleware.GetUserID(r)
		log.Println("Logged in user: ", validatedUser)

		id := rest.GetPathParam(r, "bookID")
		bookID, err := uuid.Parse(id)
		if err != nil {
			rest.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		loans, err := a.DB.GetLoanHistory(r.Context(), bookID)
		if err != nil {
			rest.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := make([]models.LoanResponse, len(loans))
		for idx, loan := range loans {
			response[idx] = models.LoanResponse{
				Id:         loan.ID,
				Lender:     loan.Lender,
				Borrower:   loan.Borrower,
				Book:       loan.Book,
				LentAt:     loan.LentAt,
				ReturnedAt: loan.ReturnedAt.Time,
			}
		}

		if r.URL.Query().Get("sort") == "desc" {
			sort.Slice(response, func(i, j int) bool { return response[i].LentAt.After(response[j].LentAt) })
		}

		rest.RespondJSON(w, http.StatusOK, response)
	}
}
