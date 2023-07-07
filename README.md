# PDF Management Project

This project is a PDF management application that allows users to sign up, log in, upload PDF files, view their uploaded PDFs, and share PDFs using sharable links. 
It is built using Go programming language and utilizes various packages and libraries for handling HTTP requests, rendering templates, interacting with a database, 
and managing user sessions.I have used the postgresql database in this project.

## Features

**User Signup:**  Users can create an account by providing their name, email, and password.

**User Login:**   Registered users can log in using their email and password.

**Home Page:**    Displays the signup form for new users.

**Login Page:**   Displays the login form for existing users.

**PDF Page:**     Allows users to upload a PDF file along with a title.

**Options Page:**   Provides options for managing PDF files, including viewing all uploaded PDFs.

**Store PDF:**    Stores the uploaded PDF file in the database associated with the logged-in user.

**Get PDF:**     Retrieves a specific PDF file from the database and displays it in the browser.

**Get All PDFs:**  Lists all PDF files uploaded by the logged-in user, along with their respective shareable links.

**Handle PDF:**   Handles the request to view a specific PDF file based on its encoded data.

**Get Link:**    Retrieves a PDF file based on a shareable link and displays it in the browser.

## Installation and Setup

  **Clone the repository:**
  
    git clone https://github.com/your-username/pdf-management.git

  **Install the dependencies:**

    go mod download

  **Configure the database connection in the database package to connect to your preferred database system.**
  

  **Generate a secret key for session management and update it in the handlers package:**

     var store = sessions.NewCookieStore([]byte("your-secret-key"))
     
  **Build and run the application:**

     go build
    ./pdf-management


  **Access the application in your web browser at http://localhost:8080.**
  

## Usage

  **Sign Up:**  Visit the home page and fill out the signup form with your name, email, and password.

  **Log In:**  Use your registered email and password to log in to the application.

  **Upload PDF:**  After logging in, navigate to the PDF page and upload a PDF file by selecting the file and providing a title.

  **View Uploaded PDFs:**  Go to the options page to see a list of all the PDFs you have uploaded. Each PDF will have a button to view it and a shareable link.

  **Share PDFs:**  You can share your uploaded PDFs with others by sharing the respective shareable link.


## Dependencies

  **github.com/go-chi/chi/v5:**  A lightweight, idiomatic, and composable router for building Go HTTP services.

  **github.com/google/uuid:**  Package uuid provides a pure Go implementation of Universally Unique Identifiers (UUIDs).

  **github.com/gorilla/sessions:**  A package that provides cookie and filesystem sessions and infrastructure for custom session backends.

  **golang.org/x/crypto/bcrypt:** A package that implements cryptographic hashing functions using the Blowfish algorithm.

  **database/sql:**  The database/sql package provides a generic interface around SQL (or SQL-like) databases.

  **encoding/base64:**  Package base64 implements base64 encoding and decoding.

  Make sure to download these dependencies using **go mod download** before running the application.


## Database Used

  ### PostgreSQL

  This project utilizes PostgreSQL as the chosen database system. PostgreSQL is a powerful, open-source, and feature-rich relational database management system (RDBMS) 
  known for its stability, extensibility, and adherence to industry standards.


## Contributing

  Contributions to this project are welcome. If you find any issues or want to add new features, please open an issue or submit a pull request on the project's 
  GitHub repository.


## Acknowledgments

  The project structure and code organization were inspired by best practices in Go programming.

  The project utilizes various open-source packages and libraries, which are mentioned in the "Dependencies" section.

  Special thanks to the contributors who helped in developing and testing this application.


    
