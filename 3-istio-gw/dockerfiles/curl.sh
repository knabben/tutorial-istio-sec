#!/bin/bash -x

URI=${2:-headers}

while true; do
   curl http://$1:8000/${URI}
   date
done
