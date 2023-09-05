<h1 align="center">Authentication Project</h1>

## Table of contents
- [Table of contents](#table-of-contents)
- [Description](#description)
- [Requirements](#requirements)
- [Usage](#usage)
- [Operation](#operation)
- [Author](#author)

## Description
This project is an authentication system implemented in the Go programming language. It provides user signup, login and user management functionality. The project uses a MongoDB database for storing user information and JWT (JSON Web Tokens) for authentication.

## Requirements

- Go 1.21.0
- MongoDB 5.0.20

**OBS:** Make sure have MongoDB installed and running on your machine.

## Usage

1. Clone the repository.
2. Edit the `.env` file in the project root directory, add a strong secret key for token generation:

        SECRET_KEY=your_secret_key_here

3. Install the necessary dependencies by running:
    
       go mod tidy

4. Run the project using the following command:

        go run main.go

The application will start, and you can access the API at http://localhost:9000. Make sure you have MongoDB installed and running on your machine with the specified URI.

## Operation

The project structure consists of several directories and files, each serving a specific purpose in building the authentication system. Below is a detailed explanation of the project structure::

```
├── go.mod
├── go.sum
├── main.go
├── controllers
│   └── userController.go
├── database
│   └── connection.go
├── helpers
│   ├── authHelper.go
│   └── tokenHelper.go
├── middleware
│   └── authMiddleware.go
├── models
│   └── models.go
└── routes
    ├── authRouter.go
    └── userRouter.go
```

- **main.go:** The main entry point of the application that sets up the server and defines API routes.
- **.env:** Configuration file for setting environment variables such as the port and MongoDB URI.
- **controllers/:** Contains controller functions for handling user signup, login, and user retrieval.
- **database/:** Contains the database connection setup code.
- **helpers/:** Contains helper functions for token generation, validation, and user type checking.
- **middleware/:** Contain middleware function for user authentication.
- **models/:** Defines the data model for the User entity.
- **routes/:** Contains route definitions for user and authentication-related endpoints.

This project is a basic example of user authentication and can be extended and customized as needed for accomplish specific requirements.

## Author
| [<img src="https://avatars.githubusercontent.com/u/105685220?v=4" width=115><br><sub>Pedro Vinicius Messetti</sub>](https://github.com/pedromessetti) |
| :---------------------------------------------------------------------------------------------------------------------------------------------------: |
