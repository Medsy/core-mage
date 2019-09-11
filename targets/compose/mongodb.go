package compose

//   mongodb:
// 	  image: mongo:3.4
// 	  restart: always
// 	  networks:
// 		- app-network

var (
	// MongoService represents a docker compose mongodb service.
	MongoService = Service{
		Image:   "mongo:3.4",
		Restart: "always",
		Logging: map[string]string{
			"driver": "none",
		},
		Volumes: []string{
			"${PWD}/data/mongo/dev:/data/db",
		},
		Ports: []string{
			"27017:27017",
			"27018:27018",
			"27019:27019",
		},
		Environment: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "user",
			"MONGO_INITDB_ROOT_PASSWORD": "passwd",
		},
		urlPattern: "%s:27017",
	}

	// MongoTestService represents a docker compose mongodb service.
	MongoTestService = Service{
		Image:   "mongo:3.4",
		Restart: "always",
		Logging: map[string]string{
			"driver": "none",
		},
		urlPattern: "%s:27017",
	}
)