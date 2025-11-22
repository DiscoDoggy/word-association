# "Word Association" -- A realtime multiplayer game
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![TailwindCSS](https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoColor=white) [![HTMX](https://img.shields.io/badge/HTMX-36C?logo=htmx&logoColor=fff)](#) 
[![Templ](https://img.shields.io/badge/Templ-36C?logo=htmx&logoColor=fff)](#) 

## How to Play
For the moment, this game is two players. At the beginning of the game, a topic is randomly chosen. For example, the topic car brands may be selected. 
One of the two players are selected to go first. They have 3 seconds to enter a car brand such as "Honda". Once entering and submitting their word, the other player must enter a word related to the topic that has not yet been entered by another player in at most 3 seconds.
The first player to enter an invalid entry or have their time expire before they enter a word loses.

## Why 
As a kid, this is a game I would play with friends in person. However, I could not find a game like this already existing online. I also made this game largely for the learning experience. I had not done a lot of concurrent programming nor had I extensively worked with websockets. 
After working on this project, I possess a greater understanding of concurrency and websockets

## Why this Tech Stack
I achknowledge that this is likely an incredibly unorthodox tech stack to build a multiplayer game with. This stack is known as the GOTTH stack meaning GOlang, Tailwind, Templ, and HTMX. Because this project is intended for pure learning purposes, one of the philopshphies
I have heard is choosing a tech stack that is not meant for the intended purpose can be one of the best ways to learn in general. This tech stack did not end up being too bad though I do feel limited on the clientside by HTMX's ability to enable a lot of interactivity without
the need to write additional javascript.

Using Golang to manage all the games taught me about concurrency and specifically about how Go handles concurrency. There are two large paradigms for concurrency: sharing memory or sending messages. Go subscribes to sending messages through channels. Each websocket client
is represented as a goroutine and each match being played is also a goroutine and each define their own channels for being communicated with. There are still some read-write locks for resources that must be shared by multiple goroutines but using goroutines makes
concurrency feel easier than it probably might be in another language.

One of the large challenges of a multiplayer game was having the clients stay in sync with the server. In a lot of senses, the game is "played" on the server and the clients see an "illusion" of the game. However, because the ground truth of the game is on the server,
the states the client see must closely mirror that of the server. Using templates through templ makes this a bit easier beacuse I can just send the appropriate snippets of HTMX to the clients and Golang and templ rendering are incredily fast so there is little lag.

Another large challenge was the mindset change about how to send data to the client. I am very much used to sending JSON to the UI and having the UI use JSON. However, with templ, HTMX, and websockets, I would be sending HTML over the websockets to the clients to show. 
This was a very large paradigm shift and at times i had troubles thinking about it.

I am clearly not a frontend developer and much of my focus was on the "backend" instead of the UI but improvements to the UI will come. This project as it stands is still a prototype. 

