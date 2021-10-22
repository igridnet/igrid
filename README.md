# igrid

Main Components (Functional)
- Auth
- MQ
- Registry

## compose

Adjust the .env file accordingly as docker-compose will pick the environmental variables automatically

```bash

docker-compose up --build 

```

This will start postgres, mosquitto and igrid

## register admin

```text
POST /admins
```
```json
{
  "name": "nameofadmin",
  "email": "demo@emailcompany.com",
  "password": "strong-password-example"
}
```

## admin login
specify username and password for Basic Auth
```text
GET /login
```

## add new node
```text
POST /nodes
```
```json
{
  "addr": "ip-address of the node",
  "name": "water sensor",
  "region": "sregion-id",
  "lat": "7337763563563",
  "long": "w7827287828289",
  "master": "master node",
  "type": 1
}
```
note: node type specification can be any of the 3 values 0,1 and 2 where
- SensorNode = 0
- ActuatorNode = 1
- ControllerNode = 2


## add new region 
```text
POST /regions
```
```json
{
  "name": "region name",
  "description": "region description"
}
```

## publish
example
- -h host
- -p port
- -t topic
- -m message

```bash

mosquitto_pub -h 192.168.1.1 -p 1885 -t sensors/temperature -m "1266193804 32"

```

note: topic are written as region-id/node-id or just region-id if the publishing/subscribing node is controller node

## subscribe
```bash

mosquitto_sub -t sensors/temperature

```