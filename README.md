# Terminal Portfolio (TUI)

> A fully interactive SSH-based portfolio built with Go, Bubble Tea, and Docker.
> **Try it live:** `ssh dev.tarunnayaka.me`

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)
![Azure](https://img.shields.io/badge/Hosted_on-Azure-0078D4?style=flat&logo=microsoftazure)

[demo video](https://github.com/user-attachments/assets/5de1f09b-2ad9-49b4-b78a-3dcf74646b0e)


## 🚀 About The Project

This isn't your average portfolio website. It is a **Terminal User Interface (TUI)** designed for developers who live in the command line. It provides a rich, keyboard-driven experience accessible via a simple SSH command—no installation required.

### Key Features
* **SSH Accessible:** No browser needed. Just `ssh <domain>`.
* **Visuals in Terminal:** Renders project images as high-quality ASCII art.
* **Responsive Layout:** Adapts to different terminal window sizes dynamically.
* **Hybrid Data Layer:** Fetches metadata from **MongoDB** and assets from **Azure Blob Storage**.
* **Performance:** Implements in-memory caching to reduce database load.

## 🛠️ Architecture

The system is containerized with Docker and hosted on an Azure VM. It uses **Charm's Wish** library to serve the application over SSH.
![demo video](https://github.com/user-attachments/assets/3de13d44-2b53-49e2-a06a-a497fddcb5cc)

### Tech Stack
* **Core:** Go (Golang)
* **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lipgloss](https://github.com/charmbracelet/lipgloss)
* **SSH Server:** [Wish](https://github.com/charmbracelet/wish)
* **Database:** MongoDB Atlas (Prisma Schema)
* **Storage:** Azure Blob Storage
* **Infrastructure:** Docker on Azure Linux VM

## ⚡ Quick Start

### For Users
Simply open your terminal and run:
```bash
ssh dev.tarunnayaka.me
```

For Developers (Running Locally)

    Clone the repo
```Bash

git clone https://github.com/Rtarun3606k/Portfolio-Bubbletea-TUI
cd portfolioTUI
```

Setup Environment Create a .env file:
Code snippet

MONGO_URI=your_mongodb_connection_string
AZURE_STORAGE_ACCOUNT=your_account_name
AZURE_STORAGE_KEY=your_account_key

Run with Docker
```Bash
docker build -t tui-app .
docker run -p 23234:23234 --env-file .env tui-app
Then ssh into localhost: ssh -p 23234 localhost
```

📂 Project Structure
Bash
❯ tree                                                                       ─╯
.
├── config
│   └── configs.go
├── database
│   └── connection.go
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── tui
│   ├── Blogs.go
│   ├── cache.go
│   ├── contact.go
│   ├── Home.go
│   ├── model.go
│   ├── Positions.go
│   ├── Project.go
│   ├── styles.go
│   ├── update.go
│   └── view.go
└── utils
    ├── AsciiImage.go
    └── tools.go


🤝 Contributing

Contributions are welcome! Please open an issue or submit a PR if you have ideas for new widgets or optimizations.
📜 License

Distributed under the MIT License. See LICENSE for more information.
