# LeaderBoard

The project is a module to fetch the top k scores from all the game data.
The game data can be feeded into the module through two ways, the module listens to a game data topic and also through an exposed API.



## Tech Stack

**Golang** is used as the main language. **Redis** is used for message queue to listen to the game topic.
Redis sorted sets are used as in memory data structure to store the top k games.