#!/bin/bash

sudo docker run -d -P --name test_sshd1 rastasheep/ubuntu-sshd:14.04
sudo docker run -d -P --name test_sshd2 rastasheep/ubuntu-sshd:14.04
sudo docker run -d -P --name test_sshd3 rastasheep/ubuntu-sshd:14.04

sudo docker port test_sshd1 22
sudo docker port test_sshd2 22
sudo docker port test_sshd3 22