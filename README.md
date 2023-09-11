# LiLiStyle_Back_Assignment_11Sep_AngJoonKitEugene

Prerequisites:
1. Ensure that Docker is installed on your system.

Instructions:
1. Clone the Repository
2. Clone the project's repository to your local machine.
3. Navigate to the Project Directory
4. Build the Docker Image using the following command

bash
docker build -t product-like-image .

5. Start the PostgreSQL Database:
Use the docker-compose.yml file to start both the application and the PostgreSQL database with a the command:

bash
docker-compose up -d

OR 
If you don't use Docker Compose, you can run the PostgreSQL container separately with the following command:

bash
docker run -d \
  --name postgres-1 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgrespass \
  -e POSTGRES_DB=product_like_db \
  postgres:latest

6. Run the Application Container using the following command:

bash
docker run -p 8080:8080 -d product-like-image

7. References command
docker stop postgres-1  # Stop the PostgreSQL container
docker rm postgres-1    # Remove the PostgreSQL container
docker stop <app_container_id>  # Stop the application container
docker rm <app_container_id>    # Remove the application container
Replace <app_container_id> with the actual ID of your application container.
