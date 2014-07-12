#!/bin/bash

fracture -addr :25565 edge &
fracture -addr :25566 -node :7947 -x 0 -z 0 chunk &
fracture -addr :25567 -node :7948 -x -8 -z 0 chunk &
fracture -addr :25568 -node :7949 -x -8 -z -8 chunk &
fracture -addr :25569 -node :7950 -x 0 -z -8 chunk &
fracture -addr :25570 -node :7951 player &
cat
killall fracture
