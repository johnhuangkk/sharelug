#!/bin/bash
PID=$$
#cd api
#
#go mod vendor
#
#go run services/dao/createTable.go
#
#process_id=`ps -ef | grep /tmp/go-build | grep -v grep | awk '{print $2}'`
#
#if [[ "" !=  "$process_id" ]]; then
#  echo "killing $process_id"
#  kill -9 $process_id
#fi
#
##go run www/main.go > log/run.log &
#
##go -v run www/main.go
#
##tail -f log/run.log
#
log_file="/var/log/restart_sh.log"

# return the current date time
TIMESTAMP(){
    echo $(date "+%Y-%m-%d %H:%M:%S")
}

stop_bsah_running() {
  echo "get Pid $1"
  process_id=`ps -ef | grep /bin/bash | awk '{print $2}'`
  if [ "" !=  "$process_id" ]
  then
    ids=$(echo $process_id | tr "\n" " ")
    for i in $ids
    do
      echo "ccc $i"
      if [ "$i" != $1 ]
      then
        echo "kill $i successfully!!!"
        kill -9 $i
      else
        echo "$i is not kill"
      fi
    done
  fi
}


stop_process_if_running(){
    # $1->process_name to grep
    process_id=`ps -ef | grep $1 | grep -v grep | awk '{print $2}'`
    echo "running $process_id" | tee -a $log_file
    if [ "$process_id" != "" ]
    then
        echo "$(TIMESTAMP) $1 is running, T'm going to kill it" | tee -a $log_file
        kill -9 $process_id
        echo "kill $1 successfully!!!" | tee -a $log_file
    else
        echo "$(TIMESTAMP) $1 is not running" | tee -a $log_file
    fi
}


restart_process_if_die(){
    process_id=`ps -ef | grep $1 | grep -v grep | awk '{print $2}'`
    echo "running $process_id" | tee -a $log_file
    if [ "$process_id" == "" ];
    then
        echo "$(TIMESTAMP) $3 got down, now I will restart it" | tee -a $log_file
        cd $2
        echo "Now I am in $PWD" | tee -a $log_file
        go run $3 > log/run.log &
        echo "$(TIMESTAMP) $3 restart successfully" | tee -a $log_file
    else
        echo "$(TIMESTAMP) $3 is running, no need to restart" | tee -a $log_file
    fi
}





process="/tmp/go-build"
file_dir=/go/src/sharelug/api
py_file=www/main.go
echo "pid $PID"

ps aux

stop_bsah_running $PID

ps aux
#when execute this shell script, if the process is running,kill it firstly
stop_process_if_running $process

# poll if the process is died, if got died then restart it.
while :
do
    restart_process_if_die $process $file_dir $py_file
    echo "$(TIMESTAMP) now I will sleep 30S" | tee -a $log_file
    sleep 30
done