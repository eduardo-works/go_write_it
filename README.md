# WritΞIt

WritΞIt is a Go command-line tool designed to help writers to maintain a consistent file structure for their writing projects, and compile the [Markdown](https://www.markdownguide.org) manuscripts into [LibreOffice](https://www.libreoffice.org) files.

## 📖 Features

- Create and manage writing projects.
- Organize manuscripts into chapters, parts and different versions.
- Maintain consistent file naming conventions.
- Keep writing projects organized using a predictable directory structure.

## 🚧 Project types

WritΞIt supports three types of writing projects:

| Type         | Description                                            |
| ------------ | ------------------------------------------------------ |
| `standalone` | A single manuscript, such as a short story or essay.   |
| `chaptered`  | A large manuscript organized into chapters.            |
| `parted`     | A large manuscript organized into parts and chapters.  |

## ⚜️ Getting started

### 🎬 Create a project

```bash
WritΞIt project create ./books "My First Novel" \
  --author="John Doe" \
  --type=chaptered
```

For a project divided into multiple parts:

```bash
WritΞIt project create ./books "My Trilogy" \
  --author="John Doe" \
  --type=parted \
  --parts=3
```

By default, projects use the `standalone` type.

### Manage your project

```bash
WritΞIt project info ./books/my-first-novel
WritΞIt project rename ./books/my-first-novel "My New Title"
WritΞIt project update ./books/my-first-novel --author="Jane Doe"
WritΞIt project add-part ./books/my-trilogy
```

### Normalize filenames

WritΞIt can identify and rename manuscript files that do not follow the project's naming conventions.

Preview the changes without modifying any files:

```bash
WritΞIt project dry-run ./books/my-first-novel
```

Apply the changes:

```bash
WritΞIt project normalize ./books/my-first-novel
```

Files that cannot be normalized are reported for manual review.

## Available commands

| Command             | Description                     |
| ------------------- | ------------------------------- |
| `project create`    | Create a new writing project.   |
| `project info`      | Display project information.    |
| `project rename`    | Rename a project.               |
| `project update`    | Update project metadata.        |
| `project add-part`  | Add a new part to a project.    |
| `project dry-run`   | Preview filename normalisation. |
| `project normalize` | Normalise manuscript filenames. |

Use `WritΞIt --help` or `WritΞIt project --help` for more information.

## Planned features

The following commands are under development.

### File management

- `file create` — Create a new manuscript file.
- `file version` — Create a new version of an existing manuscript.

### Manuscript assembly and compilation

- `master assemble` — Combine manuscript files into a single Markdown document.
- `master dry-run` — Preview the files selected for assembly.
- `master compile` — Convert the assembled manuscript into a LibreOffice document.
- `master build` — Assemble and compile the manuscript in one operation.

## Philosophy

**Write in Markdown. Organize with WritΞIt.**

WritΞIt focuses on project organization rather than text editing. Writers remain free to choose their own editor, manage their files locally and maintain control over their manuscripts.
