package manageusers

//TODO: Implement this, as of now it is just placeholder code
import (
	"net/http"

	"github.com/athosone/golib/pkg/server/routing"
	"github.com/go-chi/chi/v5"
)

const (
	ProduceUserV1              = "application/vnd.athosone.projectraven.user+*; v=1"
	ConsumeRegisterUserInputV1 = "application/vnd.athosone.projectraven.registerUser+json; v=1"
)

type UserHandler struct {
}

func NewRestUserHandler() *UserHandler {
	return &UserHandler{}
}

func (handler *UserHandler) AddRoutes(router chi.Router) {
	router.Route("/user", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Get user"))
		})
		userRouter := routing.NewRouter()
		userRouter.Post(handler.RegisterUser).Consume(ConsumeRegisterUserInputV1).Produce(ProduceUserV1).SetDefault()
		r.Mount("/register", userRouter)
	})
}

// swagger:route POST /user/register user registerUser
// Register a new user
// Consumes:
// - application/vnd.athosone.projectraven.registerUser+json; v=1
// Produces:
// - application/vnd.athosone.projectraven.user+*; v=1
// Security:
// - ApiKeyAuth: []
// Responses:
//
//	   default: errorResponse
//	   201: registerUserResponse
//	   400: validationError
//		 401: errorResponse
//		 403: errorResponse
//		 500: errorResponse
//
// Parameters:
//
//	+name: registerUserInput
//	  in: body
//	  required: true
//	  type: registerUserInput
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implemented"))
}

// Models

// Register user response
//
// swagger:response registerUserResponse
type RegisterUserResponse struct {
	// The user
	// in: body
	Body struct {
		// The user
		// Required: true
		User UserModelV1 `json:"user"`
	}
}

// UserModel
//
// swagger:model user
type UserModelV1 struct {
	// The user id
	ID string `json:"id"`
	// The user email
	Email string `json:"email"`
	// The user identity provider
	//
	// Example: google
	IDP string `json:"idp"`
}

// swagger:parameters registerUserInput
type RegisterUserInputV1 struct {
	// in: body
	// required: true
	// min length: 1
	Email string `json:"email"`
	// in: body
	// required: true
	IDP string `json:"idp"`
}
