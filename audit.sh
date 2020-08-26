#!/bin/bash

set +e

#
# Set Colors
#

bold="\e[1m"
dim="\e[2m"
underline="\e[4m"
blink="\e[5m"
reset="\e[0m"
red="\e[31m"
green="\e[32m"
blue="\e[34m"

#
# Common Output Styles
#

h1() {
  printf "\n${bold}${underline}%s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
h2() {
  printf "\n${bold}%s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
info() {
  printf "${dim}➜ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
success() {
  printf "${green}✔ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
error() {
  printf "${red}${bold}✖ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
warnError() {
  printf "${red}✖ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
warnNotice() {
  printf "${blue}✖ %s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}
note() {
  printf "\n${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$(echo "$@" | sed '/./,$!d')"
}

typeExists() {
  if [ $(type -P $1) ]; then
    return 0
  fi
  return 1
}

if [ "x${BIN_PATH}x" = "xx" ]; then
  if ! typeExists "gw-aws-audit"; then
    error "gw-aws-audit is not installed"
    note "To install run: curl https://i.jpillora.com/GoodwayGroup/gw-aws-audit! | bash"
    note "Or use BIN_PATH=<path to binary> ./audit.sh"
    exit 1
  fi
  BIN=gw-aws-audit
else
  BIN=$BIN_PATH
fi

US="us-east-1 us-east-2 us-west-1 us-west-2"
EU="eu-central-1 eu-west-1 eu-west-2 eu-west-3 eu-south-1 eu-north-1"
AP="ap-east-1 ap-south-1 ap-northeast-3 ap-northeast-2 ap-southeast-1 ap-southeast-2 ap-northeast-1"
CHINA="cn-north-1 cn-northwest-1"
ROW="af-south-1 me-south-1 sa-east-2"
ALL="$US $EU $AP $ROW $CHINA"

if [[ "x$1x" == "xx" || "$1" == "-h" || "$1" == "--help" ]]; then
  h1 "audit.sh helper script for gw-aws-audit"
  h2 "Usage:"
  cat <<EOF
    audit.sh [gw-aws-audit commands]
EOF
  h2 "Examples:"
  cat <<EOF
> This will run the 'gw-aws-audit sg detached' command for every region in the US (default)

    $ audit.sh sg detached

> This will run the 'gw-aws-audit ec2 stopped-hosts' for ONLY the us-west-2 region

    $ AWS_REGION=us-west-2 audit.sh ec2 stopped-hosts

> This will run the 'gw-aws-audit ec2 stopped-hosts' for every region in the EU

    $ REGION=eu audit.sh ec2 stopped-hosts

> This will run the 'gw-aws-audit cw monitoring' using a specific version of the tool.

    $ BIN_PATH=./bin/gw-aws-audit audit.sh cw monitoring
EOF
  note "REGION env values (default: US):"
  cat <<EOF
US: $US
EU: $EU
AP: $AP
CH: $CHINA
ROW: $ROW
ALL: All of the above combined

You can also set AWS_REGION and that will supersede the value of REGION
EOF
  success "Have fun!"
  exit 0
fi

if [ "x${AWS_REGION}x" = "xx" ]; then
  case $REGION in
  us | US)
    note "Processing for US Regions"
    CHECK_REGIONS=$US
    ;;
  ap | AP)
    note "Processing for Asia Pacific Regions"
    CHECK_REGIONS=$AP
    ;;
  eu | EU)
    note "Processing for EU Regions"
    CHECK_REGIONS=$EU
    ;;
  ch | CH | china | CHINA)
    note "Processing for China Regions"
    CHECK_REGIONS=$CHINA
    ;;
  row | ROW)
    note "Processing for Rest of World (ME, SA, AF) Regions"
    CHECK_REGIONS=$ROW
    ;;
  *)
    note "Defaulting to US Regions"
    CHECK_REGIONS=$US
    ;;
  esac
else
  CHECK_REGIONS=$AWS_REGION
fi

info "Regions: $CHECK_REGIONS"
info "Executing: $BIN ${@}"

for AWS_REGION in $CHECK_REGIONS; do
  h1 "AWS_REGION=$AWS_REGION"
  AWS_REGION=$AWS_REGION $BIN ${@}
  echo ""
done

success "Done!"
