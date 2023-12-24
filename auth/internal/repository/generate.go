package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.1.3 -i UserRepository -o ./mocks/ -s "_minimock.go"
//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.1.3 -i AuthRepository -o ./mocks/ -s "_minimock.go"
