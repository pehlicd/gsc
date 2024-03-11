# gsc (Gitlab Structured Cloner)

gsc is a tool to help you clone all the repositories from a Gitlab group in a way that you see the repositories in the same structure as they are in the Gitlab group.

## Usage

```bash
Usage of gsc:
  -all
        Clone all projects, default is true (default true)
  -concurrency int
        Number of concurrent workers, default is 10 (default 10)
  -group int
        Clone projects from the given group ID
  -host string
        GitLab hostname, default is https://gitlab.com (default "https://gitlab.com")
  -insecure
        Allow insecure connection to your GitLab instance, default is false
  -matcher string
        Clone projects that match the given regex
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
usage: gsc [flags]
```