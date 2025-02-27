module github.com/LucaIgnatescu/FlatEarthBackend

replace github.com/LucaIgnatescu/FlatEarthBackend => ./

go 1.23.1

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/time v0.8.0
)
