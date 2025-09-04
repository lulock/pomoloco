# Pomoloco :tomato:
A simple pomodoro CLI app written in Go.

```
username@hostname:~/pomoloco$ ./pomoloco

  ┌────────────────────────────────────────┐
  │                                        │
  │  Rivers know this: there is no hurry.  │
  │  We shall get there some day.          │
  │    -- A.A. Milne                       │
  │                                        │
  └────────────────────────────────────────┘

  Go go go! Time to focus.

  13:27  *  ███████████████████████████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░

```
## Roadmap
### mvp TODOs
- [x] pomo command: default 25-minute work session
- [x] loco command: default 5-minute break
- [x] exit command
- [x] help command
- [x] visual bar that shrinks with countdown

### level-up TODOs
- [x] bubbletea? lipgloss? yes
- [x] notification on end (sound?)
- [x] custom time for both pomo and loco sessions
- [x] press esc or q to end session
- [x] fetch random motivational quote from zenquotes.io
- [x] refresh quote when r is pressed
- [x] skip to next session when n or enter is pressed
- [x] loop pomo-loco sessions
- [x] refactor out styling logic
- [ ] session title optional
- [ ] notes / reflections optional
- [ ] save sessions and notes in database
- [ ] refactor out TUI / App logic
- [ ] refactor out Viper config logic
- [ ] add teatests

## Contributing

### clone the repo

```bash
git clone https://github.com/lulock/pomoloco
cd pomoloco
```

### build 

```bash
go build
```

### run executable

```bash
./pomoloco
```
### run tests 
no tests yet... 

### submit a pull request
To contribute, fork the repository and open a pull request to 'main'.

