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
- Drop-in replacement for the `terragrunt run` command in any Terragrunt-enabled repository
- Supports multiple platforms (local & GitHub, initially) – similar to Renovate
- Does the intuitive "right" thing in CI and when executed locally
- Implements merge-to-apply with failure handling

#### Non-functional

- Easy to install
- Easy to use

### Local workflow

_TBD_

### CI/CD workflow

#### Phase 1: Plan & feedback

Trigger events:

- Opening a pull request
- Updating the head branch of an open pull request
- Re-opening a closed pull request with unmerged changes

Steps:

- Determine what changed between base and head
  - Determine deleted units
  - Determine changed files
- Run "plan -destroy" for deleted units in base
  - Also plan to destroy all downstream dependencies
- Run "plan" for units affected by changed files
  - Also plan all upstream dependencies
  - Note which downstream units might be affected by changes
- Leave feedback on pull request
  - Update existing comment with latest workflow run result
  - Limit feedback to a single comment and/or job status report
  - Post a summary report for all units (succeeded, failed, skipped)
    - Mention all downstream dependencies which might be affected by
      changes in outputs and should be included in the apply phase
  - Post details for each unit that had changes or errors, only if all details
    fit within the character limit; otherwise, post only change summaries for
    each unit and link to detailed logs as workflow attachments

Notes:

- Merging the pull request should be prevented until all units and all workflow
  steps have succeeded

#### Phase 2: Apply & feedback

Trigger events:

- Merging a pull request (not just closing it unmerged)

Steps:

- Repeat all planning steps, except for feedback
  - Why no to use saved plans?
    - Mainly, because changes in the outputs of one unit can affect downstream
      units, so saved plans can be plain *wrong*
    - Secondly, because they can become stale and cause the apply phase to fail,
      if too much time has passed and other changes have been made in between
- Open/update an issue on failure
- Leave feedback on pull request
  - Summary report on success with link to details in workflow logs
  - Link to created/updated issue on failure
- Create a release on success (semantic version tag, change log)

Questions:

- Should drift detection be triggered automatically on apply failure?
- Should failures during drift detection create/update the same issue?
- Should a single issue or multiple issues be used to handle apply failures?

#### Phase 3: Deploy Terraform module releases

Triggers:

- Renovate updates Terraform module versions in production environments

Steps:

- Renovate pushes changes to a new branch and opens a pull request
- See [Phase 1: Plan & feedback](#phase-1-plan--feedback)

Notes:

- Renovate's automerge feature can be employed at this point to automate
  rollouts in production
  - Semantic versions are available to guide the decision of whether or
    not to enable automerge for a certain module

### Programming language

**Go**, because both Terraform and Terragrunt are also written in it. It’s the only language that is guaranteed to find at least some common ground among users and maintainers of the tool. As a developer, it also makes it easier to reuse and move code between the projects involved.

### Name

The name should bear some resemblance to “Terraform” and “Terragrunt” and it should start with the letter “h.” This naming scheme leads to a convenient and somewhat curious set of aliases in interactive shells: `tf`, `tg` and `th` (`f` → `g` → `h`). I’m also fairly sure that the alias `th` is not in widespread use, yet.

I’ve decided to go with “**Terrahoot**” as a working title, but intend to give it a more expressive name if and when it is becoming really useful.

Alternatives considered:

- Terra**huff**

### Tag line?

- “Making Terragrunt and Terraform actually _fun_ to use.”
- “A hoot and a half for Terragrunt and Terraform users.”
- “Hoot’n and toot’n”
