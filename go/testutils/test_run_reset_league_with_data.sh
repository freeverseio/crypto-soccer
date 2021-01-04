#!/bin/bash
cd ../synchronizer/process && go test -bench=. -run=^$ -benchmem -cpuprofile profile.out