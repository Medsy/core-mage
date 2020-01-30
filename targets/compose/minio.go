package compose

var (
	// MinioService represents a docker compose minio service.
	MinioService = Service{
		Image:   "minio/minio",
		Restart: "always",
		Logging: map[string]string{
			"driver": "none",
		},
		Volumes: []string{
			"${PWD}/data/minio/dev:/data",
		},
		Ports: []string{
			"9000:9000",
		},
		Environment: map[string]string{
			"MINIO_ACCESS_KEY": "user",
			"MINIO_SECRET_KEY": "passwd",
		},
		urlPattern: "user:passwd@%s:9000",
	}
	// MinioTestService represents a docker compose minio service for the testing environment.
	MinioTestService = Service{
		Image:   "minio/minio",
		Restart: "always",
		Logging: map[string]string{
			"driver": "none",
		},
		Ports: []string{
			"9001:9001",
		},
		Environment: map[string]string{
			"MINIO_ACCESS_KEY": "user",
			"MINIO_SECRET_KEY": "passwd",
		},
		urlPattern: "user:passwd@%s:9000",
	}
)