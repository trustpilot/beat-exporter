workflow "build-and-deploy" {
  on = "push"
  resolves = "build"
}

action "build" {
  uses="cedrickring/golang-action/go1.12@1.3.0"
}