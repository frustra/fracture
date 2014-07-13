#!/bin/bash

fracture -addr :25565 edge &
sleep 0.2
fracture -addr :25566 -node :7947 -x 0 -z 0 chunk &
fracture -addr :25567 -node :7948 -x -1 -z 0 chunk &
fracture -addr :25568 -node :7949 -x -1 -z -1 chunk &
fracture -addr :25569 -node :7950 -x 0 -z -1 chunk &
fracture -addr :25570 -node :7951 entity &
cat
killall fracture
