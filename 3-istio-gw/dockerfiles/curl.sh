#!/bin/bash -x

URI=${2:-headers}

while true; do
   curl http://$1:8000/${URI} -v
   date
   sleep 1
done
