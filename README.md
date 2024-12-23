# MentoRef Web App

MentoRef is a modern web application built with Go and HTMX (as well as _hyperscript) and includes the usage of MongoDB and TailwindCSS.

## Setup

> [!IMPORTANT]
> This app requires setting up certain environment variables.
> Also, the static assets like logos used in this project are not public and only used here for development to showcase the functionality.
> So other, self-obtained assets need to be used.

.env variables with **only placeholders**:
```
MONGODB_URI="MONGODB_URI"
SECRET_KEY="SECRET_KEY"
PRIVATE_KEY="path/to/private-key.key"
CERTIFICATE="path/to/cert.crt"
```
Run `npm i` to **install all dependencies**.
For building the **TailwindCSS** styles also run `npm run build`.

Lastly, build the **Go** application and run the binary:
Run the following from project root in your terminal to build `$ go build -o bin/app cmd/main.go` (output optional) and to run the app `$ ./bin/app`.

Creating a docker image and container is also possible by running `docker-compose up -d`.

## Features (WIP)

- User authentication - Sign up, sign in, and sign out functionality with JWT-based authentication (OAuth implementation in review)
- Dashboard - A dashboard for users (very rudimentary - more to come soon)
- Dynamic - The app uses HTMX and _hyperscript to provide dynamic interactivity without the need for classic JavaScript code

## Technologies

- The backend is built with Go
- The app uses MongoDB as its database
- The frontend uses HTMX and _hyperscript
- The app is styled using standard CSS (mainly for View Transition internal) together with TailwindCSS
