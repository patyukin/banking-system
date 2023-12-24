package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.1.3 -i UserService -o ./mocks/ -s "_minimock.go"
//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock@v3.1.3 -i AuthService -o ./mocks/ -s "_minimock.go"
