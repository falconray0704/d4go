#!/bin/bash

set -o nounset
#set -o errexit

#set -x

run_native_func()
{

    #./outBin
    #./outBin -c="./configs/appCfgs.yaml"
    DATE=`date -d "+3 minutes" +"%FT%T %z"`
    ./outBin -conn=$1 -c=$2 -ip=$3 -sm="${DATE}"
    echo "ret:$?"

}

run_native_cmd_test_func()
{
    set -x

    ./outBin -c="./configs/appCfgs.yaml" -word=opt -numb=7 -fork -svar=flag
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -word=opt
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -word=opt a1 a2 a3
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -word=opt a1 a2 a3 -numb=7
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -h
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -wat
    echo "ret:$?"

    ./outBin -c="./configs/appCfgs.yaml" -loop=true -word=opt a1 a2 a3 looping
    echo "ret:$?"

    set +x
}

run_docker_func()
{
    #docker run --rm -v $(pwd)/logDatas:/myApp/logDatas myapp:falcon
    CONNECTIONS=$1
    CONCURRENCY=$2
    REPLICAS=$3
    IP=$4

    DATE=`date -d "+2 minutes" +"%FT%T %z"`
    #go build --tags "static netgo" -o client client.go
    for (( c=0; c<${REPLICAS}; c++ ))
    do
        docker run --rm -v $(pwd)/outBin:/client --name 1mclient_$c -d ubuntu /client -conn=${CONNECTIONS} -c=${CONCURRENCY} -ip=${IP} -sm "${DATE}"
    done
}

run_clean_docker_datas_func()
{
    rm -rf logDatas/*
}

usage()
{
    echo "Run native:"
    echo "./run.sh lc"
    echo ""
    echo "Run docker:"
    echo "./run.sh dk"
    echo ""
    echo "Run clean datas:"
    echo "./run.sh clean"
}

[ $# -lt 1 ] && usage && exit

mkdir -p ./logDatas

case $1 in
    lc) echo "Run native..."
#        run_native_func $2 $3 $4
        run_native_func 10000 5 172.17.0.1
        ;;
    dk) echo "Run in docker..."
#        run_docker_func $2 $3 $4
        run_docker_func 10000 2 50 172.17.0.1
        ;;
    dkstop) echo "Stopping docker..."
        #docker stop $(docker ps -a --format '{{.ID}} {{.Names}}' | grep '1mclient_' | awk '{print $1}')
        docker ps --format '{{.Names}}' | grep "^1mclient" | awk '{print $1}' | xargs -I {} docker stop {}
        ;;
    clean) echo "Clean datas..."
        run_clean_docker_datas_func
        ;;
    *) echo "Unknown command!"
        usage
        ;;
esac



