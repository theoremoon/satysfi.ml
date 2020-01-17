name=$(cat /dev/urandom | tr -d -C 0-9A-Za-z | head -c 8)
sock="unix:/tmp/${name}.sock"
workdir="/home/ubuntu/satysfi.ml"
rsync -e "ssh -F ./ssh_config" -zz ./app "satysfi.ml:${workdir}/${name}"
ssh -F ssh_config satysfi.ml $(printf 'cd %s; bgproxyctl green -addr %s:/ -cmd "env PATH=\"\$PATH:/snap/bin\" ./%s --host %s"' "${workdir}" "${sock}" "${name}" "${sock}" )
