package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Kalshiev/lucien/internal/auth"
	"github.com/Kalshiev/lucien/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiCfg) handlerLendBook(w http.ResponseWriter, r *http.Request) {
	borrowerName := r.PathValue("borrowerName")
	bookID := r.PathValue("bookID")

	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		fmt.Println("Bad Book ID")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	validatedUser, err := auth.ValidateJWT(reqToken, cfg.tokenSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	loan, err := cfg.db.CreateLoan(r.Context(), database.CreateLoanParams{
		ID:       uuid.New(),
		Lender:   validatedUser,
		Borrower: borrowerName,
		Book:     bookUUID,
	})
	if err != nil {
		fmt.Println("Loan failed: ", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = cfg.db.UpdateBook(r.Context(), database.UpdateBookParams{
		IsAvailable: false,
		Borrower:    sql.NullString{String: borrowerName, Valid: true},
	})
	if err != nil {
		fmt.Println("Updating Book Status Failed: ", err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, loan)

}

func (cfg *apiCfg) handlerReturnBook(w http.ResponseWriter, r *http.Request) {
	bookID := r.PathValue("bookID")
	bookUUID, err := uuid.Parse(bookID)
	if err != nil {
		fmt.Println("Bad Book ID")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	validatedUser, err := auth.ValidateJWT(reqToken, cfg.tokenSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	fmt.Println("User logged in: ", validatedUser)

	loanRecord, err := cfg.db.ReturnLoan(r.Context(), bookUUID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Loan record updated: ", loanRecord)

	_, err = cfg.db.UpdateBook(r.Context(), database.UpdateBookParams{
		IsAvailable: true,
		Borrower:    sql.NullString{Valid: false},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJson(w, http.StatusOK, loanRecord)
}

func (cfg *apiCfg) handlerGetLoanHistory(w http.ResponseWriter, r *http.Request) {
	var loans []database.Loan
	var err error

	type Loan struct {
		Id         uuid.UUID `json:"id"`
		Lender     uuid.UUID `json:"lender"`
		Borrower   string    `json:"borrower"`
		Book       uuid.UUID `json:"book"`
		LentAt     time.Time `json:"lent_at"`
		ReturnedAt time.Time `json:"returned_at"`
	}

	reqToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	validatedUser, err := auth.ValidateJWT(reqToken, cfg.tokenSecret)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	log.Println("Logged in user: ", validatedUser)

	id := r.PathValue("bookID")
	if id != "" {
		bookID, Parserr := uuid.Parse(id)
		if Parserr != nil {
			respondError(w, http.StatusBadRequest, Parserr.Error())
			return
		}
		loans, err = cfg.db.GetLoanHistory(r.Context(), bookID)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	response := make([]Loan, len(loans))

	for idx, loan := range loans {
		response[idx] = Loan{
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

	respondJson(w, http.StatusOK, response)
}
