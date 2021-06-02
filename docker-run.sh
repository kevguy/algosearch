mkdir -p db-data
# https://stackoverflow.com/a/38753971
export DOCKERHOST=$(ifconfig | grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}" | grep -v 127.0.0.1 | awk '{ print $2 }' | cut -f2 -d: | head -n1)
docker-compose rm
docker-compose up --build --force-recreate --remove-orphans
