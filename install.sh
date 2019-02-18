
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
cd $cliDir

cp root.crt ~/.godrop/root.crt
echo Root Certificate is installed

