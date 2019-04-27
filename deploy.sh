#!/usr/bin/env bash

set -o nounset -o pipefail -o errexit

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

  printf "${counter}"
}

are_services_healthy() {
  if [[ "${#}" -ne 2 ]]; then
    printf "${RED}Usage: are_services_healthy [PROJECT_SHA1] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local COMPOSE_FILE="${2}"

  local runningContainers=$(docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" ps -q | wc -l)
  printf "${BLUE}${runningContainers} running${RESET}"

  # add `-a` options when issues on docker-compose is resolved
  local allContainers=$(docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" ps -q | wc -l)
  printf "${BLUE}${allContainers} total${RESET}"

  if [[ "${runningContainers}" != "${allContainers}" ]]; then
    printf "false"
    return
  fi

  local WAIT_TIMEOUT="35"

  printf "${YELLOW}Waiting ${WAIT_TIMEOUT} seconds for containers to start...${RESET}\n"
  timeout=$(date --date="${WAIT_TIMEOUT} seconds" +%s)

  local healthcheckCount=$(count_services_with_health "${PROJECT_SHA1}" "${COMPOSE_FILE}")
  printf "${BLUE}${healthcheckCount} with healthcheck${RESET}"

  local healthyCount=$(docker events --until "${timeout}" -f event="health_status: healthy" -f name="^${PROJECT_SHA1}" | wc -l)
  printf "${BLUE}${healthyCount} healthy${RESET}"

  if [[ "${healthcheckCount}" != "${healthyCount}" ]]; then
    printf "false"
    return
  fi

  printf "true"
}

revert_services() {
  if [[ "${#}" -ne 2 ]]; then
    printf "${RED}Usage: revert_services [PROJECT_SHA1] [DOCKER-COMPOSE-FILE]${RESET}\n"
    return 1
  fi

  local PROJECT_SHA1="${1}"
  local COMPOSE_FILE="${2}"

  printf "${YELLOW}Containers didn't start, reverting...${RESET}\n"

  # Force continuation in case of errors for keeping a clean state
  set +e

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

  local projectServices=($(docker ps -f name="^${PROJECT_NAME}*" -q))
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
    local containerID=$(docker ps -a -q --filter name="${PROJECT_SHA1}_${service}")
    if [[ -n "${containerID}" ]]; then
      docker rename "${containerID}" "${PROJECT_NAME}_${service}"
    fi
  done
}

clean() {
  printf "${BLUE}Cleaning docker system${RESET}\n"

  set +e
  docker system prune -f
  set -e
}

deploy() {
  if [[ "${#}" -lt 1 ]]; then
    printf "${RED}Usage: deploy [PROJECT_NAME] [SHA1] [DOCKER-COMPOSE-FILE]\n"
    printf "  where\n"
    printf "    - PROJECT_NAME         Name of your compose project\n"
    printf "    - SHA1                 Unique identifier of your project (default: git sha1 of commit)\n"
    printf "    - DOCKER_COMPOSE_FILE  Path to your compose file (default: docker-compose.yml in current dir)\n"
    printf "${RESET}"

    return 1
  fi

  if [[ "$(docker-compose version --short)" < "1.24.0" ]]; then
    printf "${RED}You need at least docker-compose@1.24.0, please upgrade${RESET}"

    return 1
  fi

  local PROJECT_NAME="${1}"
  local PROJECT_SHA1="${PROJECT_NAME}${2:-$(git rev-parse --short HEAD)}"
  local COMPOSE_FILE="${3:-docker-compose.yml}"

  start_services "${PROJECT_SHA1}" "${PROJECT_NAME}" "${COMPOSE_FILE}"

  is_healthy=$(are_services_healthy "${PROJECT_SHA1}" "${COMPOSE_FILE}")
  docker-compose -p "${PROJECT_SHA1}" -f "${COMPOSE_FILE}" logs
  
  if [[ "${is_healthy}" == "false" ]]; then
    revert_services "${PROJECT_SHA1}" "${COMPOSE_FILE}"
    return 1
  fi

  remove_old_services "${PROJECT_SHA1}" "${PROJECT_NAME}" "${COMPOSE_FILE}"
  rename_new_services "${PROJECT_SHA1}" "${PROJECT_NAME}" "${COMPOSE_FILE}"

  printf "${GREEN}Deploy successful! ${RESET}\n"

  clean
}

deploy "${@}"
