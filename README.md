# ConvCom

ConvCom is a lightweight Go tool designed to streamline the creation of commit messages following the [Conventional Commits specifications](https://www.conventionalcommits.org/en/v1.0.0/). Through an interactive command-line interface, ConvCom guides users in crafting structured and consistent commit messages.

## Commit Message Format

Each commit message should be structured as follows:

```bash
<type>(<scope>)<!>: <short summary>
  │       │     |       │
  │       │     |       └─⫸ Summary in present tense. Not capitalized. 
  |       |     |
  |       |     └─⫸ Exclamation mark for breaking change
  |       |
  No period at the end.
  │       │
  │       └─⫸ Commit Scope: api|formatter|taskmanager|taskstorage|docs
  │
  └─⫸ Commit Type: build|ci|chore|docs|feat|fix|perf|refactor|style|test
```
* `<type>`: The type of change being made (e.g., feat, fix, docs, etc.)
* `<scope>`: (Optional) The scope of the change (e.g., ui, api, build). This can be omitted if not relevant.
* `<short summary>`: A brief description of the change.


## Commit types
* `build`: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)
* `ci`: Changes to CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)
* **`chore`: Changes which doesn't change source code or tests e.g. changes to the build process, auxiliary tools, libraries**
* `docs`: Documentation only changes
* **`feat`: A new feature**
* **`fix`: A bug fix**
* `perf`: A code change that improves performance
* `refactor`:  A code change that neither fixes a bug nor adds a feature
* `revert`: Revert something
* `style`: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* `test`: Adding missing tests or correcting existing tests