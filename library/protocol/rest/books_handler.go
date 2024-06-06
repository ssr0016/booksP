package rest

import (
	"encoding/json"
	"net/http"
	"practice/api/errors"
	"practice/api/response"
	"practice/library/books"
	"practice/routing"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) NewBooksHandler(r *routing.Router) {
	groupBO := r.Group("/api/bo/books")

	groupBO.GET("/", s.searchBooks)
	groupBO.POST("/", s.createBooks)
	groupBO.GET("/:id", s.getByIDBook)
	groupBO.PUT("/:id", s.updateBook)
}

func (s *Server) searchBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		query       books.BooksQuery
		queryValues = r.URL.Query()
	)

	if len(queryValues[""]) != 1 {
		if err := decoder.Decode(&query, queryValues); err != nil {
			response.Error(http.StatusBadRequest, errors.ErrBadRequest).WriteTo(w)
			return
		}
	}

	result, err := s.Dependencies.BooksSvc.Search(r.Context(), &query)
	if err != nil {
		response.Error(http.StatusInternalServerError, err).WriteTo(w)
		return
	}

	response.JSON(http.StatusOK, result).WriteTo(w)
}
func (s *Server) createBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var cmd books.CreateBooksCommand
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		response.Error(http.StatusBadRequest, errors.ErrBadRequest).WriteTo(w)
		return
	}

	err = cmd.Validate()
	if err != nil {
		response.Error(http.StatusBadRequest, err).WriteTo(w)
		return
	}

	err = s.Dependencies.BooksSvc.Create(r.Context(), &cmd)
	if err != nil {
		response.Error(http.StatusInternalServerError, err).WriteTo(w)
	}

	response.Success("Book created").WriteTo(w)
}

func (s *Server) getByIDBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		response.Error(http.StatusBadRequest, errors.ErrBadRequest).WriteTo(w)
		return
	}

	result, err := s.Dependencies.BooksSvc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(http.StatusInternalServerError, err).WriteTo(w)
		return
	}

	response.JSON(http.StatusOK, result).WriteTo(w)
}

func (s *Server) updateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cmd books.UpdateBooksCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		response.Error(http.StatusBadRequest, errors.ErrBadRequest).WriteTo(w)
		return
	}

	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		response.Error(http.StatusBadRequest, errors.ErrBadRequest).WriteTo(w)
		return
	}

	cmd.ID = id
	if err := cmd.Validate(); err != nil {
		response.Error(http.StatusBadRequest, err).WriteTo(w)
		return
	}

	response.Success("Book updated successfully").WriteTo(w)
}
