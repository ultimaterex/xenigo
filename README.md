# Xenigo

Xenigo is a Reddit monitoring tool designed to fetch posts from specified subreddits and send notifications to configured targets such as Discord or Slack via webhooks. It features a built-in cache to prevent duplicate notifications and includes a catch-up mechanism to ensure no posts are missed if the service is temporarily down.

## Getting Started

### Prerequisites

- Docker
- Docker Compose (optional, if you want to use it)

### Building and Running the Application

1. **Clone the repository:**

    ```sh
    git clone https://github.com/yourusername/xenigo.git
    cd xenigo
    ```

2. **Build the Docker image:**

    ```sh
    docker build -t xenigo .
    ```

3. **Run the Docker container:**

    ```sh
    docker run -d -p 8080:8080 --name xenigo xenigo
    ```

### Running with Docker Compose

1. **Create a `docker-compose.yml` file with the following content**:

   ```yaml
   services:
     conscript:
       command:
         - './xenigo'
       container_name: 'xenigo'
       hostname: 'xenigo'
       image: 'ghcr.io/ultimaterex/xenigo/xenigo:latest'
       restart: 'unless-stopped'
       volumes:
         - "./config.yaml:/config.yaml" 
   ```

2. **Run the Docker Compose setup**:

   ```sh
   docker-compose up -d
   ```

### Configuration

The application uses a `config.yaml` file for configuration. you can find an example in the file repo


### Contributing

Feel free to submit issues or pull requests. For major changes, please open an issue first to discuss what you would like to change.

### License

This project is licensed under the MIT License.