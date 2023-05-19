# write a golang docker file to build the app
FROM golang:alpine

# set the working directory
WORKDIR /app

# copy the current directory contents into the container at /app
COPY . .

# build the app
RUN go build -o welcome_to_build .

# expose port 8080 to the outside world
EXPOSE 8080

# command to run the executable
CMD ["./welcome_to_build"]