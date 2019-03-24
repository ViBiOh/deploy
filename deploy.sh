#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
RESET='\033[0m'

start_services() {
  if [[ "${#}" -ne 3 ]]; then
    printf "${RED}Usage: start_services [PROJECT_SHA1] [PROJECT_NAME] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local PROJECT_NAME="${2}"
  local COMPOSE_FILE="${3}"

  printf "${GREEN}Starting services for ${PROJECT_NAME}${RESET}\n"

  docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" config -q
  docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" pull
  docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" up -d
}

count_services_with_health() {
  if [[ "${#}" -ne 2 ]]; then
    printf "${RED}Usage: count_services_with_health [PROJECT_SHA1] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local COMPOSE_FILE="${2}"

  local counter=0

  for service in $(docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" ps -q); do
    if [[ $(docker inspect --format '{{ .State.Health }}' "${service}") != '<nil>' ]]; then
      counter=$((counter+1))
    fi
  done

  echo "${counter}"
}

are_services_healthy() {
  if [[ "${#}" -ne 2 ]]; then
    printf "${RED}Usage: are_services_healthy [PROJECT_SHA1] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local COMPOSE_FILE="${2}"

  local WAIT_TIMEOUT="35"

  printf "${YELLOW}Waiting ${WAIT_TIMEOUT} seconds for containers to start...${RESET}\n"
  timeout=$(date --date="${WAIT_TIMEOUT} seconds" +%s)

  local healthcheckCount=$(count_services_with_health "${PROJECT_SHA1}" "${COMPOSE_FILE}")
  local healthyCount=$(docker events --until "${timeout}" -f event="health_status: healthy" -f name="${PROJECT_SHA1}" | wc -l)

  [[ "${healthcheckCount}" == "${healthyCount}" ]] && echo "true" || echo "false"
}

revert_services() {
  if [[ "${#}" -ne 3 ]]; then
    printf "${RED}Usage: revert_services [PROJECT_SHA1] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local COMPOSE_FILE="${2}"

  printf "${YELLOW}Containers didn't start, reverting...${RESET}\n"

  # Force continuation in case of errors for keeping a clean state
  set +e

  docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" logs

  for service in $(docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" ps -q); do
    if [[ $(docker inspect --format '{{ .State.Health }}' "${service}") != '<nil>' ]]; then
      docker inspect --format='{{ .Name }}{{ "\n" }}{{range .State.Health.Log }}code={{ .ExitCode }}, log={{ .Output }}{{ end }}' "${service}"
    fi
  done

  docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" rm --force --stop -v

  set -e
}

remove_old_services() {
  if [[ "${#}" -ne 3 ]]; then
    printf "${RED}Usage: remove_old_services [PROJECT_SHA1] [PROJECT_NAME] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local PROJECT_NAME="${2}"
  local COMPOSE_FILE="${3}"

  printf "${BLUE}Removing old containers from ${PROJECT_NAME}${RESET}\n"

  local projectServices=($(docker ps -f name="${PROJECT_NAME}*" -q))
  local composeServices=($(docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" ps -q))

  local containersToRemove=()

  for projectService in "${projectServices[@]}"; do
    local found=0

    for composeService in "${composeServices[@]}"; do
      if [[ "${projectService:0:12}" == "${composeService:0:12}" ]]; then
        found=1
        break
      fi
    done

    if [[ "${found}" -eq 0 ]]; then
      containersToRemove+=("${projectService}")
    fi
  done

  if [[ "${#containersToRemove[@]}" -gt 0 ]]; then
    docker stop --time=180 ${containersToRemove[@]}
    docker rm -f -v ${containersToRemove[@]}
  fi
}

rename_new_services() {
  if [[ "${#}" -ne 3 ]]; then
    printf "${RED}Usage: rename_new_services [PROJECT_SHA1] [PROJECT_NAME] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local PROJECT_NAME="${2}"
  local COMPOSE_FILE="${3}"

  printf "${BLUE}Renaming containers from ${PROJECT_SHA1} to ${PROJECT_NAME}${RESET}\n\n"

  for service in $(docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" ps --services); do
    local containerID=$(docker ps -q --filter name="${PROJECT_SHA1}_${service}")
    docker rename "${containerID}" "${PROJECT_NAME}_${service}"
  done
}

deploy() {
  if [[ "${#}" -lt 2 ]]; then
    printf "${RED}Usage: deploy [PROJECT_NAME] [SHA1] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_NAME="${1}"
  local PROJECT_SHA1="${PROJECT_NAME}${2}"
  local COMPOSE_FILE="${3:-$(pwd)/docker-compose.yml}"

  start_services "${PROJECT_SHA1}" "${PROJECT_NAME}" "${COMPOSE_FILE}"

  if [[ $(are_services_healthy "${PROJECT_SHA1}" "${COMPOSE_FILE}") == "false" ]]; then
    revert_services "${PROJECT_SHA1}" "${COMPOSE_FILE}"
    return 1
  fi

  remove_old_services "${PROJECT_SHA1}" "${PROJECT_NAME}" "${COMPOSE_FILE}"
  rename_new_services "${PROJECT_SHA1}" "${PROJECT_NAME}" "${COMPOSE_FILE}"

  printf "${GREEN}Deploy successful! ${RESET}\n"
}
