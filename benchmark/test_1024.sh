#!/bin/sh
# TODO: change me!
# -p: which partition do you want to run your workload on? <batch, gpu, bigmem>
# -n: how many CPU cores do you want to run your job?
# --mem: how much memory do you want?
# -t: how long do you want to run the job before it timesout <hh:mm:ss>
# --constraint=intel: required for power monitoring

#SBATCH -p batch
#SBATCH -n 1
#SBATCH --mem=1g
#SBATCH -t 60:00
#SBATCH --constraint=intel
for i in {1..1000}
	do
		echo "Loop spin:" $i
		./test_speed1024
	done
