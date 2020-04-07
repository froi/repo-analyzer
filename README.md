## Trabajo de hoy

Applicacion que ayude a analizar un repositorio de Git.

## TODO

- [x] Comenzar el app
- [x] Aprender a manejar file systems en Go
- [x] Aprender a escribir archivos en Go para escribir .gitattribute, para Git LFS
- [x] Abstract command execution.
  - Function should receive the command along with the arguments needed.
  - This would need to be used over and over for the different git commands.
- [ ] Preprocess the repo to configure Git LFS correctly
  - Should I group files by size ranges??? - ended up dividing the file in lfs and non lfs groups
  - Should size determine a processing path? - anything over 40MB will be added to LFS
- [x] Keep a list of files to add to Git
  - Large files (over 40MB) should be committed by themselves. - See above
  - Should smaller files be batched???? - The smaller files are being batched in groups of 10.
- [ ] Add function documentation
- [ ] Add CLI arguments
  - repo_analyzer <branch> --init --remote https://github.com/user/repo.git --remote_name "something"
- [x] Add dannywolfmx as a coauthor


## Git commands to execute

- [x] git add
- [x] git commit
- [x] git push
- [x] git lfs track
- [ ] git init

