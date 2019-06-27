#!/usr/bin/env bash
#
# Create the student user and allow
# them SSH access using the 21_student_id_rsa key.
#
echo "Setting flags to exit script if any command fails or if any variable is undefined.."
set -eu

echo "Checking that we are the super-user.."
[[ $(id -u) -eq 0 ]] || { echo >&2 "not running as superuser"; exit 1; }

echo "Checking that commands we require are installed.."
command -v adduser >/dev/null 2>&1 || { echo >&2 "adduser is missing"; exit 1; }
command -v mkdir >/dev/null 2>&1 || { echo >&2 "mkdir is missing"; exit 1; }
command -v chown >/dev/null 2>&1 || { echo >&2 "chown is missing"; exit 1; }
command -v chmod >/dev/null 2>&1 || { echo >&2 "chmod is missing"; exit 1; }
command -v sed >/dev/null 2>&1 || { echo >&2 "sed is missing"; exit 1; }
command -v systemctl >/dev/null 2>&1 || { echo >&2 "systemctl is missing"; exit 1; }

echo "Updating packages.."
apt update
echo "Installing useful packages.."
apt install -y mosh jq

echo "Creating student user.."
adduser --shell /bin/bash --disabled-password --gecos "" student
echo "Giving SSH access to 21_student_id_rsa key.."
echo "NOTE: the pubkey below is for the insecure shared 21_student_id_rsa privkey"
echo "get access to the server."
echo "*********************************************"
echo "DO NOT RUN THIS ON A SERVER THAT SHOULD NOT BE ACCESSIBLE TO EVERYONE"
echo "*********************************************"
echo ""
mkdir -p /home/student/.ssh
cat << EOF > /home/student/.ssh/authorized_keys
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDBLrqMonmwIVOWPN+ra1UIztlN0DebYSpu3w4AT4Q2/qxKnBByTTGJm553Cv8s4VCS7FnOQR+vusWHrpJj5vYfPCyft9qamoZ6H9ISYw2CTlMYIRF87TxbkgiQJ/Vt0DxBjCUBmpsgq3VCOPMio4zetxbpo7WHApk50keYCRrqa4iLXNk0kmmCCO4aOPDf4ZRVg5hEauW61CQyny675VoCZdZ8otPMV6+d2ZMqX2rHOhOWVpVRx/y/qzoyvB9EGeQrPPQp0hnO8yz5f43nml5eCDxqMFF2Br3GRuyctdHh0ZWVT912ZRHsVEMkQD4JWOIZiWV305LL7sej5mbXUZYl student@21
EOF
chown -R student:student /home/student/
chmod 400 /home/student/.ssh/authorized_keys

echo "Changing SSH port to 2222.."
sed -i 's/#Port 22/Port 2222/g' /etc/ssh/sshd_config
systemctl restart sshd

echo "Done setting up server."