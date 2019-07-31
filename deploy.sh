#!/usr/bin/env bash

set -o nounset -o pipefail -o errexit

deploy() {
  if [[ "${#}" -lt 1 ]]; then
    printf "${RED}Usage: deploy [PROJECT_NAME] [DOCKER-COMPOSE-FILE]\n"
    printf "  where\n"
    printf "    - PROJECT_NAME         Name of your compose project\n"
    printf "    - DOCKER_COMPOSE_FILE  Path to your compose file (default: docker-compose.yml in current dir)\n"
    printf "${RESET}"

    return 1
  fi

  local PROJECT_NAME="${1}"
  local COMPOSE_FILE="${2:-docker-compose.yml}"

  local RED='\033[0;31m'
  local GREEN='\033[0;32m'
  local BLUE='\033[0;34m'
  local RESET='\033[0m'

  printf "${GREEN}Starting services of ${PROJECT_NAME}${RESET}\n"

  docker-compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" config -q
  docker-compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" pull
  docker-compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" up -d --remove-orphans

  printf "${GREEN}Logs for services of ${PROJECT_NAME}${RESET}\n"

  docker-compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" logs

  printf "${BLUE}Cleaning${RESET}\n"

  set +e
  docker rmi $(docker images -q)
  docker network rm $(docker network ls -q)
  set -e

  printf "${GREEN}Deploy successful! ${RESET}\n"
}

deploy "${@}"
