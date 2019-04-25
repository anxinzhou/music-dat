#!/usr/bin/env bash
pvc="PvcERC998"
pbc="PbcERC998"
contractName="ComposableTopDown"

baseContractPath="../etc/contracts/"
targetAPIPath="../token/"

generateGoAPI(){
     #set gopath
     export GOPATH=$(dirname $(dirname $(pwd)))

     #install abigen
    (
       cd ${GOPATH}/src/github.com/ethereum/go-ethereum
       go install ./cmd/abigen
    )
    fileName=$1
    contractName=$2
    basePath=$3
    targetDir=$4${fileName}"/"
    contractPathPrefix=${basePath}${fileName}

    (
        cd ${basePath}
        solc --abi --bin ${fileName}".sol" -o ${fileName} --overwrite
    )

    if [[ ! -d "$targetDir" ]]
    then
         mkdir $targetDir
    fi
    abigen --abi=${contractPathPrefix}"/"${contractName}".abi" --bin=${contractPathPrefix}"/"${contractName}".bin" --pkg=${fileName} -out ${targetDir}${fileName}".go"

   #clear abi and bin file
   for file in `ls $contractPathPrefix`
   do
        rm ${contractPathPrefix}"/"${file}
   done
   rmdir ${contractPathPrefix}

}

generateGoAPI $pvc $contractName $baseContractPath $targetAPIPath
generateGoAPI $pbc $contractName $baseContractPath $targetAPIPath
