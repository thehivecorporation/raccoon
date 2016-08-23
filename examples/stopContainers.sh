#!/usr/bin/env bash

sudo docker ps | awk '!/CON/{print $1}' | xargs sudo docker stop
sudo docker ps -a | awk '!/CON/{print $1}' | xargs sudo docker rm