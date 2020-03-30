## Trabajo de hoy

Applicacion que ayude a analizar un repositorio de Git.

## TODO

- [x] Comenzar el app
- [x] Aprender a manejar file systems en Go
- [x] Aprender a escribir archivos en Go para escribir .gitattribute, para Git LFS
- [ ] Abstract command execution.
  - Function should receive the command along with the arguments needed.
  - This would need to be used over and over for the different git commands.
- [ ] Preprocess the repo to configure Git LFS correctly
  - Should I group files by size ranges???
  - Should size determine a processing path?
- [ ] Keep a list of files to add to Git
  - Large files (over 40MB) should be committed by themselves.
  - Should smaller files be batched????
- [ ] Add function documentation

## Git commands to execute

- git add
- git commit
- git push
- git lfs track
- git init

