#!/bin/sh

required() {
    var_name=$1
    var_value=$(eval echo \$$var_name)

    if [ -z "$var_value" ]; then
        echo "Error: $var_name is not set or is empty."
        exit 1
    fi
}

append() {
    var_name=$1
    var_value=$2
    [ -z "$(eval echo \$$var_name)" ] && eval "$var_name=\"$var_value\"" || eval "$var_name=\"\$$var_name $var_value\""
}
