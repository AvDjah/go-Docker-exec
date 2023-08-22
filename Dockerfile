FROM ubuntu:latest

RUN apt update && apt upgrade -y

RUN apt install python3 -y

COPY main.py /