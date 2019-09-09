default:
	@echo "Please use specific target"

# NOTE: Needs to be kept in sync with docker-compose.yml

# How many clients to spawn
N=5
# Client target for commands like att
T=4

stop:
	docker-compose stop

# Start all containers and feed stdout to terminal
up: stop
	docker-compose up --scale client=$(N)

# Start all containers in background
upbk: stop
	docker-compose up -d --scale client=$(N)

# Attach to specific containers
att:
	docker attach thummcoin_client_$(T)

# Listen to specific containers
list:
	docker logs -f thummcoin_client_$(T)

# Compile protobuffers
proto:
	protoc -I=prot --go_out=prot prot/*.proto
