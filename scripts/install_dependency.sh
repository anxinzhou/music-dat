#!/usr/bin/env bash

installGoPackage() {
	sudo apt install -y gcc
	echo "installing go package"
	export GOPATH=$(dirname $(dirname "$PWD"))
	go get -u github.com/ethereum/go-ethereum
	go get -u github.com/go-redis/redis
	go get -u github.com/gorilla/handlers
	go get -u github.com/gorilla/mux
	go get -u github.com/gorilla/websocket
	go get -u github.com/goware/disque
	echo "finish installing go package"
}

installGO() {
	echo "installing go"
	if [ -x "$(command -v go)" ];
	then
		echo "have installed go"
		return
	fi 

	if [ "$(uname)" = "Darwin" ];
	then
		brew update
		brew install golang
	elif [ "$(expr substr $(uname -s) 1 5)" = "Linux" ];
	then
		sudo snap install go --classic
	else
		echo "do not support this system"
		exit 1
	fi

	echo "finish installing go"
}

installGeth() {
	echo "installing geth"
	if [ -x "$(command -v geth)" ];
	then
		echo "have installed geth,skip"
		return
	fi 

	if [ "$(uname)" = "Darwin" ];
	then
		brew tap ethereum/ethereum
		brew install ethereum
	elif [ "$(expr substr $(uname -s) 1 5)" = "Linux" ];
	then
		sudo add-apt-repository -y ppa:ethereum/ethereum
		sudo apt-get update
		sudo apt-get install -y ethereum
	else
		echo "do not support this system"
		exit 1
	fi
	echo "finish installing geth"
}

installRedis() {
	echo "installing redis"
	#install in $HOME

	if [ -x "$(command -v redis-server)" ]
	then 
		echo "have installed redis skip"
		return
	fi

	target="redis-5.0.3.tar.gz"
	dst="redis-5.0.3"

	(	
		cd $HOME
		curl -O http://download.redis.io/releases/${target}
		tar -xzvf ${target}
		cd $dst
		sudo apt install -y  make
		sudo apt install -y  make-guile
		make
		sudo ln -s $(pwd)/src/redis-server /usr/local/bin/redis-server
	)

	echo "finish installing redis"
}

installNodeJS() {
	echo "installing node js"
	if [ -x "$(command -v node)" ];
	then
		echo "have installed geth,skip"
		return
	fi

	if [ "$(uname)" = "Darwin" ];
	then
		brew update
		brew install node
	elif [ "$(expr substr $(uname -s) 1 5)" = "Linux" ];
	then
		curl -sL https://deb.nodesource.com/setup_10.x | sudo bash -
		sudo apt install -y nodejs
	else
		echo "do not support this system"
		exit 1
	fi
}
installJsdependency(){
  sudo apt install -y g++
  npm i
}

installParity(){
	sudo snap install parity
}

installGeth 
installGO
installRedis
installNodeJS
installGoPackage
installJsdependency
installParity



