# DeBeAndo Agent for MySQL

Database monitoring tool designed for small environments, adapted for Kubernetes and send metrics to InfluxDB.

## Image Description

This image is maintained by DeBeAndo and will be updated regularly on best-effort basis. The image is based on Alpine Linux and only contains the build result of this repository.

## Run

To run container:

```bash
docker run \
	--name debeando-agent-mysql \
	--env DEBUG=true \
	--env INTERVAL=10 \
	--env INFLUXDB_HOST=http://com-env-influxdb-observability-node01.aws.com \
	--env INFLUXDB_TOKEN="abc123cde456==" \
	--env MYSQL_HOST=com-env-mysql-stack-node01.aws.com \
	--env MYSQL_USER=monitor \
	--env MYSQL_PASSWORD=monitor \
	--env SERVER=com-env-mysql-stack-node01 \
	debeando/agent-mysql
```

## MySQL Config

Create a `monitor` user to allow access to agent.

```sql
CREATE USER monitor@'%' IDENTIFIED by 'monitor';
ALTER USER monitor@'%' WITH MAX_USER_CONNECTIONS 5;
GRANT REPLICATION CLIENT ON *.* TO monitor@'%';
GRANT PROCESS ON *.* TO monitor@'%';
GRANT SELECT ON *.* TO monitor@'%';
```

Please, change default password.
