#!/bin/bash
echo "sleeping 3 seconds and then running migrations" 
sleep 3 
echo "running migrations"
migrate -database postgresql://silly:silly@postgres/silly?sslmode=disable -path migrations up