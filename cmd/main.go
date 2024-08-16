package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Reaper1994/go-serverless-app/package/handlers"
	"github.com/Reaper1994/go-serverless-app/package/middlewares"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func initializeMiddleware(handler http.Handler) http.Handler {
	treblleAPIKey := os.Getenv("TREBLLE_API_KEY")
	treblleProjectID := os.Getenv("TREBLLE_PROJECT_ID")

	// TODO: Remove this later
	fmt.Printf("TREBLLE_API_KEY %s\n", os.Getenv("TREBLLE_API_KEY"))
	fmt.Printf("TREBLLE_PROJECT_ID on port %s\n", os.Getenv("TREBLLE_PROJECT_ID"))

	handler = middlewares.TreblleMiddleware(treblleAPIKey, treblleProjectID, handler)

	return handler
}

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

const tableName = "go-serverless-crud-app"

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)})

	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)

	// Convert the Lambda handler to an http.Handler
	httpHandler := lambdaHandlerToHTTPHandler(lambdaHandler)

	// Initialize the middleware with the handler
	httpHandlerWithMiddleware := initializeMiddleware(httpHandler)

	// Start the Lambda function with the wrapped handler
	lambda.Start(httpHandlerToLambda(httpHandlerWithMiddleware))
}

// lambdaHandler handles the actual Lambda logic
func lambdaHandler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, tableName, dynaClient)
	case "POST":
		return handlers.CreateUser(req, tableName, dynaClient)
	case "PUT":
		return handlers.UpdateUser(req, tableName, dynaClient)
	case "DELETE":
		return handlers.DeleteUser(req, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}

// lambdaHandlerToHTTPHandler converts a Lambda handler to an http.Handler
func lambdaHandlerToHTTPHandler(lambdaHandler func(events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body) // Updated to use io.ReadAll
		if err != nil {
			http.Error(w, "could not read request body", http.StatusInternalServerError)
			return
		}

		// Convert the http.Request to an APIGatewayProxyRequest
		req := events.APIGatewayProxyRequest{
			HTTPMethod: r.Method,
			Path:       r.URL.Path,
			Body:       string(bodyBytes),
		}

		resp, err := lambdaHandler(req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Convert the APIGatewayProxyResponse to an http.ResponseWriter response
		for k, v := range resp.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(resp.Body))
	})
}

// httpHandlerToLambda converts an http.Handler to a Lambda-compatible function
func httpHandlerToLambda(handler http.Handler) func(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return func(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		r, _ := http.NewRequest(req.HTTPMethod, req.Path, bytes.NewBuffer([]byte(req.Body)))
		w := newLambdaResponseWriter()
		handler.ServeHTTP(w, r)
		return w.response, nil
	}
}

// lambdaResponseWriter is a custom HTTP response writer to capture the response for API Gateway.
type lambdaResponseWriter struct {
	response *events.APIGatewayProxyResponse
}

func newLambdaResponseWriter() *lambdaResponseWriter {
	return &lambdaResponseWriter{response: &events.APIGatewayProxyResponse{
		Headers: make(map[string]string),
	}}
}

func (w *lambdaResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *lambdaResponseWriter) Write(body []byte) (int, error) {
	w.response.Body = string(body)
	return len(body), nil
}

func (w *lambdaResponseWriter) WriteHeader(statusCode int) {
	w.response.StatusCode = statusCode
}
