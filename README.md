# unique-pass-gen

A unique password generator written in Go, containerized with Docker.

---

## 🚀 Quick Start

### 1️⃣ Build the Docker image

Run this in the project root:

```bash
docker build -t unique-pass-gen .
```

✅ This command will:

* Download the base Ubuntu image
* Install Go 1.24.2
* Install golangci-lint
* Download Go module dependencies
* Run linting and tests during the build
* Compile the HTTP server

---

### 2️⃣ Run the container

```bash
docker run -p 8080:8080 --rm unique-pass-gen
```

✅ This will:

* Start the built container
* Expose port 8080
* Launch the server inside the container

---

### 3️⃣ Test it in your browser

Open:

```
http://localhost:8080
```

You should see a generated password.

---

