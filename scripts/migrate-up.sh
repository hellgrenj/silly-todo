migrate -database postgresql://silly:silly@localhost/silly?sslmode=disable -path ../db/migrations up && echo $1
