# Scoring Service for Online Multiplayer Game

## Overview
This Go-based application implements a scoring service for an online multiplayer game. It provides endpoints to manage player scores, leveraging MongoDB for persistent data storage and Redis for caching to enhance performance.

## Features
- **Player Score Management**: Add or update player scores.
- **Top Players Retrieval**: Fetch a leaderboard of top players based on their scores.
- **Score Caching**: Efficiently cache player scores in Redis to reduce database load.
- **Player Information Storage**: Store player details using Redis hash.

## Requirements
- Go 1.22+
- Docker and Docker Compose (for containerized deployment)
- MongoDB
- Redis

## Setup Instructions
1. Clone the repository:
   ```bash
   git clone https://github.com/iliatalebzade/quiz_leaderboard.git
   cd quiz_leaderboard
    ```
2. Build and run the application using Docker:
   ```cd docker
   docker-compose up --build
   ```

## API Endpoints
- ```POST /points/add_or_update:``` Add or update a player's score.
- ```GET /points/top_players:``` Retrieve the top players.
- ```GET /points/get_points/:id```: Get the score for a specific player.

## License
### This project is licensed under the MIT License.