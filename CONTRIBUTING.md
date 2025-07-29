# Contributing

## Local development workflow

### Prerequisites

- Docker installed (`docker` command available)

### Building the Docker image

To verify locally that the Docker image can be built:

```sh
docker build .
```

### Verifying the Docker image

To also verify that the app runs in the resulting image:

```sh
docker run --rm "`docker build -q .`"
```

To verify that the installation of the default Terragrunt and Terraform versions works as expected:

```sh
docker run --rm --entrypoint bash "`docker build -q .`" -c 'terragrunt run -- version'
```

## Design decisions

### Requirements

#### Functional

- Detects changed files and only runs Terraform in Terragrunt configurations that depend on them
- Can act as a drop-in replacement for the `terragrunt run` command in any Terragrunt-enabled repository
- Supports multiple hosting platforms (local & GitHub, initially) – similar to Renovate

#### Non-functional

- Easy to install
- Easy to use

### Programming language

**Go**, because both Terraform and Terragrunt are also written in it. It’s the only language that is guaranteed to find at least some common ground among users and maintainers of the tool. As a developer, it also makes it easier to reuse and move code between the projects involved.

### Name

The name should bear some resemblance to “Terraform” and “Terragrunt” and it should start with the letter “h.” This naming scheme leads to a convenient and somewhat curious set of aliases in interactive shells: `tf`, `tg` and `th` (`f` → `g` → `h`). I’m also fairly sure that the alias `th` is not in widespread use, yet.

I’ve decided to go with “**Terrahoot**,” initially.

Possible alternatives:

- Terra**huff**

### Tag line?

- “Making Terragrunt and Terraform actually _fun_ to use.”
- “A hoot and a half for Terragrunt and Terraform users.”
- “Hoot’n and toot’n”
