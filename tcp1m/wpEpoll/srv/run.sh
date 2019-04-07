#!/bin/bash

set -o nounset
#set -o errexit

#set -x

run_native_func()
{
    sudo sysctl -w fs.file-max=3000000
    sudo sysctl -w fs.nr_open=3000000
    sudo sysctl -w net.nf_conntrack_max=3000000
    ulimit -n 3000000

    sudo sysctl -w net.ipv4.tcp_tw_recycle=1
    sudo sysctl -w net.ipv4.tcp_tw_reuse=1

    ./outBin -c=$1 -ec=$2
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
    docker run --rm -v $(pwd)/logDatas:/myApp/logDatas myapp:falcon
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
        run_native_func 20 10
        ;;
    dk) echo "Run in docker..."
#        run_docker_func
        ;;
    clean) echo "Clean datas..."
        run_clean_docker_datas_func
        ;;
    *) echo "Unknown command!"
        usage
        ;;
esac



