package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"learning/internal/graphql/dataloaders"
	"learning/internal/graphql/generated"
	"learning/internal/graphql/resolvers"
	"learning/internal/order"
	"learning/internal/product"
	"learning/internal/user"
)

// Config holds GraphQL server configuration
type Config struct {
	PlaygroundEnabled    bool
	IntrospectionEnabled bool
}

// Server represents GraphQL server
type Server struct {
	config     *Config
	handler    http.Handler
	playground http.Handler
}

// NewServer creates a new GraphQL server
func NewServer(
	config *Config,
	userRepo user.Repository,
	productRepo product.Repository,
	orderRepo order.Repository,
) *Server {
	// Create dataloaders
	loaders := dataloaders.NewLoaders(userRepo, productRepo, orderRepo)

	// Create resolver with dependencies
	resolver := resolvers.NewResolver(userRepo, productRepo, orderRepo, loaders)

	// Create GraphQL config
	gqlConfig := generated.Config{Resolvers: resolver}

	// Create GraphQL handler
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))

	// Add middleware for complexity analysis, query timeout, etc.
	// TODO: Add these middleware for production
	// gqlHandler.Use(extension.Introspection{})
	// gqlHandler.Use(extension.AutomaticPersistedQuery{})

	var playgroundHandler http.Handler
	if config.PlaygroundEnabled {
		playgroundHandler = playground.Handler("GraphQL Playground", "/graphql")
	}

	return &Server{
		config:     config,
		handler:    dataloaders.Middleware(loaders)(gqlHandler),
		playground: playgroundHandler,
	}
}

// Handler returns the GraphQL HTTP handler
func (s *Server) Handler() http.Handler {
	return s.handler
}

// PlaygroundHandler returns the GraphQL Playground handler
func (s *Server) PlaygroundHandler() http.Handler {
	return s.playground
}

// ServeHTTP implements http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
