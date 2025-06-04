# qg-manager

`qg-manager` is a backend service written in Go, designed to manage core functionalities within the Quality-Gamer ecosystem. This service is part of a modular architecture that includes additional components such as `qg-usuario`, `qg-ranking` and `qg-frontend`.

## 📁 Project Structure

The repository is organized as follows:

- `main.go`: Application entry point.
- `conf/`: Configuration files.
- `database/`: Database scripts and migrations.
- `endpoint/`: API routes and handlers.
- `model/`: Data structures and models.
- `Makefile`: Automation for common tasks.
- `Procfile`: Deployment configuration for platforms like Heroku.

## 🚀 Getting Started

### Prerequisites

- Go 1.18 or higher
- [Dep](https://github.com/golang/dep) for dependency management (if using `Gopkg.toml`)

### Running Locally

1. Clone the repository:

   ```bash
   git clone https://github.com/Quality-Gamer/qg-manager.git
   cd qg-manager
   ```

2. Install dependencies:
   ```bash
   dep ensure
   ```
3. Set up environment variables as needed (refer to the files in the conf/ folder for reference).
4. Run the application:
   ```bash
   go run main.go
   ```

For more information about the Quality-Gamer ecosystem, visit the Quality-Gamer GitHub organization.
