#!/bin/bash

vagrant halt
vagrant destroy -f
vagrant up --provider virtualbox
