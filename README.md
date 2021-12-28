# README

This is a simple test that I'm creating to learn Go.

This the list of stuff I wanted to try:

- use of nested structure
- use mongo database and unmarshalling of data
- telegram bot (async calls)
- loading environment variables from .env file
- Dockerize

The software is connected to a Integration Service database 
(the one I'm using in my workdays) when a new message is received
by Telegram, the software will connect to DB and print all secret 
keys stored in it.