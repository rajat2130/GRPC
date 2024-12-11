package controller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"grpc/go-grpc-crud-api/conifg"
	pb "grpc/go-grpc-crud-api/proto"
	"pti_api/config"
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID        string
	Title     string
	Genre     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Server struct {
	pb.UnimplementedMovieServiceServer
}

// CreateMovie handles the creation of a movie record
func (*Server) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	fmt.Println("Create Movie")
	movie := req.GetMovie()
	movie.Id = uuid.New().String()

	// SQL query to insert the movie into the database
	query := "INSERT INTO movies (id, title, genre, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := conifg.DB.Exec(query, movie.GetId(), movie.GetTitle(), movie.GetGenre(), time.Now(), time.Now())
	if err != nil {
		return nil, errors.New("movie creation unsuccessful: " + err.Error())
	}

	return &pb.CreateMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.GetId(),
			Title: movie.GetTitle(),
			Genre: movie.GetGenre(),
		},
	}, nil
}

// GetMovie retrieves a movie by its ID
func (*Server) GetMovie(ctx context.Context, req *pb.ReadMovieRequest) (*pb.ReadMovieResponse, error) {
	fmt.Println("Read Movie", req.GetId())
	var movie Movie

	// SQL query to retrieve the movie by ID
	query := "SELECT id, title, genre, created_at, updated_at FROM movies WHERE id = ?"
	row := conifg.DB.QueryRow(query, req.GetId())

	// Scan the result into the Movie struct
	err := row.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}

	return &pb.ReadMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

// GetMovies retrieves all movies
func (*Server) GetMovies(ctx context.Context, req *pb.ReadMoviesRequest) (*pb.ReadMoviesResponse, error) {
	fmt.Println("Read Movies")
	var movies []Movie

	// SQL query to retrieve all movies
	query := "SELECT id, title, genre, created_at, updated_at FROM movies"
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and populate the movies slice
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert movies to protobuf response
	moviesProto := make([]*pb.Movie, len(movies))
	for i, movie := range movies {
		moviesProto[i] = &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		}
	}

	return &pb.ReadMoviesResponse{
		Movies: moviesProto,
	}, nil
}

// UpdateMovie updates an existing movie
func (*Server) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	fmt.Println("Update Movie")
	reqMovie := req.GetMovie()

	// SQL query to update the movie
	query := "UPDATE movies SET title = ?, genre = ?, updated_at = ? WHERE id = ?"
	_, err := config.DB.Exec(query, reqMovie.Title, reqMovie.Genre, time.Now(), reqMovie.Id)
	if err != nil {
		return nil, errors.New("movie update unsuccessful: " + err.Error())
	}

	return &pb.UpdateMovieResponse{
		Movie: &pb.Movie{
			Id:    reqMovie.Id,
			Title: reqMovie.Title,
			Genre: reqMovie.Genre,
		},
	}, nil
}

// DeleteMovie deletes a movie by its ID
func (*Server) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	fmt.Println("Delete Movie")

	// SQL query to delete the movie by ID
	query := "DELETE FROM movies WHERE id = ?"
	_, err := config.DB.Exec(query, req.GetId())
	if err != nil {
		return nil, errors.New("movie deletion unsuccessful: " + err.Error())
	}

	return &pb.DeleteMovieResponse{
		Success: true,
	}, nil
}
