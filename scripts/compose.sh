#!/usr/bin/env bash

SCRIPT_PATH=$( cd "$(dirname "$0")" ; pwd -P )
COMPOSE_PATH=${SCRIPT_PATH}/../docker
COMPOSE_FILE=docker-compose.${1}.yml

if [[ ! -e "${COMPOSE_PATH}/${COMPOSE_FILE}" ]];then
  echo "${COMPOSE_FILE} not found."
  exit 1
fi

function clean(){
  cd "${COMPOSE_PATH}" && docker-compose -f ${COMPOSE_FILE} down -v
}

function build(){
  cd "${COMPOSE_PATH}" && docker-compose -f ${COMPOSE_FILE} build
}

function up(){
  cd "${COMPOSE_PATH}" && docker-compose -f ${COMPOSE_FILE} up --force-recreate -d
}

function down(){
  cd "${COMPOSE_PATH}" && docker-compose -f ${COMPOSE_FILE} down -v
}

case "${2}" in
    up)
        up
        ;;
    build)
        build
        ;;
    down)
        down
        ;;
    stop)
        down
        ;;
    start)
        up
        ;;
    clean)
        clean
        ;;
    restart)
        down
        clean
        up
        ;;
    *)
        echo $"Usage: $0 {up|down|build|start|stop|clean|restart}"
        exit 1
esac