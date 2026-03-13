# Readers App

A personal project designed to combat information overload and create a focused, personal knowledge base from online content.

## The Problem

It's difficult and time consuming to organice, keep track and mantain a centralized point of access for all the blog posts, articles, tutorials and reading material that you come across in the internet.

## The Solution

Readers App aims to be your personal digital library. Instead of just bookmarking links, you can save the actual content, building a rich, searchable archive of the material that matters to you.

### Core Features

*   **Web Content Archiving:** Provide a link, and the app will scrape and store the content locally.
*   **Centralized Storage:** All your saved articles and posts in one place.
*   **Full-Text Search:** A powerful index allows you to instantly search through the entire body of your collected content.

### Goal

The primary goal is to create an efficient tool for studying, research, and compiling technical knowledge. By gathering and indexing content, you can easily revisit and connect ideas from various sources, turning a collection of articles into a true personal knowledge base.

## Next steps

#### 1. Server-Side Authentication with Firebase

Currently the client is sending credentials to firebase via js to generate tokens and authenticate. This should move to the backend and let the server communicate with Firebase.

**Changes:**
- Remove auth logic in JS (`login.js`, `register.js`).
- Backend uses `Firebase Admin SDK` to handle user creation and token generation.
- Improve error handling with global alerts.

#### 2. Full-Text Search in PostgreSQL

Search content in saved resources. Implement stemming and vectoring to enhance the searching capabilities of the application.

**Changes:**
- Add `tsvector` column to `resources` table.
- Add `pg_trgm` extension and create GIN index.
- Create trigger to automatically update search vector.
- Add search method to repository.
- Add `/api/resource/search` endpoint.

---

*This is a personal project and a work in progress.*
