FROM golang

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /appa
COPY . /app
RUN go build
RUN ./go-soccer db init
RUN ./go-soccer staker enroll 0x3893e05Ea1Ae26bED0541274942a398aF9385EB7
RUN ./go-soccer staker info

CMD ["./go-soccer", "service", "start"]
