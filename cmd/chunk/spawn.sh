#!/bin/zsh

for x in {-2..2}
do
	for z in {-2..2}
	do
		./chunk -x=$x -z=$z &
	done
done

