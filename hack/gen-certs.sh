#!/bin/bash
#
# Copyright (c) 2015 W. Mark Kubacki <wmark@hurrikane.de>
# Licensed under the terms of the RPL 1.5 for all usages
#   http://www.opensource.org/licenses/rpl1.5
#

set -e -o pipefail

CAsubj="/C=DE/ST=Niedersachsen/L=Hannover/O=Dummy CA/CN=Sign-It-All"
CApath="dummy/CA"
ClientSubj="/C=DE/O=Dummy Corp/CN=" # the CN value gets appended
ClientPath="dummy/"                 # the last /* will be stripped

mkdir "${CApath%/*}"

# Makes up a CSR and generates an unique key.
#
# @param1	path without ext, will create files $1.{key,csr}
# @param2	string to be used as 'subj' for a CSR
function makecsr() {
	umask 0177
	openssl ecparam -genkey -name prime256v1 -out "${1}.key"
	umask 0022

	# CSR, in case you want to submit it to any known CA
	# see: openssl req -in web.csr -noout -text
	openssl req -new -nodes -sha384 \
		-key "${1}.key" -subj "${2}" -out "${1}.csr"
}

# Creates a dummy CA.
# Uses ${CApath} and ${CAsubj}.
function create_CA() {
	makecsr "${CApath}" "${CAsubj}"
	# We issue ourselves a self-signed cert for the CA
	# without any key constraints or extended usages (,=all permitted):
	openssl req -new -x509 -sha384 -set_serial 1 -days 3 \
		-key "${CApath}.key" -subj "${CAsubj}" -out "${CApath}.crt"
}

# Signs a CSR. Uses ${CApath}.* as CA.
#
# @param1	path to the certificate to be issued, without the ext;
#               $1.csr will be used as 'signing request'
function issue_cert() {
	local random_serial=$(tr -dc '0-9' < /dev/urandom | head -c 8 || true)

	## Your in-house CA would use:
	##     openssl ca -sha384 -config … -name … -extensions …
	openssl x509 -req -sha384 -set_serial ${random_serial} -days 1 \
		-CAkey "${CApath}.key" -CA "${CApath}.crt" \
		-extfile "extensions.cnf" -extensions "for_a_node" \
		-in "${1}.csr" -out "${1}.crt"
}

# create a dummy CA…
create_CA

# … and certificates for nodes {A,B,C}
for handle in "server"; do
	makecsr "${ClientPath%/*}/${handle}" "${ClientSubj}${handle}"
	issue_cert "${ClientPath%/*}/${handle}"

	# view it by: openssl x509 -noout -text -in …/….crt
done

# fin
echo DONE
