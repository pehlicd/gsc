# gsc (Gitlab Structured Cloner)

gsc is a tool to help you clone all the repositories from a Gitlab group in a way that you see the repositories in the same structure as they are in the Gitlab group.

## How to install

### Using Homebrew

```bash
brew tap pehlicd/tap
brew install gsc
```

### Using go

```bash
go install github.com/pehlicd/gsc@latest
```

## Usage

```bash
gsc - GitLab Structured Cloner

gsc is a tool to help you clone all the repositories from a Gitlab group in a way that you see the repositories in the same structure as they are in the Gitlab group.

usage: gsc [flags]

  -all
        Clone all projects, default is true
  -concurrency int
        Number of concurrent workers, default is 10 (default 10)
  -group int
        Clone projects from the given group ID
  -host string
        GitLab hostname, default is https://gitlab.com (default "https://gitlab.com")
  -insecure
        Allow insecure connection to your GitLab instance, default is false
  -matcher string
        Clone projects that match project name with the given regex
  -quiet
        Quiet bypasses the confirmation prompt, default is false
  -recursive
        Clone projects recursively, default is false
  -token string
        GitLab token for authentication
  -username string
        GitLab username for authentication
  -verbose
        Verbose output, default is false
  -version
        Print version information
```
