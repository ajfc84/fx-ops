if [ -z "${CI}" ] \
&& [ "${CI_ENVIRONMENT_NAME}" != "local" ];
then
  if [ -z "${DOCKER_HUB_PASSWORD}" ];
  then
    read_passwd "DOCKER_HUB_PASSWORD"
  fi
fi
