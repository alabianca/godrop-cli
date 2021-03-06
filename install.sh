
# Basic Install of godrop-cli. installs godrop into $GOPATH/bin
# Then creates the .godrop dir in the home directory of the user
# Finally copy the root certfificate into the .godrop directory

cd godrop
echo installing CLI ...
go install
echo CLI Installed
cd ..

cliDir="$(pwd)"
echo $cliDir
cd ~
if [ -d ".godrop" ]; then
    echo deleting godrop files
    rm -rf .godrop
fi
mkdir .godrop
cd .godrop
echo Downloading root certfificate ...
curl -G http://104.248.183.179:80/root > root.crt
echo Root Certfificate installed
# cd $cliDir

# cp root.crt ~/.godrop/root.crt
# echo Root Certificate is installed

