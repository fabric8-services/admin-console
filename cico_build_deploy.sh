#!/bin/bash

. cico_setup.sh

load_jenkins_vars
 run_tests_without_coverage;	if [ ! -f .cico-prepare ]; then
    install_deps
    prepare
     run_tests_without_coverage;
     touch .cico-prepare
fi

deploy;
