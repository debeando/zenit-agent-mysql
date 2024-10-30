# DeBeAndo Zenit Agent MySQL

Database monitoring tool designed for small environments, adapted for Kubernetes and send metrics to InfluxDB.

## Image Description

This image is maintained by DeBeAndo and will be updated regularly on best-effort basis. The image is based on Alpine Linux and only contains the build result of this repository.

## Run

To run container:

```bash
docker run \
	--name zenit-agent-mysql \
	--env DEBUG=true \
	--env INTERVAL=10 \
	--env HOSTNAME=com-env-mysql-stack-node01 \
	--env INFLUXDB_HOST=com-env-influxdb-observability-node01.aws.com \
	--env INFLUXDB_TOKEN="abc123cde456==" \
	--env MYSQL_HOST=com-env-mysql-stack-node01.aws.com \
	--env MYSQL_USER=monitor \
	--env MYSQL_PASSWORD=passmon \
	debeando/zenit-agent-mysql
```
