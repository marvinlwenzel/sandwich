#!/bin/bash
# https://stackoverflow.com/questions/75696690/how-to-resolve-tls-failed-to-verify-certificate-x509-certificate-signed-by-un
echo "CA candicates"
strace curl https://discordapp.com/ |& grep open | grep -E "(crt|ca)" | cut -d"\"" -f 2