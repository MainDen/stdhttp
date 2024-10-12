#!/bin/sh

error() {
    [ -z "$1" ] && error "error: argument is empty."
    echo "ERROR: $1" >&2
    exit 1
}

warn() {
    [ -z "$1" ] && error "info: argument is empty."
    echo "WARN:  $1" >&2
}

info() {
    [ -z "$1" ] && error "info: argument is empty."
    echo "INFO:  $1" >&2
}

debug() {
    [ -z "$1" ] && error "debug: argument is empty."
    echo "DEBUG: $1" >&2
}

required() {
    [ -z "$1" ] && error "required: argument is empty."
    var_name=$1
    var_value=$(eval echo \$$var_name)
    [ -z "$var_value" ] && error "Variable $var_name is not set or is empty."
}

append() {
    [ -z "$1" ] && error "append: argument is empty."
    [ -z "$2" ] && error "append: argument is empty."
    var_name=$1
    var_value=$2
    [ -z "$(eval echo \$$var_name)" ] && eval "$var_name=\"$var_value\"" || eval "$var_name=\"\$$var_name $var_value\""
}
